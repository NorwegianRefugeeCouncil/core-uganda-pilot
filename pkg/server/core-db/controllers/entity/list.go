package entity

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
)

func (c *Controller) List() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (c *Controller) RestfulList(request *restful.Request, response *restful.Response) {
	handler := c.List()
	handler(response.ResponseWriter, request.Request)
}
