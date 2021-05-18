package server

import (
	"github.com/nrc-no/core/api/pkg/apis/core"
	"github.com/nrc-no/core/api/pkg/apis/core/install"
	v1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/api/pkg/generated/openapi"
	"github.com/nrc-no/core/api/pkg/registry/core/formdefinition"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/endpoints/discovery"
	"k8s.io/apiserver/pkg/endpoints/filters"
	openapi2 "k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/storageversion"
	"k8s.io/kube-openapi/pkg/common"
	"net/http"
	"strings"
	"time"
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	install.Install(Scheme)
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
	apiServerHandler := server.NewAPIServerHandler("", c.Serializer, handlerChainBuilder, http.NotFoundHandler())

	s := &APIServer{
		Serializer:                 c.Serializer,
		DiscoveryGroupManager:      discovery.NewRootAPIsHandler(discovery.DefaultAddresses{DefaultAddress: c.ExternalAddress}, c.Serializer),
		EquivalentResourceRegistry: runtime.NewEquivalentResourceRegistry(),
		minRequestTimeout:          time.Duration(1800) * time.Second,
		maxRequestBodyBytes:        int64(3 * 1024 * 1024),
		StorageVersionManager:      storageversion.NewDefaultManager(),
		Handler:                    apiServerHandler,
		RESTOptionsGetter:          c.RESTOptionsGetter,
	}

	s.openAPIConfig = server.DefaultOpenAPIConfig(openapi.GetOpenAPIDefinitions, openapi2.NewDefinitionNamer(Scheme))
	s.openAPIConfig.Info.Title = "Core API"
	s.openAPIConfig.Info.Version = "0.1"

	apiResourceConfig := server.NewDefaultAPIGroupInfo(core.GroupName, Scheme, metav1.ParameterCodec, Codecs)
	storage := map[string]rest.Storage{}
	formDefinitionsStorage, err := formdefinition.NewREST(Scheme, s.RESTOptionsGetter)
	if err != nil {
		return nil, err
	}
	storage["formdefinitions"] = formDefinitionsStorage
	apiResourceConfig.VersionedResourcesStorageMap[v1.SchemeGroupVersion.Version] = storage

	if err := s.InstallAPIGroups(&apiResourceConfig); err != nil {
		return nil, err
	}

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
