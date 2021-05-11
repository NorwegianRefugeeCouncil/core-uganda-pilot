package server

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restclient "github.com/nrc-no/core/apps/api/pkg/client/rest"
	"github.com/nrc-no/core/apps/api/pkg/endpoints"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/runtime/serializer"
	"github.com/nrc-no/core/apps/api/pkg/storageversion"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
	"net/http"
	"strings"
)

// Info about an API group.
type APIGroupInfo struct {
	PrioritizedVersions []schema.GroupVersion
	// Info about the resources in this group. It's a map from version to resource to the storage.
	VersionedResourcesStorageMap map[string]map[string]rest.Storage
	// OptionsExternalVersion controls the APIVersion used for common objects in the
	// schema like api.Status, api.DeleteOptions, and metav1.ListOptions. Other implementors may
	// define a version "v1beta1" but want to use the Kubernetes "v1" internal objects.
	// If nil, defaults to groupMeta.GroupVersion.
	// TODO: Remove this when https://github.com/kubernetes/kubernetes/issues/19018 is fixed.
	OptionsExternalVersion *schema.GroupVersion
	// MetaGroupVersion defaults to "meta.k8s.io/v1" and is the scheme group version used to decode
	// common API implementations like ListOptions. Future changes will allow this to vary by group
	// version (for when the inevitable meta/v2 group emerges).
	MetaGroupVersion *schema.GroupVersion

	// Scheme includes all of the types used by this group and how to convert between them (or
	// to convert objects from outside of this group that are accepted in this API).
	// TODO: replace with interfaces
	Scheme *runtime.Scheme
	// NegotiatedSerializer controls how this group encodes and decodes data
	NegotiatedSerializer runtime.NegotiatedSerializer
	// ParameterCodec performs conversions for query parameters passed to API calls
	ParameterCodec runtime.ParameterCodec

	// StaticOpenAPISpec is the spec derived from the definitions of all resources installed together.
	// It is set during InstallAPIGroups, InstallAPIGroup, and InstallLegacyAPIGroup.
	// StaticOpenAPISpec *spec.Swagger
}

type Server struct {
	Container                  *restful.Container
	handlerChain               http.Handler
	LoopbackClientConfig       *restclient.Config
	Serializer                 runtime.NegotiatedSerializer
	delegationTarget           DelegationTarget
	EquivalentResourceRegistry runtime.EquivalentResourceRegistry
	Handler                    *APIServerHandler
	postStartHooks             map[string]postStartHookEntry
	preShutdownHooks           map[string]preShutdownHookEntry
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.handlerChain.ServeHTTP(w, req)
}

type HandlerChainBuilderFn func(apiHandler http.Handler) http.Handler

func CreateServerConfig() {

}

type APIAggregator struct {
	apiServer       *Server
	delegateHandler http.Handler
	handledGroups   sets.String
}

// PostStartHookContext provides information about this API server to a PostStartHookFunc
type PostStartHookContext struct {
	// LoopbackClientConfig is a config for a privileged loopback connection to the API server
	LoopbackClientConfig *restclient.Config
	// StopCh is the channel that will be closed when the server stops
	StopCh <-chan struct{}
}

// PreShutdownHookFunc is a function that can be added to the shutdown logic.
type PreShutdownHookFunc func() error

// PostStartHookFunc is a function that is called after the server has started.
// It must properly handle cases like:
//  1. asynchronous start in multiple API server processes
//  2. conflicts between the different processes all trying to perform the same action
//  3. partially complete work (API server crashes while running your hook)
//  4. API server access **BEFORE** your hook has completed
// Think of it like a mini-controller that is super privileged and gets to run in-process
// If you use this feature, tag @deads2k on github who has promised to review code for anyone's PostStartHook
// until it becomes easier to use.
type PostStartHookFunc func(context PostStartHookContext) error

type postStartHookEntry struct {
	hook PostStartHookFunc
	// originatingStack holds the stack that registered postStartHooks. This allows us to show a more helpful message
	// for duplicate registration.
	originatingStack string

	// done will be closed when the postHook is finished
	done chan struct{}
}

type preShutdownHookEntry struct {
	hook PreShutdownHookFunc
}

type preparedGenericAPIServer struct {
	*Server
}

// DelegationTarget is an interface which allows for composition of API servers with top level handling that works
// as expected.
type DelegationTarget interface {
	// UnprotectedHandler returns a handler that is NOT protected by a normal chain
	UnprotectedHandler() http.Handler

	// PostStartHooks returns the post-start hooks that need to be combined
	PostStartHooks() map[string]postStartHookEntry

	// PreShutdownHooks returns the pre-stop hooks that need to be combined
	PreShutdownHooks() map[string]preShutdownHookEntry

	// HealthzChecks returns the healthz checks that need to be combined
	// HealthzChecks() []healthz.HealthChecker

	// ListedPaths returns the paths for supporting an index
	ListedPaths() []string

	// NextDelegate returns the next delegationTarget in the chain of delegations
	NextDelegate() DelegationTarget

	// PrepareRun does post API installation setup steps. It calls recursively the same function of the delegates.
	PrepareRun() preparedGenericAPIServer
}

// NewDefaultAPIGroupInfo returns an APIGroupInfo stubbed with "normal" values
// exposed for easier composition from other packages
func NewDefaultAPIGroupInfo(group string, scheme *runtime.Scheme, parameterCodec runtime.ParameterCodec, codecs serializer.CodecFactory) APIGroupInfo {
	return APIGroupInfo{
		PrioritizedVersions:          scheme.PrioritizedVersionsForGroup(group),
		VersionedResourcesStorageMap: map[string]map[string]rest.Storage{},
		// TODO unhardcode this.  It was hardcoded before, but we need to re-evaluate
		OptionsExternalVersion: &schema.GroupVersion{Version: "v1"},
		Scheme:                 scheme,
		ParameterCodec:         parameterCodec,
		NegotiatedSerializer:   codecs,
	}
}

func (s *Server) InstallAPIGroups(apiGroupInfos ...*APIGroupInfo) error {

	for _, apiGroupInfo := range apiGroupInfos {
		if len(apiGroupInfo.PrioritizedVersions[0].Group) == 0 {
			return fmt.Errorf("cannot register handler with an empty group for %#v", *apiGroupInfo)
		}
		if len(apiGroupInfo.PrioritizedVersions[0].Version) == 0 {
			return fmt.Errorf("cannot register handler with an empty version for %#v", *apiGroupInfo)
		}
	}

	for _, apiGroupInfo := range apiGroupInfos {
		if err := s.installAPIResources("/apis", apiGroupInfo); err != nil {
			return fmt.Errorf("unable to install api resources; %v", err)
		}
	}

	return nil

}

func (s *Server) installAPIResources(apiPrefix string, apiGroupInfo *APIGroupInfo) error {

	var resourceInfos []*storageversion.ResourceInfo
	for _, groupVersion := range apiGroupInfo.PrioritizedVersions {

		if len(apiGroupInfo.VersionedResourcesStorageMap[groupVersion.Version]) == 0 {
			logrus.Warnf("skipping api %s because it has no resources", groupVersion)
			continue
		}

		apiGroupVersion := s.getApiGroupVersion(apiGroupInfo, groupVersion, apiPrefix)
		// apiGroupVersion.MaxRequestBodyBytes =

		r, err := apiGroupVersion.InstallREST(s.Handler.GoRestfulContainer)
		if err != nil {
			return fmt.Errorf("unable to setup api %v: %v", apiGroupInfo, err)
		}
		resourceInfos = append(resourceInfos, r...)

	}

	return nil

}

func (s *Server) getApiGroupVersion(apiGroupInfo *APIGroupInfo, groupVersion schema.GroupVersion, prefix string) *endpoints.APIGroupVersion {
	storage := make(map[string]rest.Storage)
	for k, v := range apiGroupInfo.VersionedResourcesStorageMap[groupVersion.Version] {
		storage[strings.ToLower(k)] = v
	}
	version := s.newAPIGroupVersion(apiGroupInfo, groupVersion)
	version.Root = "/apis"
	version.Storage = storage
	return version
}

func (s *Server) newAPIGroupVersion(apiGroupInfo *APIGroupInfo, groupVersion schema.GroupVersion) *endpoints.APIGroupVersion {
	return &endpoints.APIGroupVersion{
		GroupVersion:     groupVersion,
		MetaGroupVersion: apiGroupInfo.MetaGroupVersion,

		ParameterCodec: apiGroupInfo.ParameterCodec,
		Serializer:     apiGroupInfo.NegotiatedSerializer,
		Creater:        apiGroupInfo.Scheme,
		Convertor:      apiGroupInfo.Scheme,
		//ConvertabilityChecker: apiGroupInfo.Scheme,
		UnsafeConvertor: runtime.UnsafeObjectConvertor(apiGroupInfo.Scheme),
		//Defaulter:             apiGroupInfo.Scheme,
		Typer: apiGroupInfo.Scheme,
		//Linker:                runtime.SelfLinker(meta.NewAccessor()),

		EquivalentResourceRegistry: s.EquivalentResourceRegistry,

		//Admit:             s.admissionControl,
		//MinRequestTimeout: s.minRequestTimeout,
		//Authorizer:        s.Authorizer,
	}
}

type emptyDelegate struct {
}

func NewEmptyDelegate() DelegationTarget {
	return emptyDelegate{}
}

func (s emptyDelegate) UnprotectedHandler() http.Handler {
	return nil
}
func (s emptyDelegate) PostStartHooks() map[string]postStartHookEntry {
	return map[string]postStartHookEntry{}
}
func (s emptyDelegate) PreShutdownHooks() map[string]preShutdownHookEntry {
	return map[string]preShutdownHookEntry{}
}

//func (s emptyDelegate) HealthzChecks() []healthz.HealthChecker {
//  return []healthz.HealthChecker{}
//}
func (s emptyDelegate) ListedPaths() []string {
	return []string{}
}
func (s emptyDelegate) NextDelegate() DelegationTarget {
	return nil
}
func (s emptyDelegate) PrepareRun() preparedGenericAPIServer {
	return preparedGenericAPIServer{nil}
}
