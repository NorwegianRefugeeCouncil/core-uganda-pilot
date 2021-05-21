package server2

import (
	"context"
	"github.com/nrc-no/core/api/pkg/apis/core"
	"github.com/nrc-no/core/api/pkg/apis/core/install"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	v1 "github.com/nrc-no/core/api/pkg/client/core/v1"
	restclient "github.com/nrc-no/core/api/pkg/client/rest"
	"github.com/nrc-no/core/api/pkg/server2/endpoints/handlers"
	"github.com/nrc-no/core/api/pkg/server2/registry/core/customresourcedefinition"
	"github.com/nrc-no/core/api/pkg/server2/registry/core/formdefinitions"
	"github.com/nrc-no/core/api/pkg/server2/registry/generic"
	"github.com/nrc-no/core/api/pkg/server2/registry/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"net"
	"net/http"
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	install.Install(Scheme)
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

// New creates a new Server from the CompletedConfig
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
	apiServerHandler := NewAPIServerHandler(handlerChainBuilder)

	// Installs the known resource handlers in the API
	// These include FormDefinitions and CustomResourceDefinitions
	if err := c.installApiGroups(apiServerHandler); err != nil {
		return nil, err
	}

	v1Client, err := v1.NewForConfig(c.LoopbackClientConfig)
	if err != nil {
		return nil, err
	}

	// Installs the CustomResource handler in the API
	// This is ran after go-restful tries to match the route
	// The CustomResource handler is able to serve new resources
	// dynamically (eg. created at runtime)
	if err := c.installCustomResources(apiServerHandler, v1Client.CustomResourceDefinitions()); err != nil {
		return nil, err
	}

	// Create the API server
	server := &Server{
		ctx:      ctx,
		handler:  apiServerHandler,
		listener: c.Listener,
	}

	return server, nil
}

// installCustomResources installs the required handlers to serve dynamic
// CustomResources from the API
func (c *CompletedConfig) installCustomResources(apiServerHandler *APIServerHandler, cli v1.CustomResourceDefinitionInterface) error {
	crdHandler, err := handlers.NewCustomResourceDefinitionHandler(
		c.CRDRestOptionsGetter,
		Codecs,
		cli)
	if err != nil {
		return err
	}
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
func createCoreAPIGroupInfo(restOptionsGetter generic.RESTOptionsGetter) (APIGroupInfo, error) {
	coreApiGroup := NewDefaultAPIGroup(core.GroupName, Scheme, metav1.ParameterCodec, Codecs)
	storage := map[string]rest.Storage{}

	// Storage for FormDefinitions
	formDefinitionsStorage, err := formdefinitions.NewREST(Scheme, restOptionsGetter)
	if err != nil {
		return APIGroupInfo{}, err
	}
	storage["formdefinitions"] = formDefinitionsStorage

	// Storage for CustomResourceDefinitions
	customResourceDefinitionsStorage, err := customresourcedefinition.NewREST(Scheme, restOptionsGetter)
	if err != nil {
		return APIGroupInfo{}, err
	}
	storage["customresourcedefinitions"] = customResourceDefinitionsStorage

	// register the storage for the "core v1" version
	coreApiGroup.VersionedResourcesStorageMap[corev1.SchemeGroupVersion.Version] = storage

	return coreApiGroup, nil
}
