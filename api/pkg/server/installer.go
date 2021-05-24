package server

import (
	"fmt"
	"github.com/emicklei/go-restful"
	v1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	handlers2 "github.com/nrc-no/core/api/pkg/endpoints/handlers"
	rest2 "github.com/nrc-no/core/api/pkg/registry/rest"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
	"reflect"
	"sort"
)

type APIInstaller struct {
	group  *APIGroupVersion
	prefix string
}

func (i *APIInstaller) Install() (*restful.WebService, []v1.APIResource, error) {

	ws := new(restful.WebService)
	ws.Path(i.prefix)
	ws.Doc("API at " + i.prefix)
	ws.Consumes("application/json")
	ws.Produces("application/json")
	ws.ApiVersion(i.group.GroupVersion.String())

	paths := make([]string, len(i.group.Storage))
	var idx int = 0
	for path := range i.group.Storage {
		paths[idx] = path
		idx++
	}
	sort.Strings(paths)

	var apiResources []v1.APIResource
	for _, path := range paths {
		storage := i.group.Storage[path]
		apiResource, err := i.registerResourceHandlers(path, storage, ws)
		if err != nil {
			return nil, nil, err
		}
		apiResources = append(apiResources, *apiResource)
	}

	return ws, apiResources, nil
}

func (i *APIInstaller) registerResourceHandlers(path string, storage rest2.Storage, ws *restful.WebService) (*v1.APIResource, error) {

	resource := path
	creater, isCreater := storage.(rest2.Creater)
	lister, isLister := storage.(rest2.Lister)
	updater, isUpdater := storage.(rest2.Updater)
	getter, isGetter := storage.(rest2.Getter)
	deleter, isDeleter := storage.(rest2.Deleter)
	watcher, _ := storage.(rest2.Watcher)

	fqKindToRegister, err := GetResourceKind(i.group.GroupVersion, storage, i.group.Typer)
	if err != nil {
		return nil, err
	}
	kind := fqKindToRegister.Kind
	versionedPtr, err := i.group.Creater.New(fqKindToRegister)
	if err != nil {
		return nil, err
	}
	defaultVersionedObject := indirectArbitraryPointer(versionedPtr)

	var versionedList interface{}
	if isLister {
		list := lister.NewList()
		listGVKs, _, err := i.group.Typer.ObjectKinds(list)
		if err != nil {
			return nil, err
		}
		versionedListPtr, err := i.group.Creater.New(i.group.GroupVersion.WithKind(listGVKs[0].Kind))
		if err != nil {
			return nil, err
		}
		versionedList = indirectArbitraryPointer(versionedListPtr)
	}

	reqScope := handlers2.RequestScope{
		Serializer:      i.group.Serializer,
		ParameterCodec:  i.group.ParameterCodec,
		Creater:         i.group.Creater,
		Convertor:       i.group.Convertor,
		Defaulter:       i.group.Defaulter,
		Typer:           i.group.Typer,
		Resource:        i.group.GroupVersion.WithResource(resource),
		Kind:            fqKindToRegister,
		HubGroupVersion: schema.GroupVersion{Group: fqKindToRegister.Group, Version: runtime.APIVersionInternal},
		Namer:           handlers2.ContextBasedNaming{},
	}

	var verbs []string

	var routes []*restful.RouteBuilder
	if isGetter {
		verbs = append(verbs, "get")
		handler := restfulGetResource(getter, reqScope)
		route := ws.GET("/"+resource+"/{name}").To(handler).
			Doc("Gets a "+kind).
			Operation("read"+kind).
			Produces("application/json", "application/yaml").
			Returns(http.StatusOK, "OK", defaultVersionedObject).
			Writes(defaultVersionedObject)
		routes = append(routes, route)
	}

	if isLister {
		verbs = append(verbs, "list")
		handler := restfulListResource(lister, watcher, reqScope)
		route := ws.GET("/"+resource).To(handler).
			Doc("List object of kind "+kind).
			Operation("list"+kind).
			Produces("application/json", "application/yaml").
			Returns(http.StatusOK, "OK", versionedList).
			Writes(versionedList)
		routes = append(routes, route)
	}

	if isUpdater {
		verbs = append(verbs, "update")
		handler := restfulUpdateResource(updater, reqScope)
		route := ws.PUT("/"+resource+"/{name}").To(handler).
			Doc("Replaces a "+kind).
			Operation("replace"+kind).
			Consumes("application/json", "application/yaml").
			Produces("application/json", "application/yaml").
			Returns(http.StatusOK, "OK", defaultVersionedObject).
			Reads(defaultVersionedObject).
			Writes(defaultVersionedObject)
		routes = append(routes, route)
	}

	if isCreater {
		verbs = append(verbs, "create")
		handler := restfulCreateResource(creater, reqScope)
		route := ws.POST("/"+resource+"/").To(handler).
			Doc("Create a "+kind).
			Operation("create"+kind).
			Consumes("application/json", "application/yaml").
			Produces("application/json", "application/yaml").
			Returns(http.StatusCreated, "OK", defaultVersionedObject).
			Reads(defaultVersionedObject).
			Writes(defaultVersionedObject)
		routes = append(routes, route)
	}

	if isDeleter {
		verbs = append(verbs, "delete")
		handler := restfulDeleteResource(deleter, reqScope)
		route := ws.DELETE("/"+resource+"/{name}").To(handler).
			Doc("Delete a "+kind).
			Operation("delete"+kind).
			Consumes("application/json", "application/yaml").
			Produces("application/json", "application/yaml").
			Returns(http.StatusOK, "OK", defaultVersionedObject).
			Writes(defaultVersionedObject)
		routes = append(routes, route)
	}

	for _, route := range routes {
		ws.Route(route)
	}

	apiResource := &v1.APIResource{
		Group:      i.group.GroupVersion.Group,
		Version:    i.group.GroupVersion.Version,
		Name:       path,
		Namespaced: false,
		Kind:       fqKindToRegister.Kind,
		Verbs:      verbs,
	}

	return apiResource, nil
}

// indirectArbitraryPointer returns *ptrToObject for an arbitrary pointer
func indirectArbitraryPointer(ptrToObject interface{}) interface{} {
	return reflect.Indirect(reflect.ValueOf(ptrToObject)).Interface()
}

func GetResourceKind(groupVersion schema.GroupVersion, storage rest2.Storage, typer runtime.ObjectTyper) (schema.GroupVersionKind, error) {
	object := storage.New()
	fqKinds, _, err := typer.ObjectKinds(object)
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	fqKindToRegister := schema.GroupVersionKind{}
	for _, fqKind := range fqKinds {
		if fqKind.Group == groupVersion.Group {
			fqKindToRegister = groupVersion.WithKind(fqKind.Kind)
			break
		}
	}
	if fqKindToRegister.Empty() {
		return schema.GroupVersionKind{}, fmt.Errorf("could not locate fully qualified kind for %v: found %v when registering for %v", reflect.TypeOf(object), fqKinds, groupVersion)
	}

	return fqKindToRegister, nil
}

func restfulGetResource(r rest2.Getter, scope handlers2.RequestScope) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		handlers2.GetResource(&scope, r)(response.ResponseWriter, request.Request)
	}
}
func restfulCreateResource(r rest2.Creater, scope handlers2.RequestScope) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		handlers2.CreateResource(&scope, r)(response.ResponseWriter, request.Request)
	}
}
func restfulDeleteResource(r rest2.Deleter, scope handlers2.RequestScope) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		handlers2.DeleteResource(&scope, r)(response.ResponseWriter, request.Request)
	}
}
func restfulUpdateResource(r rest2.Updater, scope handlers2.RequestScope) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		handlers2.UpdateResource(&scope, r)(response.ResponseWriter, request.Request)
	}
}
func restfulListResource(r rest2.Lister, watcher rest2.Watcher, scope handlers2.RequestScope) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		handlers2.ListResource(&scope, r, watcher)(response.ResponseWriter, request.Request)
	}
}
