package endpoints

import (
  "fmt"
  "github.com/emicklei/go-restful"
  "github.com/nrc-no/core/apps/api/apis/meta"
  "github.com/nrc-no/core/apps/api/pkg/endpoints/handlers"
  "github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/negotiation"
  "github.com/nrc-no/core/apps/api/pkg/registry/rest"
  "k8s.io/apimachinery/pkg/runtime"
  "k8s.io/apimachinery/pkg/runtime/schema"
  "net/http"
  "reflect"
  "sort"
  "strings"
)

type APIInstaller struct {
  prefix string
  group  *APIGroupVersion
}

func (a *APIInstaller) Install(container *restful.Container) ([]meta.APIResource, *restful.WebService, []error) {
  var apiResources []meta.APIResource
  var errors []error

  var ws *restful.WebService
  var shouldAddWebService = false
  for _, service := range container.RegisteredWebServices() {
    if service.RootPath() == a.prefix {
      ws = service
      break
    }
  }
  if ws == nil {
    shouldAddWebService = true
    ws = a.newWebService()
  }

  paths := make([]string, len(a.group.Storage))
  var i = 0
  for path := range a.group.Storage {
    paths[i] = path
    i++
  }
  sort.Strings(paths)

  for _, path := range paths {
    apiResource, err := a.registerResourceHandlers(path, a.group.Storage[path], ws)
    if err != nil {
      errors = append(errors, fmt.Errorf("error in registering resource [%v]: %v", path, err))
    }
    if apiResource != nil {
      apiResources = append(apiResources, *apiResource)
    }
  }

  if shouldAddWebService {
    container.Add(ws)
  }

  return apiResources, ws, errors
}

func (a *APIInstaller) newWebService() *restful.WebService {
  ws := new(restful.WebService)
  ws.Path(a.prefix)
  ws.Doc("API at " + a.prefix)
  ws.Consumes("application/json")
  mediaTypes, streamMediaTypes := negotiation.MediaTypesForSerializer(a.group.Serializer)
  ws.Produces(append(mediaTypes, streamMediaTypes...)...)
  ws.ApiVersion(a.group.GroupVersion.String())

  return ws
}

func (a *APIInstaller) registerResourceHandlers(path string, storage rest.Storage, ws *restful.WebService) (*meta.APIResource, error) {
  var resource string
  switch parts := strings.Split(path, "/"); len(parts) {
  case 1:
    resource = parts[0]
  default:
    return nil, fmt.Errorf("invalid resource")
  }

  // group, version := a.group.GroupVersion.Group, a.group.GroupVersion.Version
  fqKindToRegister, err := GetResourceKind(a.group.GroupVersion, storage, a.group.Typer)
  if err != nil {
    return nil, err
  }

  versionedPtr, err := a.group.Creater.New(fqKindToRegister)
  if err != nil {
    return nil, err
  }

  defaultVersionedObject := indirectArbitraryPointer(versionedPtr)
  kind := fqKindToRegister.Kind

  creater, isCreater := storage.(rest.Creater)
  // lister, isLister := storage.(rest.Lister)
  getter, isGetter := storage.(rest.Getter)
  // updater, isUpdater := storage.(rest.Updater)

  var apiResource meta.APIResource

  apiResource.Name = path
  apiResource.Kind = kind

  reqScope := handlers.RequestScope{
    Namer:          nil,
    Serializer:     a.group.Serializer,
    ParameterCoder: a.group.ParameterCodec,
    Creater:        a.group.Creater,
    Resource:       a.group.GroupVersion.WithResource(resource),
    Kind:           fqKindToRegister,
  }

  mediaTypes, _ := negotiation.MediaTypesForSerializer(a.group.Serializer)

  routes := []*restful.RouteBuilder{}

  if isGetter {
    getRoute := ws.GET(resource).To(restfulGetResource(getter, reqScope)).
      Doc("read the specified "+kind).
      Operation("read"+kind).
      Produces(mediaTypes...).
      Returns(http.StatusOK, "OK", defaultVersionedObject).
      Writes(defaultVersionedObject)
    routes = append(routes, getRoute)
    ws.Route(getRoute)
  }

  if isCreater {
    var postHandler = restfulCreateResource(creater, reqScope)
    createRoute := ws.POST(resource).To(postHandler).
      Doc("create "+kind).
      Operation("create"+kind).
      Produces(mediaTypes...).
      Returns(http.StatusOK, "OK", defaultVersionedObject).
      Returns(http.StatusCreated, "OK", defaultVersionedObject).
      Returns(http.StatusAccepted, "OK", defaultVersionedObject).
      Reads(defaultVersionedObject).
      Writes(defaultVersionedObject)
    routes = append(routes, createRoute)
    ws.Route(createRoute)
  }

  if gvkProvider, ok := storage.(rest.GroupVersionKindProvider); ok {
    gvk := gvkProvider.GroupVersionKind(a.group.GroupVersion)
    apiResource.Group = gvk.Group
    apiResource.Version = gvk.Version
    apiResource.Kind = gvk.Kind

  }

  return &apiResource, nil

}

// indirectArbitraryPointer returns *ptrToObject for an arbitrary pointer
func indirectArbitraryPointer(ptrToObject interface{}) interface{} {
  return reflect.Indirect(reflect.ValueOf(ptrToObject)).Interface()
}

// GetResourceKind returns the external group version kind registered for the given storage
// object. If the storage object is a subresource and has an override supplied for it, it returns
// the group version kind supplied in the override.
func GetResourceKind(groupVersion schema.GroupVersion, storage rest.Storage, typer runtime.ObjectTyper) (schema.GroupVersionKind, error) {
  // Let the storage tell us exactly what GVK it has
  //if gvkProvider, ok := storage.(rest.GroupVersionKindProvider); ok {
  //  return gvkProvider.GroupVersionKind(groupVersion), nil
  //}

  object := storage.New()
  fqKinds, _, err := typer.ObjectKinds(object)
  if err != nil {
    return schema.GroupVersionKind{}, err
  }

  // a given go type can have multiple potential fully qualified kinds.  Find the one that corresponds with the group
  // we're trying to register here
  fqKindToRegister := schema.GroupVersionKind{}
  for _, fqKind := range fqKinds {
    if fqKind.Group == groupVersion.Group {
      fqKindToRegister = groupVersion.WithKind(fqKind.Kind)
      break
    }
  }
  if fqKindToRegister.Empty() {
    return schema.GroupVersionKind{}, fmt.Errorf("unable to locate fully qualified kind for %v: found %v when registering for %v", reflect.TypeOf(object), fqKinds, groupVersion)
  }

  // group is guaranteed to match based on the check above
  return fqKindToRegister, nil
}

func restfulGetResource(r rest.Getter, scope handlers.RequestScope) restful.RouteFunction {
  return func(req *restful.Request, res *restful.Response) {
    handlers.GetResource(r, &scope)(res.ResponseWriter, req.Request)
  }
}

func restfulCreateResource(r rest.Creater, scope handlers.RequestScope) restful.RouteFunction {
  return func(req *restful.Request, res *restful.Response) {
    handlers.CreateResource(r, &scope)(res.ResponseWriter, req.Request)
  }
}
