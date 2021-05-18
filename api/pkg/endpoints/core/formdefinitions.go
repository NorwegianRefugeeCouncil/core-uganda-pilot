package core

import (
	"github.com/emicklei/go-restful/v3"
	v1 "github.com/nrc-no/coreapi/pkg/apis/core/v1"
	"net/http"
)

func Install(container *restful.Container) {

	ws := &restful.WebService{}
	ws.Doc("API for formdefinitions")
	ws.Path("/apis/core/v1/formdefinitions")
	ws.ApiVersion("core/v1")

	ws.Route(ws.GET("").To(restfulList()).
		Doc("list FormDefinitions").
		Produces("application/json").
		Writes(&v1.FormDefinitionList{}).
		Returns(http.StatusOK, "OK", &v1.FormDefinitionList{}).
		Operation("listFormDefinitions"))

	ws.Route(ws.GET("/{name}").To(restfulGet()).
		Doc("get a FormDefinition").
		Produces("application/json").
		Writes(&v1.FormDefinition{}).
		Returns(http.StatusOK, "OK", &v1.FormDefinition{}).
		Operation("getFormDefinition"))

	ws.Route(ws.POST("").To(restfulCreate()).
		Doc("create FormDefinition").
		Produces("application/json").
		Writes(&v1.FormDefinition{}).
		Returns(http.StatusOK, "OK", &v1.FormDefinition{}).
		Returns(http.StatusCreated, "OK", &v1.FormDefinition{}).
		Operation("postFormDefinition"))

	ws.Route(ws.PUT("/{name}").To(restfulPut()).
		Doc("put FormDefinition").
		Produces("application/json").
		Writes(&v1.FormDefinition{}).
		Returns(http.StatusOK, "OK", &v1.FormDefinition{}).
		Operation("putFormDefinition"))

	ws.Route(ws.DELETE("/{name}").To(restfulDelete()).
		Doc("delete a FormDefinition").
		Produces("application/json").
		Writes(&v1.FormDefinition{}).
		Returns(http.StatusOK, "OK", nil).
		Operation("deleteFormDefinition"))

}

func restfulDelete() func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {

	}
}

func restfulPut() func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {

	}
}

func restfulGet() func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {

	}
}

func restfulCreate() func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {

	}
}

func restfulList() func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {

	}
}
