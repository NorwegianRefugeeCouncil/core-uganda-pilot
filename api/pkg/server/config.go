package server

import (
	"context"
	"github.com/nrc-no/core/api/pkg/apis/core"
	coreinstall "github.com/nrc-no/core/api/pkg/apis/core/install"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/api/pkg/apis/discovery"
	discoveryinstall "github.com/nrc-no/core/api/pkg/apis/discovery/install"
	discoveryv1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	"github.com/nrc-no/core/api/pkg/auth/keycloak"
	"github.com/nrc-no/core/api/pkg/client/informers"
	restclient "github.com/nrc-no/core/api/pkg/client/rest"
	"github.com/nrc-no/core/api/pkg/client/typed"
	"github.com/nrc-no/core/api/pkg/controllers/customresource"
	formdefinitions3 "github.com/nrc-no/core/api/pkg/controllers/formdefinitions"
	"github.com/nrc-no/core/api/pkg/controllers/operatingscope"
	"github.com/nrc-no/core/api/pkg/controllers/registration"
	discoveryhandlers "github.com/nrc-no/core/api/pkg/endpoints/discovery"
	customresourcedefinitionstorage "github.com/nrc-no/core/api/pkg/registry/core/customresourcedefinition"
	formdefinitionsstorage "github.com/nrc-no/core/api/pkg/registry/core/formdefinitions"
	operatingscopestorage "github.com/nrc-no/core/api/pkg/registry/core/operatingscope"
	"github.com/nrc-no/core/api/pkg/registry/core/users"
	"github.com/nrc-no/core/api/pkg/registry/discovery/apiservice"
	"github.com/nrc-no/core/api/pkg/registry/generic"
	"github.com/nrc-no/core/api/pkg/registry/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"net"
	"net/http"
	"strings"
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
	RESTOptionsGetter     generic.RESTOptionsGetter
	BuildHandlerChainFunc func(handler http.Handler, config *Config) http.Handler
	CRDRestOptionsGetter  generic.RESTOptionsGetter
	LoopbackClientConfig  *restclient.Config
	Listener              net.Listener
	OidcClientID          string
	OidcClientSecret      string
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
		listedPathProvider:   apiServerHandler,
	}

	keycloakClient, err := keycloak.NewKeycloakClient("http://localhost:8080", c.OidcClientID, c.OidcClientSecret)
	if err != nil {
		return nil, err
	}

	// Installs the known resource handlers in the API (in the go-restful container)
	// These include FormDefinitions and CustomResourceDefinitions
	c.installApiGroupsOrDie(apiServerHandler, keycloakClient, "Core")

	// Create the core.nrc.no/v1 client that will be used to create
	// the controllers/informers. It's using the LoopbackClientConfig
	// so that it can reach localhost
	cli := c.createClientSetOrDie()

	// starts the core informers on server startup
	c.startInformersOrDie(server, cli)

	// Installs the CustomResources in the API
	c.installCustomResourcesOrDie(server, apiServerHandler)

	// Creates the FormDefinitionController that maintains the
	// CRDs corresponding to the form definitions
	c.startFormDefinitionControllerOrDie(server, cli)

	// Creates the APIService registration controller
	// Maintains the list of APIServices registered for the API
	autoRegisterController := c.createAutoRegisterController(server, cli)

	c.registerApiServices(server, autoRegisterController)

	// Creates the CRDRegistrationController
	// Registers the APIServices corresponding to the
	// CustomResourceDefinitions
	crdRegistrationController := c.createCRDRegistrationController(server, autoRegisterController)

	// Start the CRD and APIService controllers
	c.startAutoRegistrationControllers(server, crdRegistrationController, autoRegisterController)

	operatingScopeController := operatingscope.NewOperatingScopeController(
		server.informers.Core().V1().OperatingScopes(),
		keycloakClient,
		"Core")

	server.AddPostStartHookOrDie("operatingscope-controllers", func(context PostStartHookContext) error {
		syncedCh := make(chan struct{})
		operatingScopeController.Run(context.StopCh, syncedCh)
		select {
		case <-context.StopCh:
		case <-syncedCh:
		}
		return nil
	})

	apisHandler := discoveryhandlers.NewApisHandler(
		Codecs,
		server.informers.Discovery().V1().APIServices().Lister())

	apiServerHandler.NonGoRestfulMux.Handle("/apis", apisHandler)
	// apiServerHandler.NonGoRestfulMux.UnlistedHandle("/apis/", apisHandler)

	return server, nil
}

func (c *CompletedConfig) registerApiServices(
	server *Server,
	registration registration.AutoAPIServiceRegistration,
) {
	for _, curr := range server.listedPathProvider.ListedPaths() {
		if !strings.HasPrefix(curr, "/apis") {
			continue
		}
		tokens := strings.Split(curr, "/")
		if len(tokens) != 4 {
			continue
		}
		apiService := makeAPIService(schema.GroupVersion{Group: tokens[2], Version: tokens[3]})
		if apiService == nil {
			continue
		}
		registration.AddAPIServiceToSyncOnStart(apiService)
	}
}

func makeAPIService(gv schema.GroupVersion) *discoveryv1.APIService {
	return &discoveryv1.APIService{
		ObjectMeta: metav1.ObjectMeta{Name: gv.Version + "." + gv.Group},
		Spec: discoveryv1.APIServiceSpec{
			Group:   gv.Group,
			Version: gv.Version,
		},
	}
}

func (c *CompletedConfig) startAutoRegistrationControllers(
	server *Server,
	crdRegistrationController *customresource.CRDRegistrationController,
	autoRegistrationController *registration.AutoRegisterController,
) {
	server.AddPostStartHookOrDie("autoregistration-controller", func(context PostStartHookContext) error {
		go crdRegistrationController.Run(context.StopCh)
		go func() {
			crdRegistrationController.WaitForInitialSync()
			autoRegistrationController.Run(context.StopCh)
		}()
		return nil
	})
}

func (c *CompletedConfig) createCRDRegistrationController(server *Server, autoRegisterController *registration.AutoRegisterController) *customresource.CRDRegistrationController {
	return customresource.NewCRDRegistrationController(
		server.informers.Core().V1().CustomResourceDefinitions(),
		autoRegisterController,
	)
}

func (c *CompletedConfig) createAutoRegisterController(server *Server, cli typed.Interface) *registration.AutoRegisterController {
	return registration.NewAutoRegisterController(
		server.informers.Discovery().V1().APIServices(),
		cli.DiscoveryV1(),
	)
}

func (c *CompletedConfig) createClientSetOrDie() typed.Interface {
	cli, err := typed.NewForConfig(c.LoopbackClientConfig)
	if err != nil {
		panic(err)
	}
	return cli
}

func (c *CompletedConfig) startInformersOrDie(server *Server, cli typed.Interface) {
	// starts the core informers on server startup
	server.informers = informers.NewSharedInformerFactory(cli, time.Minute*5)
	server.AddPostStartHookOrDie("informers", func(context PostStartHookContext) error {
		server.informers.Start(context.StopCh)
		return nil
	})
}

func (c *CompletedConfig) startFormDefinitionControllerOrDie(server *Server, cli typed.Interface) {
	formDefController := formdefinitions3.NewFormDefinitionController(
		cli,
		server.informers.Core().V1().FormDefinitions(),
		server.informers.Core().V1().CustomResourceDefinitions(),
	)
	server.AddPostStartHookOrDie("formdefinition-controllers", func(context PostStartHookContext) error {
		syncedCh := make(chan struct{})
		formDefController.Run(context.StopCh, syncedCh)
		select {
		case <-context.StopCh:
		case <-syncedCh:
		}
		return nil
	})
}

// installCustomResourcesOrDie installs the required handlers to serve dynamic
// CustomResources from the API
// This will install the Discovery handlers as well as the regular resource handlers
func (c *CompletedConfig) installCustomResourcesOrDie(server *Server, apiServerHandler *APIServerHandler) {

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
		server.informers.Core().V1().CustomResourceDefinitions(),
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
		server.informers.Core().V1().CustomResourceDefinitions(),
		crdVersionDiscoveryHandler,
		crdGroupDiscoveryHandler)
	if err != nil {
		panic(err)
	}

	// Registers the crdHandler inside of the non-goRestfulMux
	//apiServerHandler.NonGoRestfulMux.Handle("/apis", crdHandler)
	apiServerHandler.NonGoRestfulMux.HandlePrefix("/apis/", crdHandler)
}

// installApiGroups installs the known API groups in the HTTP handler
func (c *CompletedConfig) installApiGroupsOrDie(
	apiServerHandler *APIServerHandler,
	keycloakClient *keycloak.KeycloakClient,
	realmName string,
) {

	var apiGroups []*APIGroupInfo
	coreApiGroup, err := createCoreAPIGroupInfo(c.RESTOptionsGetter, keycloakClient, realmName)
	if err != nil {
		panic(err)
	}

	discoveryApiGroup, err := createDiscoveryAPIGroupInfo(c.RESTOptionsGetter)
	if err != nil {
		panic(err)
	}

	apiGroups = append(apiGroups, &coreApiGroup, &discoveryApiGroup)

	if err := installApiGroups(apiServerHandler.GoRestfulContainer, "/apis", apiGroups...); err != nil {
		panic(err)
	}

}

func createDiscoveryAPIGroupInfo(restOptionsGetter generic.RESTOptionsGetter) (APIGroupInfo, error) {
	discoveryApiGroup := NewDefaultAPIGroup(discovery.GroupName, Scheme, metav1.ParameterCodec, Codecs)
	storage := map[string]rest.Storage{}

	apiServicesStorage, err := apiservice.NewRESTStorage(Scheme, restOptionsGetter)
	if err != nil {
		return APIGroupInfo{}, err
	}
	storage["apiservices"] = apiServicesStorage
	discoveryApiGroup.VersionedResourcesStorageMap[discoveryv1.SchemeGroupVersion.Version] = storage

	return discoveryApiGroup, nil
}

// createCoreAPIGroupInfo creates the APIGroupInfo for the core API
func createCoreAPIGroupInfo(
	restOptionsGetter generic.RESTOptionsGetter,
	keycloakClient *keycloak.KeycloakClient,
	realmName string,
) (APIGroupInfo, error) {
	coreApiGroup := NewDefaultAPIGroup(core.GroupName, Scheme, metav1.ParameterCodec, Codecs)
	storage := map[string]rest.Storage{}

	// Storage for FormDefinitions
	formDefinitionsStorage, err := formdefinitionsstorage.NewREST(Scheme, restOptionsGetter)
	if err != nil {
		return APIGroupInfo{}, err
	}
	storage["formdefinitions"] = formDefinitionsStorage

	// Storage for CustomResourceDefinitions
	customResourceDefinitionsStorage, err := customresourcedefinitionstorage.NewREST(Scheme, restOptionsGetter)
	if err != nil {
		return APIGroupInfo{}, err
	}
	storage["customresourcedefinitions"] = customResourceDefinitionsStorage

	// Storage for OperatingScopes
	operatingScopeStorage, err := operatingscopestorage.NewREST(Scheme, restOptionsGetter)
	if err != nil {
		return APIGroupInfo{}, err
	}
	storage["operatingscopes"] = operatingScopeStorage

	userStorage := users.NewREST(keycloakClient, realmName)
	storage["users"] = userStorage

	// register the storage for the "core v1" version
	coreApiGroup.VersionedResourcesStorageMap[corev1.SchemeGroupVersion.Version] = storage

	return coreApiGroup, nil
}
