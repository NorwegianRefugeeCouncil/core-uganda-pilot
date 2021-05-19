package server

import (
	"fmt"
	"github.com/nrc-no/core/api/pkg/apis/core"
	coreinstall "github.com/nrc-no/core/api/pkg/apis/core/install"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/api/pkg/generated/openapi"
	"github.com/nrc-no/core/api/pkg/registry/core/formdefinition"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsinstall "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/install"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiextensionsinformers "k8s.io/apiextensions-apiserver/pkg/client/informers/externalversions"
	"k8s.io/apiextensions-apiserver/pkg/controller/apiapproval"
	"k8s.io/apiextensions-apiserver/pkg/controller/establish"
	"k8s.io/apiextensions-apiserver/pkg/controller/finalizer"
	"k8s.io/apiextensions-apiserver/pkg/controller/nonstructuralschema"
	"k8s.io/apiextensions-apiserver/pkg/controller/status"
	"k8s.io/apiextensions-apiserver/pkg/registry/customresourcedefinition"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/endpoints/discovery"
	"k8s.io/apiserver/pkg/endpoints/filters"
	openapi2 "k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/apiserver/pkg/storageversion"
	restclient "k8s.io/client-go/rest"
	"k8s.io/kube-openapi/pkg/common"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	coreinstall.Install(Scheme)
	apiextensionsinstall.Install(Scheme)
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
}

type RecommendedConfig struct {
	Config
}

func NewRecommendedConfig(codecs serializer.CodecFactory) *RecommendedConfig {
	return &RecommendedConfig{
		Config: *NewConfig(codecs),
	}
}

type Config struct {
	RESTOptionsGetter      generic.RESTOptionsGetter
	OpenAPIConfig          *common.Config
	BuildHandlerChainFunc  func(apiHandler http.Handler, c *Config) (secure http.Handler)
	Serializer             runtime.NegotiatedSerializer
	ExternalAddress        string
	RequestInfoResolver    request.RequestInfoResolver
	LegacyAPIGroupPrefixes sets.String
	LoopbackClientConfig   *restclient.Config
}

func (c *Config) Complete() CompletedConfig {
	if c.RequestInfoResolver == nil {
		c.RequestInfoResolver = NewRequestInfoResolver(c)
	}
	return CompletedConfig{
		*c,
	}
}

type CompletedConfig struct {
	Config
}

func NewConfig(codecs serializer.CodecFactory) *Config {
	return &Config{
		Serializer:             codecs,
		BuildHandlerChainFunc:  DefaultBuildHandlerChain,
		LegacyAPIGroupPrefixes: sets.NewString(DefaultLegacyAPIPrefix),
	}
}

func (c CompletedConfig) New() (*APIServer, error) {

	handlerChainBuilder := func(handler http.Handler) http.Handler {
		return c.BuildHandlerChainFunc(handler, &c.Config)
	}
	delegate := http.NotFoundHandler()
	apiServerHandler := server.NewAPIServerHandler("", c.Serializer, handlerChainBuilder, delegate)

	s := &APIServer{
		Serializer:                 c.Serializer,
		DiscoveryGroupManager:      discovery.NewRootAPIsHandler(discovery.DefaultAddresses{DefaultAddress: c.ExternalAddress}, c.Serializer),
		EquivalentResourceRegistry: runtime.NewEquivalentResourceRegistry(),
		minRequestTimeout:          time.Duration(1800) * time.Second,
		maxRequestBodyBytes:        int64(3 * 1024 * 1024),
		StorageVersionManager:      storageversion.NewDefaultManager(),
		Handler:                    apiServerHandler,
		RESTOptionsGetter:          c.RESTOptionsGetter,
		preShutdownHooks:           map[string]preShutdownHookEntry{},
		postStartHooks:             map[string]postStartHookEntry{},
		disabledPostStartHooks:     sets.NewString(),
		healthzChecks:              []healthz.HealthChecker{},
		healthzChecksInstalled:     false,
		healthzLock:                sync.Mutex{},
		livezChecks:                []healthz.HealthChecker{},
		livezChecksInstalled:       false,
		livezClock:                 clock.RealClock{},
		livezGracePeriod:           10 * time.Second,
		livezLock:                  sync.Mutex{},
		postStartHookLock:          sync.Mutex{},
		postStartHooksCalled:       false,
		preShutdownHookLock:        sync.Mutex{},
		preShutdownHooksCalled:     false,
		readyzChecks:               []healthz.HealthChecker{},
		readyzChecksInstalled:      false,
		readyzLock:                 sync.Mutex{},
		readinessStopCh:            make(chan struct{}),
	}

	s.openAPIConfig = server.DefaultOpenAPIConfig(openapi.GetOpenAPIDefinitions, openapi2.NewDefinitionNamer(Scheme))
	s.openAPIConfig.Info.Title = "Core API"
	s.openAPIConfig.Info.Version = "0.1"

	// Install FormDefinitions
	apiResourceConfig := server.NewDefaultAPIGroupInfo(core.GroupName, Scheme, metav1.ParameterCodec, Codecs)
	storage := map[string]rest.Storage{}
	formDefinitionsStorage, err := formdefinition.NewREST(Scheme, s.RESTOptionsGetter)
	if err != nil {
		return nil, err
	}
	storage["formdefinitions"] = formDefinitionsStorage
	apiResourceConfig.VersionedResourcesStorageMap[corev1.SchemeGroupVersion.Version] = storage

	// Install CustomResourceDefinitions
	extensionsGroupInfo := server.NewDefaultAPIGroupInfo(apiextensions.GroupName, Scheme, metav1.ParameterCodec, Codecs)
	extensionsStorage := map[string]rest.Storage{}
	customResourceDefinitionStorage, err := customresourcedefinition.NewREST(Scheme, c.RESTOptionsGetter)
	if err != nil {
		return nil, err
	}
	extensionsStorage["customresourcedefinitions"] = customResourceDefinitionStorage
	extensionsStorage["customresourcedefinitions/status"] = customresourcedefinition.NewStatusREST(Scheme, customResourceDefinitionStorage)
	extensionsGroupInfo.VersionedResourcesStorageMap[apiextensionsv1.SchemeGroupVersion.Version] = extensionsStorage

	// Install REST resources
	if err := s.InstallAPIGroups(&apiResourceConfig, &extensionsGroupInfo); err != nil {
		return nil, err
	}

	// Create API Extensions Client
	crdClient, err := clientset.NewForConfig(c.LoopbackClientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %v", err)
	}
	s.ApiExtensionsInformers = apiextensionsinformers.NewSharedInformerFactory(crdClient, 5*time.Minute)

	// Handles CRD version discovery
	versionDiscoveryHandler := &versionDiscoveryHandler{
		discovery: map[schema.GroupVersion]*discovery.APIVersionHandler{},
		delegate:  delegate,
	}
	// Handles CRD group discovery
	groupDiscoveryHandler := &groupDiscoveryHandler{
		discovery: map[string]*discovery.APIGroupHandler{},
		delegate:  delegate,
	}

	// EstablishingController controls when CRD are
	// considered as "established" (it's a runtime.Object condition)
	establishingController := establish.NewEstablishingController(s.ApiExtensionsInformers.Apiextensions().V1().CustomResourceDefinitions(), crdClient.ApiextensionsV1())

	// HTTP Handler for CRDs
	crdHandler, err := NewCustomResourceDefinitionHandler(
		versionDiscoveryHandler,
		groupDiscoveryHandler,
		s.ApiExtensionsInformers.Apiextensions().V1().CustomResourceDefinitions(),
		delegate,
		c.RESTOptionsGetter,
		nil,
		establishingController,
		nil,
		nil,
		1,
		nil,
		5*time.Second,
		2*time.Second,
		extensionsGroupInfo.StaticOpenAPISpec,
		(10 * 1024 * 1024),
	)

	// Handle the CRDs without go-restful container
	s.Handler.NonGoRestfulMux.Handle("/apis", crdHandler)
	s.Handler.NonGoRestfulMux.HandlePrefix("/apis/", crdHandler)

	// Create the CRD controllers
	discoveryController := NewDiscoveryController(s.ApiExtensionsInformers.Apiextensions().V1().CustomResourceDefinitions(), versionDiscoveryHandler, groupDiscoveryHandler)
	namingController := status.NewNamingConditionController(s.ApiExtensionsInformers.Apiextensions().V1().CustomResourceDefinitions(), crdClient.ApiextensionsV1())
	nonStructuralSchemaController := nonstructuralschema.NewConditionController(s.ApiExtensionsInformers.Apiextensions().V1().CustomResourceDefinitions(), crdClient.ApiextensionsV1())
	apiApprovalController := apiapproval.NewKubernetesAPIApprovalPolicyConformantConditionController(s.ApiExtensionsInformers.Apiextensions().V1().CustomResourceDefinitions(), crdClient.ApiextensionsV1())
	finalizingController := finalizer.NewCRDFinalizer(
		s.ApiExtensionsInformers.Apiextensions().V1().CustomResourceDefinitions(),
		crdClient.ApiextensionsV1(),
		crdHandler,
	)
	// openapiController := openapicontroller.NewController(s.ApiExtensionsInformers.Apiextensions().V1().CustomResourceDefinitions())

	// Start the api extension informers
	s.AddPostStartHookOrDie("start-apiextension-informers", func(context PostStartHookContext) error {
		s.ApiExtensionsInformers.Start(context.StopCh)
		return nil
	})

	// Start the api extensions controllers
	s.AddPostStartHookOrDie("start-apiextensions-controllers", func(context PostStartHookContext) error {

		//TODO: openapi
		go namingController.Run(context.StopCh)
		go establishingController.Run(context.StopCh)
		go nonStructuralSchemaController.Run(5, context.StopCh)
		go apiApprovalController.Run(5, context.StopCh)
		go finalizingController.Run(5, context.StopCh)

		discoverySyncedCh := make(chan struct{})
		go discoveryController.Run(context.StopCh, discoverySyncedCh)
		select {
		case <-context.StopCh:
		case <-discoverySyncedCh:
		}
		return nil
	})

	// Wait for CRDs to sync
	s.AddPostStartHookOrDie("crd-informer-synced", func(context PostStartHookContext) error {
		return wait.PollImmediateUntil(100*time.Millisecond, func() (done bool, err error) {
			return s.ApiExtensionsInformers.Apiextensions().V1().CustomResourceDefinitions().Informer().HasSynced(), nil
		}, context.StopCh)
	})

	return s, nil

}

func DefaultBuildHandlerChain(apiHandler http.Handler, config *Config) http.Handler {
	handler := filters.WithRequestInfo(apiHandler, config.RequestInfoResolver)
	return handler
}

func NewRequestInfoResolver(c *Config) *request.RequestInfoFactory {
	apiPrefixes := sets.NewString(strings.Trim(APIGroupPrefix, "/")) // all possible API prefixes
	legacyAPIPrefixes := sets.String{}                               // APIPrefixes that won't have groups (legacy)
	for legacyAPIPrefix := range c.LegacyAPIGroupPrefixes {
		apiPrefixes.Insert(strings.Trim(legacyAPIPrefix, "/"))
		legacyAPIPrefixes.Insert(strings.Trim(legacyAPIPrefix, "/"))
	}
	return &request.RequestInfoFactory{
		APIPrefixes:          apiPrefixes,
		GrouplessAPIPrefixes: legacyAPIPrefixes,
	}
}
