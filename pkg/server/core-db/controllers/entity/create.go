package entity

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
)

func (c *Controller) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (c *Controller) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := c.Create()
	handler(response.ResponseWriter, request.Request)
}
