package server

import (
	"context"
	"github.com/nrc-no/core/api/pkg/apis/core"
	coreinstall "github.com/nrc-no/core/api/pkg/apis/core/install"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	discoveryinstall "github.com/nrc-no/core/api/pkg/apis/discovery/install"
	coreclient "github.com/nrc-no/core/api/pkg/client/core"
	"github.com/nrc-no/core/api/pkg/client/informers"
	restclient "github.com/nrc-no/core/api/pkg/client/rest"
	"github.com/nrc-no/core/api/pkg/controllers/customresource"
	formdefinitions3 "github.com/nrc-no/core/api/pkg/controllers/formdefinitions"
	customresourcedefinition2 "github.com/nrc-no/core/api/pkg/registry/core/customresourcedefinition"
	formdefinitions2 "github.com/nrc-no/core/api/pkg/registry/core/formdefinitions"
	generic2 "github.com/nrc-no/core/api/pkg/registry/generic"
	rest2 "github.com/nrc-no/core/api/pkg/registry/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"net"
	"net/http"
	"time"
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	coreinstall.Install(Scheme)
	discoveryinstall.Install(Scheme)
}

type MongoConfig struct {
	Address string
}

type Config struct {
	ListenAddress         net.IP
	RESTOptionsGetter     generic2.RESTOptionsGetter
	BuildHandlerChainFunc func(handler http.Handler, config *Config) http.Handler
	CRDRestOptionsGetter  generic2.RESTOptionsGetter
	LoopbackClientConfig  *restclient.Config
	Listener              net.Listener
}

func (c *Config) Complete() *CompletedConfig {
	completedConfig := &CompletedConfig{
		Config: c,
	}
	return completedConfig
}

type CompletedConfig struct {
	*Config
}

// New creates and configures a new Server from the CompletedConfig
func (c *CompletedConfig) New(ctx context.Context) (*Server, error) {

	// Builds the handler chain. This will register all the filters and middlewares and so on
	// that need to be ran before the dispatching the request to go-restful or non-go-restful
	handlerChainBuilder := func(handler http.Handler) http.Handler {
		return c.BuildHandlerChainFunc(handler, c.Config)
	}

	// Creates the API server HTTP handler
	// The API server handler has both a
	// go-restful container, that tries to match the request first.
	// it then tries to match the request with a non-go-restful handler.
	//
	// This handler is ran after the handler chain (after the filters),
	// and either dispatches the request to a go-restful container,
	// or to a non-gorestful container.
	//
	// The go-restful container is for regular restful WebService that
	// are built in advance. The non-go-restful mux is for either
	// custom resource definitions, or for additional endpoints
	// that we want to register in the API
	apiServerHandler := NewAPIServerHandler(handlerChainBuilder)

	// Create the API server
	server := &Server{
		ctx:                  ctx,
		handler:              apiServerHandler,
		listener:             c.Listener,
		postStartHooks:       map[string]postStartHookEntry{},
		LoopbackClientConfig: c.LoopbackClientConfig,
	}

	// Installs the known resource handlers in the API (in the go-restful container)
	// These include FormDefinitions and CustomResourceDefinitions
	if err := c.installApiGroups(apiServerHandler); err != nil {
		return nil, err
	}

	// Create the core.nrc.no/v1 client that will be used to create
	// the controllers/informers. It's using the LoopbackClientConfig
	// so that it can reach localhost
	cli, err := coreclient.NewForConfig(c.LoopbackClientConfig)
	if err != nil {
		return nil, err
	}

	// starts the informers on server startup
	server.coreInformers = informers.NewSharedInformerFactory(cli, time.Minute*5)
	server.AddPostStartHookOrDie("core-informers", func(context PostStartHookContext) error {
		server.coreInformers.Start(context.StopCh)
		return nil
	})

	// Installs the CustomResources in the API
	if err := c.installCustomResources(server, apiServerHandler); err != nil {
		return nil, err
	}

	// Creates the FormDefinitionController that maintains the
	// CRDs corresponding to the form definitions
	formDefController := formdefinitions3.NewFormDefinitionController(
		cli,
		server.coreInformers.Core().V1().FormDefinitions(),
		server.coreInformers.Core().V1().CustomResourceDefinitions(),
	)
	server.AddPostStartHookOrDie("formdefinition-controllers", func(context PostStartHookContext) error {
		formDefController.Run(context.StopCh)
		return nil
	})

	return server, nil
}

// installCustomResources installs the required handlers to serve dynamic
// CustomResources from the API
// This will install the Discovery handlers as well as the regular resource handlers
func (c *CompletedConfig) installCustomResources(server *Server, apiServerHandler *APIServerHandler) error {

	// This creates the HTTP handler for custom resource version discovery.
	// It's able to dynamically serve the different custom resource APIVersions
	crdVersionDiscoveryHandler := customresource.NewCRDVersionDiscoveryHandler(http.NotFoundHandler())

	// This creates the HTTP handler for custom resource group discovery
	// It's able to dynamically serve the different versions present in a custom APIGroup
	crdGroupDiscoveryHandler := customresource.NewCRDGroupDiscoveryHandler(http.NotFoundHandler())

	// This creates the controller that reconciles the CRDs and maintains
	// the crdVersionDiscoveryHandler and crdGroupDiscoveryHandler endpoints
	crdDiscoveryController := customresource.NewCRDDiscoveryController(
		crdVersionDiscoveryHandler,
		crdGroupDiscoveryHandler,
		Codecs,
		server.coreInformers.Core().V1().CustomResourceDefinitions(),
	)

	// Starts the crdDiscoveryHandler
	server.AddPostStartHookOrDie("crd-controllers", func(context PostStartHookContext) error {
		discoverySyncedCh := make(chan struct{})
		go crdDiscoveryController.Run(context.StopCh, discoverySyncedCh)
		select {
		case <-context.StopCh:
		case <-discoverySyncedCh:
		}
		return nil
	})

	// This creates the HTTP handler that dynamically serves the request aimed
	// at custom resources. It will dispatch the request to the crdVersionDiscoveryHandler
	// or the crdGroupDiscoveryHandler if applicable. Otherwise, it will perform the
	// generic REST operations to serve the custom resource from storage.
	crdHandler, err := customresource.NewCustomResourceDefinitionHandler(
		c.CRDRestOptionsGetter,
		Codecs,
		server.coreInformers.Core().V1().CustomResourceDefinitions(),
		crdVersionDiscoveryHandler,
		crdGroupDiscoveryHandler)
	if err != nil {
		return err
	}

	// Registers the crdHandler inside of the non-goRestfulMux
	apiServerHandler.NonGoRestfulMux.Handle("/apis", crdHandler)
	apiServerHandler.NonGoRestfulMux.HandlePrefix("/apis/", crdHandler)

	return nil
}

// installApiGroups installs the known API groups in the HTTP handler
func (c *CompletedConfig) installApiGroups(apiServerHandler *APIServerHandler) error {
	var apiGroups []*APIGroupInfo
	coreApiGroup, err := createCoreAPIGroupInfo(c.RESTOptionsGetter)
	if err != nil {
		return err
	}
	apiGroups = append(apiGroups, &coreApiGroup)

	if err := installApiGroups(apiServerHandler.GoRestfulContainer, "/apis", apiGroups...); err != nil {
		return err
	}

	return nil
}

// createCoreAPIGroupInfo creates the APIGroupInfo for the core API
func createCoreAPIGroupInfo(restOptionsGetter generic2.RESTOptionsGetter) (APIGroupInfo, error) {
	coreApiGroup := NewDefaultAPIGroup(core.GroupName, Scheme, metav1.ParameterCodec, Codecs)
	storage := map[string]rest2.Storage{}

	// Storage for FormDefinitions
	formDefinitionsStorage, err := formdefinitions2.NewREST(Scheme, restOptionsGetter)
	if err != nil {
		return APIGroupInfo{}, err
	}
	storage["formdefinitions"] = formDefinitionsStorage

	// Storage for CustomResourceDefinitions
	customResourceDefinitionsStorage, err := customresourcedefinition2.NewREST(Scheme, restOptionsGetter)
	if err != nil {
		return APIGroupInfo{}, err
	}
	storage["customresourcedefinitions"] = customResourceDefinitionsStorage

	// register the storage for the "core v1" version
	coreApiGroup.VersionedResourcesStorageMap[corev1.SchemeGroupVersion.Version] = storage

	return coreApiGroup, nil
}
