package extensions_apiserver

import (
  "fmt"
  "github.com/nrc-no/core/apps/api/pkg/apis/apiextensions"
  "github.com/nrc-no/core/apps/api/pkg/apis/apiextensions/install"
  v1 "github.com/nrc-no/core/apps/api/pkg/apis/apiextensions/v1"
  metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
  "github.com/nrc-no/core/apps/api/pkg/registry/generic"
  "github.com/nrc-no/core/apps/api/pkg/registry/rest"
  "github.com/nrc-no/core/apps/api/pkg/runtime"
  "github.com/nrc-no/core/apps/api/pkg/runtime/schema"
  "github.com/nrc-no/core/apps/api/pkg/runtime/serializer"
  "github.com/nrc-no/core/apps/api/pkg/server"
  "github.com/nrc-no/core/apps/api/pkg/server/extensions-apiserver/client/clientset"
  "github.com/nrc-no/core/apps/api/pkg/server/extensions-apiserver/client/informers"
  "k8s.io/apiextensions-apiserver/pkg/controller/apiapproval"
  "k8s.io/apiextensions-apiserver/pkg/controller/establish"
  "k8s.io/apiextensions-apiserver/pkg/controller/finalizer"
  "k8s.io/apiextensions-apiserver/pkg/controller/nonstructuralschema"
  "k8s.io/apiextensions-apiserver/pkg/controller/status"
  "k8s.io/apiextensions-apiserver/pkg/registry/customresourcedefinition"
  "k8s.io/apimachinery/pkg/util/wait"
  "k8s.io/apiserver/pkg/endpoints/discovery"
  "k8s.io/apiserver/pkg/util/webhook"
  "net/http"
  "time"
)

var (
  Scheme = runtime.NewScheme()
  Codecs = serializer.NewCodecFactory(Scheme)

  // if you modify this, make sure you update the crEncoder
  unversionedVersion = schema.GroupVersion{Group: "", Version: "v1"}
  unversionedTypes   = []runtime.Object{
    &metav1.Status{},
    &metav1.WatchEvent{},
    &metav1.APIVersions{},
    &metav1.APIGroupList{},
    &metav1.APIGroup{},
    &metav1.APIResourceList{},
  }
)

func init() {
  install.Install(Scheme)

  // we need to add the options to empty v1
  metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Group: "", Version: "v1"})

  Scheme.AddUnversionedTypes(unversionedVersion, unversionedTypes...)
}

type ExtraConfig struct {
  CRDRESTOptionsGetter generic.RESTOptionsGetter

  // MasterCount is used to detect whether cluster is HA, and if it is
  // the CRD Establishing will be hold by 5 seconds.
  MasterCount int

  // ServiceResolver is used in CR webhook converters to resolve webhook's service names
  ServiceResolver webhook.ServiceResolver
  // AuthResolverWrapper is used in CR webhook converters
  AuthResolverWrapper webhook.AuthenticationInfoResolverWrapper
}

type Config struct {
  GenericConfig *server.Config
  ExtraConfig   ExtraConfig
}

type completedConfig struct {
  GenericConfig server.CompletedConfig
  ExtraConfig   *ExtraConfig
}

type CompletedConfig struct {
  // Embed a private pointer that cannot be instantiated outside of this package.
  *completedConfig
}

type CustomResourceDefinitions struct {
  GenericAPIServer *server.Server

  // provided for easier embedding
  Informers informers.SharedInformerFactory
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
  c := completedConfig{
    cfg.GenericConfig.Complete(),
    &cfg.ExtraConfig,
  }

  //c.GenericConfig.EnableDiscovery = false
  //c.GenericConfig.Version = &version.Info{
  //  Major: "0",
  //  Minor: "1",
  //}

  return CompletedConfig{&c}
}

// New returns a new instance of CustomResourceDefinitions from the given config.
func (c completedConfig) New(delegationTarget server.DelegationTarget) (*CustomResourceDefinitions, error) {

  name := "apiextensions-apiserver"

  genericServer, err := c.GenericConfig.New(
    name,
    delegationTarget)
  if err != nil {
    return nil, err
  }

  s := &CustomResourceDefinitions{
    GenericAPIServer: genericServer,
  }

  // used later  to filter the served resource by those that have expired.
  //resourceExpirationEvaluator, err := server.NewResourceExpirationEvaluator(*c.GenericConfig.Version)
  //if err != nil {
  //  return nil, err
  //}

  apiResourceConfig := c.GenericConfig.MergedResourceConfig
  apiGroupInfo := server.NewDefaultAPIGroupInfo(apiextensions.GroupName, Scheme, metav1.ParameterCodec, Codecs)
  //if resourceExpirationEvaluator.ShouldServeForVersion(1, 22) && apiResourceConfig.VersionEnabled(v1beta1.SchemeGroupVersion) {
  //  storage := map[string]rest.Storage{}
  //  // customresourcedefinitions
  //  customResourceDefinitionStorage, err := customresourcedefinition.NewREST(Scheme, c.GenericConfig.RESTOptionsGetter)
  //  if err != nil {
  //    return nil, err
  //  }
  //  storage["customresourcedefinitions"] = customResourceDefinitionStorage
  //  storage["customresourcedefinitions/status"] = customresourcedefinition.NewStatusREST(Scheme, customResourceDefinitionStorage)
  //
  //  apiGroupInfo.VersionedResourcesStorageMap[v1beta1.SchemeGroupVersion.Version] = storage
  //}
  if apiResourceConfig.VersionEnabled(v1.SchemeGroupVersion) {
    storage := map[string]rest.Storage{}
    // customresourcedefinitions
    customResourceDefinitionStorage, err := customresourcedefinition.NewREST(Scheme, c.GenericConfig.RESTOptionsGetter)
    if err != nil {
      return nil, err
    }
    storage["customresourcedefinitions"] = customResourceDefinitionStorage
    storage["customresourcedefinitions/status"] = customresourcedefinition.NewStatusREST(Scheme, customResourceDefinitionStorage)

    apiGroupInfo.VersionedResourcesStorageMap[v1.SchemeGroupVersion.Version] = storage
  }

  if err := s.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
    return nil, err
  }

  crdClient, err := clientset.NewForConfig(s.GenericAPIServer.LoopbackClientConfig)
  if err != nil {
    // it's really bad that this is leaking here, but until we can fix the test (which I'm pretty sure isn't even testing what it wants to test),
    // we need to be able to move forward
    return nil, fmt.Errorf("failed to create clientset: %v", err)
  }
  s.Informers = informers.NewSharedInformerFactory(crdClient, 5*time.Minute)

  delegateHandler := delegationTarget.UnprotectedHandler()
  if delegateHandler == nil {
    delegateHandler = http.NotFoundHandler()
  }

  versionDiscoveryHandler := &versionDiscoveryHandler{
    discovery: map[schema.GroupVersion]*discovery.APIVersionHandler{},
    delegate:  delegateHandler,
  }
  groupDiscoveryHandler := &groupDiscoveryHandler{
    discovery: map[string]*discovery.APIGroupHandler{},
    delegate:  delegateHandler,
  }
  establishingController := establish.NewEstablishingController(s.Informers.Apiextensions().V1().CustomResourceDefinitions(), crdClient.ApiextensionsV1())
  crdHandler, err := NewCustomResourceDefinitionHandler(
    versionDiscoveryHandler,
    groupDiscoveryHandler,
    s.Informers.Apiextensions().V1().CustomResourceDefinitions(),
    delegateHandler,
    c.ExtraConfig.CRDRESTOptionsGetter,
    c.GenericConfig.AdmissionControl,
    establishingController,
    c.ExtraConfig.ServiceResolver,
    c.ExtraConfig.AuthResolverWrapper,
    c.ExtraConfig.MasterCount,
    s.GenericAPIServer.Authorizer,
    c.GenericConfig.RequestTimeout,
    time.Duration(c.GenericConfig.MinRequestTimeout)*time.Second,
    apiGroupInfo.StaticOpenAPISpec,
    c.GenericConfig.MaxRequestBodyBytes,
  )
  if err != nil {
    return nil, err
  }
  s.GenericAPIServer.Handler.NonGoRestfulMux.Handle("/apis", crdHandler)
  s.GenericAPIServer.Handler.NonGoRestfulMux.HandlePrefix("/apis/", crdHandler)

  discoveryController := NewDiscoveryController(s.Informers.Apiextensions().V1().CustomResourceDefinitions(), versionDiscoveryHandler, groupDiscoveryHandler)
  namingController := status.NewNamingConditionController(s.Informers.Apiextensions().V1().CustomResourceDefinitions(), crdClient.ApiextensionsV1())
  nonStructuralSchemaController := nonstructuralschema.NewConditionController(s.Informers.Apiextensions().V1().CustomResourceDefinitions(), crdClient.ApiextensionsV1())
  apiApprovalController := apiapproval.NewKubernetesAPIApprovalPolicyConformantConditionController(s.Informers.Apiextensions().V1().CustomResourceDefinitions(), crdClient.ApiextensionsV1())
  finalizingController := finalizer.NewCRDFinalizer(
    s.Informers.Apiextensions().V1().CustomResourceDefinitions(),
    crdClient.ApiextensionsV1(),
    crdHandler,
  )
  openapiController := openapicontroller.NewController(s.Informers.Apiextensions().V1().CustomResourceDefinitions())

  s.GenericAPIServer.AddPostStartHookOrDie("start-apiextensions-informers", func(context server.PostStartHookContext) error {
    s.Informers.Start(context.StopCh)
    return nil
  })
  s.GenericAPIServer.AddPostStartHookOrDie("start-apiextensions-controllers", func(context server.PostStartHookContext) error {
    // OpenAPIVersionedService and StaticOpenAPISpec are populated in generic apiserver PrepareRun().
    // Together they serve the /openapi/v2 endpoint on a generic apiserver. A generic apiserver may
    // choose to not enable OpenAPI by having null openAPIConfig, and thus OpenAPIVersionedService
    // and StaticOpenAPISpec are both null. In that case we don't run the CRD OpenAPI controller.
    if s.GenericAPIServer.OpenAPIVersionedService != nil && s.GenericAPIServer.StaticOpenAPISpec != nil {
      go openapiController.Run(s.GenericAPIServer.StaticOpenAPISpec, s.GenericAPIServer.OpenAPIVersionedService, context.StopCh)
    }

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
  // we don't want to report healthy until we can handle all CRDs that have already been registered.  Waiting for the informer
  // to sync makes sure that the lister will be valid before we begin.  There may still be races for CRDs added after startup,
  // but we won't go healthy until we can handle the ones already present.
  s.GenericAPIServer.AddPostStartHookOrDie("crd-informer-synced", func(context server.PostStartHookContext) error {
    return wait.PollImmediateUntil(100*time.Millisecond, func() (bool, error) {
      return s.Informers.Apiextensions().V1().CustomResourceDefinitions().Informer().HasSynced(), nil
    }, context.StopCh)
  })

  return s, nil
}
