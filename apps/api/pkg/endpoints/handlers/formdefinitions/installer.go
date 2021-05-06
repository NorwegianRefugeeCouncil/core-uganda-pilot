package formdefinitions

import (
	"github.com/emicklei/go-restful"
	"github.com/nrc-no/core/apps/api/pkg/apis"
	"github.com/nrc-no/core/apps/api/pkg/storage"
)

func Install(container *restful.Container, storage storage.Interface) {

	ws := new(restful.WebService)
	ws.Path("/apis/core.nrc.no/v1/formdefinitions")
	ws.Doc("API at /apis/core.nrc.no/v1/formdefinitions")

	handler := Handler{
		storage: storage,
	}

	ws.Route(ws.GET("/").To(func(request *restful.Request, response *restful.Response) {
		handler.Get(response.ResponseWriter, request.Request)
	}).
		Operation("getFormDefinitions").
		Produces("application/json").
		Writes(&apis.FormDefinition{}),
	)

	ws.Route(ws.POST("/").To(func(request *restful.Request, response *restful.Response) {
		handler.Post(response.ResponseWriter, request.Request)
	}).
		Operation("postFormDefinition").
		Produces("application/json").
		Consumes("application/json").
		Writes(&apis.FormDefinition{}).
		Reads(&apis.FormDefinition{}))

	ws.Route(ws.PUT("/{id}").To(func(request *restful.Request, response *restful.Response) {
		handler.Update(response.ResponseWriter, request.Request)
	}).
		Operation("updateFormDefinition").
		Produces("application/json").
		Consumes("application/json").
		Writes(&apis.FormDefinition{}).
		Reads(&apis.FormDefinition{}))

	container.Add(ws)

}
