package database

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/utils"
	"net/http"
)

func (h *Handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		dbs, err := h.store.List()
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		utils.JSONResponse(w, http.StatusOK, dbs)
	}
}

func (h *Handler) RestfulList(request *restful.Request, response *restful.Response) {
	handler := h.List()
	handler(response.ResponseWriter, request.Request)
}
