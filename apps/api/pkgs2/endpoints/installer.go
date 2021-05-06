package endpoints

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful"
	v1 "github.com/nrc-no/core/apps/api/pkgs2/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkgs2/runtime"
	"github.com/nrc-no/core/apps/api/pkgs2/runtime/schema"
	"net/http"
	"path"
	"reflect"
	"sort"
)

type APIGroupInfo struct {
	VersionedResourcesStorageMap map[string]map[string]Storage
	Scheme                       *runtime.Scheme
	PrioritizedVersions          []schema.GroupVersion
}

type APIGroupVersion struct {
	Storage      map[string]Storage
	GroupVersion schema.GroupVersion
	Typer        runtime.ObjectTyper
	Creater      runtime.ObjectCreater
}

type ResourceInfo struct {
	GroupResource schema.GroupResource
}

func (g *APIGroupVersion) InstallREST(handler http.Handler, container *restful.Container) ([]*ResourceInfo, error) {
	prefix := path.Join(g.GroupVersion.Group, g.GroupVersion.Version)
	installer := &APIInstaller{
		group:  g,
		prefix: prefix,
	}
}

type Storage interface {
	New() runtime.Object
}

type ValidateObjectFunc func(ctx context.Context, obj runtime.Object) error

type Creater interface {
	New() runtime.Object
	Create(ctx context.Context, obj runtime.Object, createValidation ValidateObjectFunc, options *v1.CreateOptions) (runtime.Object, error)
}

type Lister interface {
	NewList() runtime.Object
	List(ctx context.Context, options *v1.ListOptions) (runtime.Object, error)
}

type Getter interface {
	Get(ctx context.Context, name string, options runtime.Object) (runtime.Object, error)
}

type GroupVersionKindProvider interface {
	GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind
}

type APIInstaller struct {
	group  *APIGroupVersion
	prefix string
}

func (a *APIInstaller) Install() ([]v1.APIResource, *restful.WebService, error) {

	var apiResources []v1.APIResource
	ws := a.newWebService()

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
			return nil, nil, err
		}
		if apiResource != nil {
			apiResources = append(apiResources, *apiResource)
		}
	}
	return apiResources, ws, nil
}

func (a *APIInstaller) newWebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(a.prefix)
	ws.Doc("API at " + a.prefix)
	ws.Consumes("application/json")
	ws.Produces("application/json")
	ws.ApiVersion(a.group.GroupVersion.String())
	return ws
}

func (a *APIInstaller) registerResourceHandlers(path string, storage Storage, ws *restful.WebService) (*v1.APIResource, error) {

	kindToRegister, err := GetResourceKind(a.group.GroupVersion, storage, a.group.Typer)
	if err != nil {
		return nil, err
	}

	creater, isCreater := storage.(Creater)
	lister, isLister := storage.(Lister)
	getter, isGetter := storage.(Getter)

	var apiResource v1.APIResource
	{
	}
	apiResource.Name = path
	apiResource.Kind = kindToRegister.Kind

	requestScope := RequestScope{
		Creater:  a.group.Creater,
		Typer:    a.group.Typer,
		Resource: a.group.GroupVersion.WithResource(path),
		Kind:     kindToRegister,
	}

	versionedPtr, err := a.group.Creater.New(kindToRegister)
	if err != nil {
		return nil, err
	}
	defaultVersionedObject := reflect.Indirect(reflect.ValueOf(versionedPtr)).Interface()

	if isGetter {
		handler := restfulGetResource(getter, requestScope)
		ws.Route(ws.GET(path).To(handler).
			Operation("read"+kindToRegister.Kind).
			Produces("application/json").
			Returns(http.StatusOK, "OK", defaultVersionedObject).
			Writes(defaultVersionedObject))
	}
	if isCreater {

	}

}

func GetResourceKind(groupVersion schema.GroupVersion, storage Storage, typer runtime.ObjectTyper) (schema.GroupVersionKind, error) {
	if gvkProvider, ok := storage.(GroupVersionKindProvider); ok {
		return gvkProvider.GroupVersionKind(groupVersion), nil
	}
	object := storage.New()
	kinds, _, err := typer.ObjectKinds(object)
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	kindToRegister := schema.GroupVersionKind{}
	for _, kind := range kinds {
		if kind.Group == groupVersion.Group {
			kindToRegister = groupVersion.WithKind(kind.Kind)
			break
		}
	}
	if kindToRegister.Empty() {
		return schema.GroupVersionKind{}, fmt.Errorf("unable to locate fully qualified name for %v: found %v when registering for %v", reflect.TypeOf(object), kinds, groupVersion)
	}

	return kindToRegister, nil

}

type RequestScope struct {
	Creater  runtime.ObjectCreater
	Typer    runtime.ObjectTyper
	Resource schema.GroupVersionResource
	Kind     schema.GroupVersionKind
}
