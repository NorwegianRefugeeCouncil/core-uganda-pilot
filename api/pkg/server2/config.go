package server2

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/nrc-no/core/api/pkg/apis/core"
	"github.com/nrc-no/core/api/pkg/apis/core/install"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/api/pkg/server2/registry/core/formdefinitions"
	"github.com/nrc-no/core/api/pkg/server2/registry/generic"
	"github.com/nrc-no/core/api/pkg/server2/registry/rest"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
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
	ListenAddress         string
	RESTOptionsGetter     generic.RESTOptionsGetter
	BuildHandlerChainFunc func(handler http.Handler, config *Config) http.Handler
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

func (c *CompletedConfig) New() (*Server, error) {
	handlerChainBuilder := func(handler http.Handler) http.Handler {
		return c.BuildHandlerChainFunc(handler, c.Config)
	}
	apiServerHandler := NewAPIServerHandler(handlerChainBuilder)

	var apiGroups []*APIGroupInfo
	coreApiGroup, err := createCoreAPIGroupInfo(c.RESTOptionsGetter)
	if err != nil {
		return nil, err
	}
	apiGroups = append(apiGroups, &coreApiGroup)

	if err := installApiGroups(apiServerHandler.GoRestfulContainer, "/apis", apiGroups...); err != nil {
		return nil, err
	}

	server := &Server{
		listenAddress: c.ListenAddress,
		handler:       apiServerHandler,
	}
	return server, nil
}

// installApiGroups registers the API groups into go-restful container
// this method will register the necessary routes and handlers
func installApiGroups(goRestfulContainer *restful.Container, apiPrefix string, apiGroupInfos ...*APIGroupInfo) error {

	for _, apiGroupInfo := range apiGroupInfos {
		if len(apiGroupInfo.PrioritizedVersions[0].Group) == 0 {
			return fmt.Errorf("cannot register handler with an empty group for %#v", *apiGroupInfo)
		}
		if len(apiGroupInfo.PrioritizedVersions[0].Version) == 0 {
			return fmt.Errorf("cannot register handler with an empty version for %#v", *apiGroupInfo)
		}
	}

	for _, apiGroupInfo := range apiGroupInfos {
		if err := installApiResources(goRestfulContainer, apiPrefix, apiGroupInfo); err != nil {
			return err
		}
	}

	return nil
}

func installApiResources(goRestfulContainer *restful.Container, apiPrefix string, apiGroupInfo *APIGroupInfo) error {
	for _, groupVersion := range apiGroupInfo.PrioritizedVersions {

		if len(apiGroupInfo.VersionedResourcesStorageMap[groupVersion.Version]) == 0 {
			logrus.Warnf("skipping api %v because it has no resources", groupVersion)
			continue
		}

		apiGroupVersion := apiGroupInfo.GetAPIGroupVersion(groupVersion, apiPrefix)
		if err := apiGroupVersion.InstallREST(goRestfulContainer); err != nil {
			return err
		}

	}
	return nil
}

// installApiGroup registers an API group into the go-restful container.
// see installApiGroups
func installApiGroup(goRestfulContainer *restful.Container, apiPrefix string, apiGroupInfo *APIGroupInfo) error {
	return installApiGroups(goRestfulContainer, apiPrefix, apiGroupInfo)
}

// createCoreAPIGroupInfo creates the APIGroupInfo for the core API
func createCoreAPIGroupInfo(restOptionsGetter generic.RESTOptionsGetter) (APIGroupInfo, error) {
	coreApiGroup := NewDefaultAPIGroup(core.GroupName, Scheme, metav1.ParameterCodec, Codecs)
	storage := map[string]rest.Storage{}
	formDefinitionsStorage, err := formdefinitions.NewREST(Scheme, restOptionsGetter)
	if err != nil {
		return APIGroupInfo{}, err
	}
	storage["formdefinitions"] = formDefinitionsStorage
	coreApiGroup.VersionedResourcesStorageMap[corev1.SchemeGroupVersion.Version] = storage
	return coreApiGroup, nil
}
