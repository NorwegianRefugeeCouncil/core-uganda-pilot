package server

import (
  "fmt"
  "github.com/go-openapi/spec"
  restclient "github.com/nrc-no/core/apps/api/pkg/client/rest"
  "github.com/nrc-no/core/apps/api/pkg/endpoints"
  "github.com/nrc-no/core/apps/api/pkg/registry/rest"
  "github.com/nrc-no/core/apps/api/pkg/runtime"
  "github.com/nrc-no/core/apps/api/pkg/runtime/schema"
  "github.com/nrc-no/core/apps/api/pkg/runtime/serializer"
  "github.com/nrc-no/core/apps/api/pkg/storageversion"
  "github.com/sirupsen/logrus"
  "k8s.io/apimachinery/pkg/util/clock"
  "k8s.io/apimachinery/pkg/util/sets"
  "k8s.io/apimachinery/pkg/util/waitgroup"
  "k8s.io/apiserver/pkg/admission"
  "k8s.io/apiserver/pkg/audit"
  "k8s.io/apiserver/pkg/authorization/authorizer"
  "k8s.io/apiserver/pkg/endpoints/discovery"
  "k8s.io/apiserver/pkg/server/healthz"
  "k8s.io/apiserver/pkg/server/routes"
  openapicommon "k8s.io/kube-openapi/pkg/common"
  "k8s.io/kube-openapi/pkg/handler"
  "net/http"
  "strings"
  "sync"
  "time"
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
  StaticOpenAPISpec *spec.Swagger
}

type Server struct {
  // discoveryAddresses is used to build cluster IPs for discovery.
  discoveryAddresses discovery.Addresses

  // LoopbackClientConfig is a config for a privileged loopback connection to the API server
  LoopbackClientConfig *restclient.Config

  // minRequestTimeout is how short the request timeout can be.  This is used to build the RESTHandler
  minRequestTimeout time.Duration

  // ShutdownTimeout is the timeout used for server shutdown. This specifies the timeout before server
  // gracefully shutdown returns.
  ShutdownTimeout time.Duration

  // legacyAPIGroupPrefixes is used to set up URL parsing for authorization and for validating requests
  // to InstallLegacyAPIGroup
  legacyAPIGroupPrefixes sets.String

  // admissionControl is used to build the RESTStorage that backs an API Group.
  admissionControl admission.Interface

  // SecureServingInfo holds configuration of the TLS server.
  SecureServingInfo *SecureServingInfo

  // ExternalAddress is the address (hostname or IP and port) that should be used in
  // external (public internet) URLs for this GenericAPIServer.
  ExternalAddress string

  // Serializer controls how common API objects not in a group/version prefix are serialized for this server.
  // Individual APIGroups may define their own serializers.
  Serializer runtime.NegotiatedSerializer

  // "Outputs"
  // Handler holds the handlers being used by this API server
  Handler *APIServerHandler

  // listedPathProvider is a lister which provides the set of paths to show at /
  listedPathProvider routes.ListedPathProvider

  // DiscoveryGroupManager serves /apis
  DiscoveryGroupManager discovery.GroupManager

  // Enable swagger and/or OpenAPI if these configs are non-nil.
  openAPIConfig *openapicommon.Config

  // SkipOpenAPIInstallation indicates not to install the OpenAPI handler
  // during PrepareRun.
  // Set this to true when the specific API Server has its own OpenAPI handler
  // (e.g. kube-aggregator)
  skipOpenAPIInstallation bool

  // OpenAPIVersionedService controls the /openapi/v2 endpoint, and can be used to update the served spec.
  // It is set during PrepareRun if `openAPIConfig` is non-nil unless `skipOpenAPIInstallation` is true.
  OpenAPIVersionedService *handler.OpenAPIService

  // StaticOpenAPISpec is the spec derived from the restful container endpoints.
  // It is set during PrepareRun.
  StaticOpenAPISpec *spec.Swagger

  // PostStartHooks are each called after the server has started listening, in a separate go func for each
  // with no guarantee of ordering between them.  The map key is a name used for error reporting.
  // It may kill the process with a panic if it wishes to by returning an error.
  postStartHookLock      sync.Mutex
  postStartHooks         map[string]postStartHookEntry
  postStartHooksCalled   bool
  disabledPostStartHooks sets.String

  preShutdownHookLock    sync.Mutex
  preShutdownHooks       map[string]preShutdownHookEntry
  preShutdownHooksCalled bool

  // healthz checks
  healthzLock            sync.Mutex
  healthzChecks          []healthz.HealthChecker
  healthzChecksInstalled bool
  // livez checks
  livezLock            sync.Mutex
  livezChecks          []healthz.HealthChecker
  livezChecksInstalled bool
  // readyz checks
  readyzLock            sync.Mutex
  readyzChecks          []healthz.HealthChecker
  readyzChecksInstalled bool
  livezGracePeriod      time.Duration
  livezClock            clock.Clock
  // the readiness stop channel is used to signal that the apiserver has initiated a shutdown sequence, this
  // will cause readyz to return unhealthy.
  readinessStopCh chan struct{}

  // auditing. The backend is started after the server starts listening.
  AuditBackend audit.Backend

  // Authorizer determines whether a user is allowed to make a certain request. The Handler does a preliminary
  // authorization check using the request URI but it may be necessary to make additional checks, such as in
  // the create-on-update case
  Authorizer authorizer.Authorizer

  // EquivalentResourceRegistry provides information about resources equivalent to a given resource,
  // and the kind associated with a given resource. As resources are installed, they are registered here.
  EquivalentResourceRegistry runtime.EquivalentResourceRegistry

  // delegationTarget is the next delegate in the chain. This is never nil.
  delegationTarget DelegationTarget

  // HandlerChainWaitGroup allows you to wait for all chain handlers finish after the server shutdown.
  HandlerChainWaitGroup *waitgroup.SafeWaitGroup

  // ShutdownDelayDuration allows to block shutdown for some time, e.g. until endpoints pointing to this API server
  // have converged on all node. During this time, the API server keeps serving, /healthz will return 200,
  // but /readyz will return failure.
  ShutdownDelayDuration time.Duration

  // The limit on the request body size that would be accepted and decoded in a write request.
  // 0 means no limit.
  maxRequestBodyBytes int64

  // APIServerID is the ID of this API server
  APIServerID string

  // StorageVersionManager holds the storage versions of the API resources installed by this server.
  StorageVersionManager storageversion.Manager
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

// Exposes the given api group in the API.
func (s *Server) InstallAPIGroup(apiGroupInfo *APIGroupInfo) error {
  return s.InstallAPIGroups(apiGroupInfo)
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
