package entity

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
)

func (c *Controller) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (c *Controller) RestfulGet(request *restful.Request, response *restful.Response) {
	handler := c.Get()
	handler(response.ResponseWriter, request.Request)
}
