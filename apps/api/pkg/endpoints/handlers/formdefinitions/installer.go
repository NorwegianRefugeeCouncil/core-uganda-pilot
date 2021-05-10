package formdefinitions

import (
	"github.com/emicklei/go-restful"
	"github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/storage"
	"path"
)

func Install(
	container *restful.Container,
	storage storage.Interface,
	getter rest.Getter,
	kind schema.GroupVersionKind,
	resource schema.GroupVersionResource,
	creater runtime.ObjectCreater,
	typer runtime.ObjectTyper,
	serializer runtime.Serializer,
	scheme *runtime.Scheme,
) {

	scope := handlers.NewRequestScope(
		kind,
		resource,
		creater,
		typer,
		serializer,
		scheme,
		func() runtime.Object { return &v1.FormDefinitionList{} },
	)

	ws := new(restful.WebService)
	apiPath := path.Join("apis", scope.Resource.GroupVersion().String(), scope.Resource.Resource)
	ws.Path(apiPath)
	ws.Doc("API at " + apiPath)

	handler := Handler{
		storage: storage,
		scope:   scope,
		getter:  getter,
	}

	ws.Route(ws.GET("/").To(func(request *restful.Request, response *restful.Response) {
		handler.List(response.ResponseWriter, request.Request)
	}).
		Operation("getFormDefinitions").
		Produces("application/json").
		Writes(&v1.FormDefinition{}),
	)

	ws.Route(ws.POST("/").To(func(request *restful.Request, response *restful.Response) {
		handler.Post(response.ResponseWriter, request.Request)
	}).
		Operation("postFormDefinition").
		Produces("application/json").
		Consumes("application/json").
		Writes(&v1.FormDefinition{}).
		Reads(&v1.FormDefinition{}))

	ws.Route(ws.PUT("/{id}").To(func(request *restful.Request, response *restful.Response) {
		handler.Update(response.ResponseWriter, request.Request)
	}).
		Operation("updateFormDefinition").
		Produces("application/json").
		Consumes("application/json").
		Writes(&v1.FormDefinition{}).
		Reads(&v1.FormDefinition{}))

	ws.Route(ws.GET("/watch").To(func(request *restful.Request, response *restful.Response) {
		handler.Watch(response.ResponseWriter, request.Request)
	}).
		Operation("watchFormDefinition").
		Produces("application/json").
		Consumes("application/json").
		Writes(&v1.FormDefinition{}).
		Reads(&v1.FormDefinition{}))

	container.Add(ws)

}
