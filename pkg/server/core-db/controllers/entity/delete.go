package entity

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
)

func (c *Controller) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (c *Controller) RestfulDelete(request *restful.Request, response *restful.Response) {
	handler := c.Delete()
	handler(response.ResponseWriter, request.Request)
}
