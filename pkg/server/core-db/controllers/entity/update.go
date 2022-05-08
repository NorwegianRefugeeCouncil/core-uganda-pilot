package entity

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
)

func (c *Controller) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (c *Controller) RestfulUpdate(request *restful.Request, response *restful.Response) {
	handler := c.Update()
	handler(response.ResponseWriter, request.Request)
}
