package endpoints

import (
	"fmt"
	"github.com/emicklei/go-restful"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/negotiation"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/storageversion"
	"github.com/nrc-no/core/apps/api/pkg/util/conversion"
	"k8s.io/apimachinery/pkg/util/sets"
	"net/http"
	path2 "path"
	"reflect"
	"sort"
	"strings"
	"time"
)

type APIInstaller struct {
	group             *APIGroupVersion
	prefix            string // Path prefix where API resources are to be registered.
	minRequestTimeout time.Duration
}

// Struct capturing information about an action ("GET", "POST", "WATCH", "PROXY", etc).
type action struct {
	Verb          string               // Verb identifying the action ("GET", "POST", "WATCH", "PROXY", etc).
	Path          string               // The path of the action
	Params        []*restful.Parameter // List of parameters associated with the action.
	Namer         handlers.ScopeNamer
	AllNamespaces bool // true iff the action is namespaced but works on aggregate result for all namespaces
}

func (a *APIInstaller) Install() ([]metav1.APIResource, []*storageversion.ResourceInfo, *restful.WebService, []error) {

	var apiResources []metav1.APIResource
	var resourceInfos []*storageversion.ResourceInfo
	var errors []error

	ws := a.newWebService()

	paths := make([]string, len(a.group.Storage))
	var i int = 0
	for path := range a.group.Storage {
		paths[i] = path
		i++
	}

	sort.Strings(paths)
	for _, path := range paths {

		apiResource, resourceInfo, err := a.registerResourceHandlers(path, a.group.Storage[path], ws)
		if err != nil {
			errors = append(errors, fmt.Errorf("error in registering resource %v: %v", path, err))
		}
		if apiResource != nil {
			apiResources = append(apiResources, *apiResource)
		}
		if resourceInfo != nil {
			resourceInfos = append(resourceInfos, resourceInfo)
		}
	}
	return apiResources, resourceInfos, ws, errors
}

// newWebService creates a new restful webservice with the api installer's prefix and version.
func (a *APIInstaller) newWebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(a.prefix)
	// a.prefix contains "prefix/group/version"
	ws.Doc("API at " + a.prefix)
	// Backwards compatibility, we accepted objects with empty content-type at V1.
	// If we stop using go-restful, we can default empty content-type to application/json on an
	// endpoint by endpoint basis
	ws.Consumes("*/*")
	mediaTypes, streamMediaTypes := negotiation.MediaTypesForSerializer(a.group.Serializer)
	ws.Produces(append(mediaTypes, streamMediaTypes...)...)
	ws.ApiVersion(a.group.GroupVersion.String())

	return ws
}

func (a *APIInstaller) registerResourceHandlers(path string, storage rest.Storage, ws *restful.WebService) (*metav1.APIResource, *storageversion.ResourceInfo, error) {
	//
	optionsExternalVersion := a.group.GroupVersion
	//
	resource, _, err := splitSubresource(path)
	//if err != nil {
	//  return nil, nil, err
	//}

	//group, version := a.group.GroupVersion.Group, a.group.GroupVersion.Version

	fqKindToRegister, err := GetResourceKind(a.group.GroupVersion, storage, a.group.Typer)
	if err != nil {
		return nil, nil, err
	}

	versionedPtr, err := a.group.Creater.New(fqKindToRegister)
	if err != nil {
		return nil, nil, err
	}

	defaultVersionedObject := indirectArbitraryPointer(versionedPtr)
	kind := fqKindToRegister.Kind
	// isSubresource := len(subresource) > 0

	creater, isCreater := storage.(rest.Creater)
	lister, isLister := storage.(rest.Lister)
	watcher, isWatcher := storage.(rest.Watcher)
	getter, isGetter := storage.(rest.Getter)
	deleter, isDeleter := storage.(rest.Deleter)
	updater, isUpdater := storage.(rest.Updater)
	storageMeta := defaultStorageMetadata{}

	var versionedList interface{}
	if isLister {
		list := lister.NewList()
		listGVKs, _, err := a.group.Typer.ObjectKinds(list)
		if err != nil {
			return nil, nil, err
		}
		versionedListPtr, err := a.group.Creater.New(a.group.GroupVersion.WithKind(listGVKs[0].Kind))
		if err != nil {
			return nil, nil, err
		}
		versionedList = indirectArbitraryPointer(versionedListPtr)
	}

	versionedListOptions, err := a.group.Creater.New(optionsExternalVersion.WithKind("ListOptions"))
	if err != nil {
		return nil, nil, err
	}
	versionedCreateOptions, err := a.group.Creater.New(optionsExternalVersion.WithKind("CreateOptions"))
	if err != nil {
		return nil, nil, err
	}
	versionedUpdateOptions, err := a.group.Creater.New(optionsExternalVersion.WithKind("UpdateOptions"))
	if err != nil {
		return nil, nil, err
	}
	//var versionedDeleteOptions runtime.Object
	// var versionedDeleterObject interface{}
	//if isDeleter {
	//versionedDeleteOptions, err = a.group.Creater.New(optionsExternalVersion.WithKind("DeleteOptions"))
	//if err != nil {
	//  return nil, nil, err
	//}
	//versionedDeleterObject = indirectArbitraryPointer(versionedDeleteOptions)
	//}

	versionedStatusPtr, err := a.group.Creater.New(optionsExternalVersion.WithKind("Status"))
	if err != nil {
		return nil, nil, err
	}
	versionedStatus := indirectArbitraryPointer(versionedStatusPtr)

	//var (
	//  getOptions             runtime.Object
	//  versionedGetOptions    runtime.Object
	//  getOptionsInternalKind schema.GroupVersionKind
	//  getSubPath             bool
	//)

	var idParam = ws.PathParameter("id", "id of the "+kind).DataType("string")
	// var pathParam = ws.PathParameter("path", "path to the resource").DataType("string")

	params := []*restful.Parameter{}
	actions := []action{}

	var resourceKind string
	kindProvider, ok := storage.(rest.KindProvider)
	if ok {
		resourceKind = kindProvider.Kind()
	} else {
		resourceKind = kind
	}

	var apiResource metav1.APIResource

	resourcePath := resource
	resourceParams := params
	itemPath := resourcePath + "/{id}"
	idParams := append(params, idParam)
	suffix := ""

	apiResource.Name = path
	apiResource.Namespaced = false
	apiResource.Kind = resourceKind
	namer := handlers.ContextBasedNaming{
		ClusterScoped:      true,
		SelfLinkPathPrefix: path2.Join(a.prefix, resource) + "/",
		SelfLinkPathSuffix: suffix,
	}

	actions = appendIf(actions, action{"LIST", resourcePath, resourceParams, namer, false}, isLister)
	actions = appendIf(actions, action{"POST", resourcePath, resourceParams, namer, false}, isCreater)
	actions = appendIf(actions, action{"GET", itemPath, idParams, namer, false}, isGetter)
	actions = appendIf(actions, action{"PUT", itemPath, idParams, namer, false}, isUpdater)
	actions = appendIf(actions, action{"DELETE", itemPath, idParams, namer, false}, isDeleter)

	var resourceInfo *storageversion.ResourceInfo

	for _, serializerInfo := range a.group.Serializer.SupportedMediaTypes() {
		if len(serializerInfo.MediaTypeSubType) == 0 || len(serializerInfo.MediaTypeType) == 0 {
			return nil, nil, fmt.Errorf("all serializers in the group Serializer must have MediaTypeType and MediaTypeSubType set: %s", serializerInfo.MediaType)
		}
	}

	mediaTypes, _ := negotiation.MediaTypesForSerializer(a.group.Serializer)
	allMediaTypes := append(mediaTypes)

	ws.Produces(allMediaTypes...)

	kubeVerbs := map[string]struct{}{}
	reqScope := handlers.RequestScope{
		Serializer:               a.group.Serializer,
		ParameterCodec:           a.group.ParameterCodec,
		Creater:                  a.group.Creater,
		Convertor:                a.group.Convertor,
		Typer:                    a.group.Typer,
		EquivalentResourceMapper: a.group.EquivalentResourceRegistry,
		Resource:                 a.group.GroupVersion.WithResource(resource),
		Kind:                     fqKindToRegister,
		HubGroupVersion:          schema.GroupVersion{Group: fqKindToRegister.Group, Version: runtime.APIVersionInternal},
	}

	for _, action := range actions {

		producedObject := storageMeta.ProducesObject(action.Verb)
		if producedObject == nil {
			producedObject = defaultVersionedObject
		}

		reqScope.Namer = action.Namer
		if kubeVerb, found := toDiscoveryKubeVerb[action.Verb]; found {
			if len(kubeVerb) != 0 {
				kubeVerbs[kubeVerb] = struct{}{}
			}
		} else {
			return nil, nil, fmt.Errorf("unknown action verb for discovery: %v", action.Verb)
		}

		routes := []*restful.RouteBuilder{}

		switch action.Verb {
		case "GET":
			doc := "read the specified " + kind
			route := ws.GET(action.Path).To(restfulGetResource(getter, reqScope)).
				Doc(doc).
				Param(ws.QueryParameter("pretty", "If 'true', then the output is pretty-printed.")).
				Operation("read"+kind).
				Produces(append(storageMeta.ProducesMIMETypes(action.Verb), mediaTypes...)...).
				Returns(http.StatusOK, "OK", producedObject).
				Writes(producedObject)
			addParams(route, action.Params)
			routes = append(routes, route)
		case "LIST":
			route := ws.GET(action.Path).To(restfulListResource(lister, watcher, reqScope)).
				Doc("list objects of kind "+kind).
				Param(ws.QueryParameter("pretty", "if 'true', then the output is pretty-printed")).
				Operation("list"+kind).
				Produces(append(storageMeta.ProducesMIMETypes(action.Verb), allMediaTypes...)...).
				Returns(http.StatusOK, "OK", versionedList).
				Writes(versionedList)
			if err := AddObjectParams(ws, route, versionedListOptions); err != nil {
				return nil, nil, err
			}
			switch {
			case isLister && isWatcher:
				route.Doc("list or watch objects of kind " + kind)
			case isWatcher:
				route.Doc("watch objects of kind " + kind)
			}
			addParams(route, action.Params)
			routes = append(routes, route)
		case "PUT":
			route := ws.PUT(action.Path).To(restfulUpdateResource(updater, reqScope)).
				Doc("replace the specified "+kind).
				Param(ws.QueryParameter("pretty", "if 'true', then the output is pretty-printed.")).
				Operation("replace"+kind).
				Produces(append(storageMeta.ProducesMIMETypes(action.Verb), mediaTypes...)...).
				Returns(http.StatusOK, "OK", producedObject).
				Returns(http.StatusCreated, "Created", producedObject).
				Reads(defaultVersionedObject).
				Writes(producedObject)
			if err := AddObjectParams(ws, route, versionedUpdateOptions); err != nil {
				return nil, nil, err
			}
			addParams(route, action.Params)
			routes = append(routes, route)
		case "PATCH":
		case "POST":
			route := ws.POST(action.Path).To(restfulCreateResource(creater, reqScope)).
				Doc("create "+kind).
				Param(ws.QueryParameter("pretty", "if 'true', then the output is pretty-printed.")).
				Operation("create"+kind).
				Produces(append(storageMeta.ProducesMIMETypes(action.Verb), mediaTypes...)...).
				Returns(http.StatusOK, "OK", producedObject).
				Returns(http.StatusCreated, "Created", producedObject).
				Returns(http.StatusAccepted, "Accepted", producedObject).
				Reads(defaultVersionedObject).
				Writes(producedObject)
			if err := AddObjectParams(ws, route, versionedCreateOptions); err != nil {
				return nil, nil, err
			}
			addParams(route, action.Params)
			routes = append(routes, route)

		case "DELETE":
			deleteReturnType := versionedStatus
			route := ws.DELETE(action.Path).To(restfulDeleteResource(deleter, reqScope)).
				Param(ws.QueryParameter("pretty", "if 'true', then the output is pretty-printed.")).
				Operation("delete"+kind).
				Produces(append(storageMeta.ProducesMIMETypes(action.Verb), mediaTypes...)...).
				Writes(deleteReturnType).
				Returns(http.StatusOK, "OK", deleteReturnType).
				Returns(http.StatusAccepted, "Accepted", deleteReturnType)
			addParams(route, action.Params)
			routes = append(routes, route)
		default:
			return nil, nil, fmt.Errorf("unrecognized action verb: %s", action.Verb)
		}

		for _, route := range routes {
			ws.Route(route)
		}

	}

	apiResource.Verbs = make([]string, 0, len(kubeVerbs))
	for kubeVerb := range kubeVerbs {
		apiResource.Verbs = append(apiResource.Verbs, kubeVerb)
	}
	sort.Strings(apiResource.Verbs)

	if gvkProvider, ok := storage.(rest.GroupVersionKindProvider); ok {
		gvk := gvkProvider.GroupVersionKind(a.group.GroupVersion)
		apiResource.Group = gvk.Group
		apiResource.Version = gvk.Version
		apiResource.Kind = gvk.Kind
	}

	a.group.EquivalentResourceRegistry.RegisterKindFor(reqScope.Resource, "", fqKindToRegister)
	return &apiResource, resourceInfo, nil

}

func addParams(route *restful.RouteBuilder, params []*restful.Parameter) {
	for _, param := range params {
		route.Param(param)
	}
}

func appendIf(actions []action, a action, shouldAppend bool) []action {
	if shouldAppend {
		actions = append(actions, a)
	}
	return actions
}

// indirectArbitraryPointer returns *ptrToObject for an arbitrary pointer
func indirectArbitraryPointer(ptrToObject interface{}) interface{} {
	return reflect.Indirect(reflect.ValueOf(ptrToObject)).Interface()
}

// splitSubresource checks if the given storage path is the path of a subresource and returns
// the resource and subresource components.
func splitSubresource(path string) (string, string, error) {
	var resource, subresource string
	switch parts := strings.Split(path, "/"); len(parts) {
	case 2:
		resource, subresource = parts[0], parts[1]
	case 1:
		resource = parts[0]
	default:
		// TODO: support deeper paths
		return "", "", fmt.Errorf("api_installer allows only one or two segment paths (resource or resource/subresource)")
	}
	return resource, subresource, nil
}

// GetResourceKind returns the external group version kind registered for the given storage
// object. If the storage object is a subresource and has an override supplied for it, it returns
// the group version kind supplied in the override.
func GetResourceKind(groupVersion schema.GroupVersion, storage rest.Storage, typer runtime.ObjectTyper) (schema.GroupVersionKind, error) {
	// Let the storage tell us exactly what GVK it has
	if gvkProvider, ok := storage.(rest.GroupVersionKindProvider); ok {
		return gvkProvider.GroupVersionKind(groupVersion), nil
	}

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

// defaultStorageMetadata provides default answers to rest.StorageMetadata.
type defaultStorageMetadata struct{}

// defaultStorageMetadata implements rest.StorageMetadata
var _ rest.StorageMetadata = defaultStorageMetadata{}

func (defaultStorageMetadata) ProducesMIMETypes(verb string) []string {
	return nil
}

func (defaultStorageMetadata) ProducesObject(verb string) interface{} {
	return nil
}

// toDiscoveryKubeVerb maps an action.Verb to the logical kube verb, used for discovery
var toDiscoveryKubeVerb = map[string]string{
	"CONNECT":          "", // do not list in discovery.
	"DELETE":           "delete",
	"DELETECOLLECTION": "deletecollection",
	"GET":              "get",
	"LIST":             "list",
	"PATCH":            "patch",
	"POST":             "create",
	"PROXY":            "proxy",
	"PUT":              "update",
	"WATCH":            "watch",
	"WATCHLIST":        "watch",
}

// An interface to see if an object supports swagger documentation as a method
type documentable interface {
	SwaggerDoc() map[string]string
}

// AddObjectParams converts a runtime.Object into a set of go-restful Param() definitions on the route.
// The object must be a pointer to a struct; only fields at the top level of the struct that are not
// themselves interfaces or structs are used; only fields with a json tag that is non empty (the standard
// Go JSON behavior for omitting a field) become query parameters. The name of the query parameter is
// the JSON field name. If a description struct tag is set on the field, that description is used on the
// query parameter. In essence, it converts a standard JSON top level object into a query param schema.
func AddObjectParams(ws *restful.WebService, route *restful.RouteBuilder, obj interface{}, excludedNames ...string) error {
	sv, err := conversion.EnforcePtr(obj)
	if err != nil {
		return err
	}
	st := sv.Type()
	excludedNameSet := sets.NewString(excludedNames...)
	switch st.Kind() {
	case reflect.Struct:
		for i := 0; i < st.NumField(); i++ {
			name := st.Field(i).Name
			sf, ok := st.FieldByName(name)
			if !ok {
				continue
			}
			switch sf.Type.Kind() {
			case reflect.Interface, reflect.Struct:
			case reflect.Ptr:
				// TODO: This is a hack to let metav1.Time through. This needs to be fixed in a more generic way eventually. bug #36191
				if (sf.Type.Elem().Kind() == reflect.Interface || sf.Type.Elem().Kind() == reflect.Struct) && strings.TrimPrefix(sf.Type.String(), "*") != "metav1.Time" {
					continue
				}
				fallthrough
			default:
				jsonTag := sf.Tag.Get("json")
				if len(jsonTag) == 0 {
					continue
				}
				jsonName := strings.SplitN(jsonTag, ",", 2)[0]
				if len(jsonName) == 0 {
					continue
				}
				if excludedNameSet.Has(jsonName) {
					continue
				}
				var desc string
				if docable, ok := obj.(documentable); ok {
					desc = docable.SwaggerDoc()[jsonName]
				}
				route.Param(ws.QueryParameter(jsonName, desc).DataType(typeToJSON(sf.Type.String())))
			}
		}
	}
	return nil
}

// TODO: this is incomplete, expand as needed.
// Convert the name of a golang type to the name of a JSON type
func typeToJSON(typeName string) string {
	switch typeName {
	case "bool", "*bool":
		return "boolean"
	case "uint8", "*uint8", "int", "*int", "int32", "*int32", "int64", "*int64", "uint32", "*uint32", "uint64", "*uint64":
		return "integer"
	case "float64", "*float64", "float32", "*float32":
		return "number"
	case "metav1.Time", "*metav1.Time":
		return "string"
	case "byte", "*byte":
		return "string"
	case "v1.DeletionPropagation", "*v1.DeletionPropagation":
		return "string"
	case "v1.ResourceVersionMatch", "*v1.ResourceVersionMatch":
		return "string"
	case "v1.IncludeObjectPolicy", "*v1.IncludeObjectPolicy":
		return "string"

	// TODO: Fix these when go-restful supports a way to specify an array query param:
	// https://github.com/emicklei/go-restful/issues/225
	case "[]string", "[]*string":
		return "string"
	case "[]int32", "[]*int32":
		return "integer"

	default:
		return typeName
	}
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

func restfulListResource(r rest.Lister, rw rest.Watcher, scope handlers.RequestScope) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		handlers.ListResource(r, rw, &scope)(response.ResponseWriter, request.Request)
	}
}

func restfulUpdateResource(r rest.Updater, scope handlers.RequestScope) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		handlers.UpdateResource(r, &scope)(response.ResponseWriter, request.Request)
	}
}

func restfulDeleteResource(r rest.Deleter, scope handlers.RequestScope) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.DeleteResource(r, &scope)(res.ResponseWriter, req.Request)
	}
}
