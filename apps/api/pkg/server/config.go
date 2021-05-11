package server

import (
	"fmt"
	"github.com/nrc-no/core/apps/api/pkg/client/rest"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/filters"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/request"
	"github.com/nrc-no/core/apps/api/pkg/registry/generic"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/runtime/serializer"
	"github.com/nrc-no/core/apps/api/pkg/server/storage"
	"net/http"
)

const (
	// DefaultLegacyAPIPrefix is where the legacy APIs will be located.
	DefaultLegacyAPIPrefix = "/api"

	// APIGroupPrefix is where non-legacy API group will be located.
	APIGroupPrefix = "/apis"
)

type Config struct {
	Serializer                 runtime.NegotiatedSerializer
	RESTOptionsGetter          generic.RESTOptionsGetter
	MergedResourceConfig       *storage.ResourceConfig
	LoopbackClientConfig       *rest.Config
	EquivalentResourceRegistry runtime.EquivalentResourceRegistry
	BuildHandlerChainFunc      func(apiHandler http.Handler, c *Config) http.Handler
	RequestInfoResolver        request.RequestInfoResolver
}

func (c *Config) Complete() CompletedConfig {

	if c.EquivalentResourceRegistry == nil {
		if c.RESTOptionsGetter == nil {
			c.EquivalentResourceRegistry = runtime.NewEquivalentResourceRegistry()
		} else {
			c.EquivalentResourceRegistry = runtime.NewEquivalentResourceRegistryWithIdentity(func(gr schema.GroupResource) string {
				if opts, err := c.RESTOptionsGetter.GetRESTOptions(gr); err == nil {
					return opts.ResourcePrefix
				}
				return ""
			})
		}
	}

	if c.RequestInfoResolver == nil {
		c.RequestInfoResolver = NewRequestInfoResolver(c)
	}

	return CompletedConfig{&completedConfig{c}}
}

func NewConfig(codecs serializer.CodecFactory) *Config {
	return &Config{
		Serializer:            codecs,
		BuildHandlerChainFunc: DefaultBuildHandlerChain,
	}
}

type completedConfig struct {
	*Config
}

type CompletedConfig struct {
	*completedConfig
}

func (c CompletedConfig) New(delegationTarget DelegationTarget) (*Server, error) {

	if c.Serializer == nil {
		return nil, fmt.Errorf("Server.Serializer == nil")
	}
	if c.LoopbackClientConfig == nil {
		return nil, fmt.Errorf("Server.LoopbackClientConfig == nil")
	}
	if c.EquivalentResourceRegistry == nil {
		return nil, fmt.Errorf("Server.EquivalentResourceRegistry == nil")
	}

	handlerChainBuilder := func(handler http.Handler) http.Handler {
		return c.BuildHandlerChainFunc(handler, c.Config)
	}
	apiServerHandler := NewAPIServerHandler(c.Serializer, handlerChainBuilder, delegationTarget.UnprotectedHandler())

	s := &Server{
		LoopbackClientConfig:       c.LoopbackClientConfig,
		Serializer:                 c.Serializer,
		delegationTarget:           delegationTarget,
		EquivalentResourceRegistry: c.EquivalentResourceRegistry,
		Handler:                    apiServerHandler,
		postStartHooks:             map[string]postStartHookEntry{},
		preShutdownHooks:           map[string]preShutdownHookEntry{},
	}

	for k, v := range delegationTarget.PostStartHooks() {
		s.postStartHooks[k] = v
	}
	for k, v := range delegationTarget.PreShutdownHooks() {
		s.preShutdownHooks[k] = v
	}

	return s, nil
}

func DefaultBuildHandlerChain(apiHandler http.Handler, c *Config) http.Handler {
	handler := filters.WithRequestInfo(apiHandler, c.RequestInfoResolver)
	return handler
}

func NewRequestInfoResolver(c *Config) *request.RequestInfoFactory {
	//apiPrefixes := sets.NewString(strings.Trim(APIGroupPrefix, "/")) // all possible API prefixes
	//legacyAPIPrefixes := sets.String{}                               // APIPrefixes that won't have groups (legacy)
	//for legacyAPIPrefix := range c.LegacyAPIGroupPrefixes {
	//  apiPrefixes.Insert(strings.Trim(legacyAPIPrefix, "/"))
	//  legacyAPIPrefixes.Insert(strings.Trim(legacyAPIPrefix, "/"))
	//}

	return &request.RequestInfoFactory{
		//
		//APIPrefixes:          apiPrefixes,
		//GrouplessAPIPrefixes: legacyAPIPrefixes,
	}
}
