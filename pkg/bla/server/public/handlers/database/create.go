package database

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/bla/types"
	"github.com/nrc-no/core/pkg/utils"
	"net/http"
)

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var db types.Database
		if err := utils.BindJSON(req, &db); err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		respDB, err := h.store.Create(&db)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		utils.JSONResponse(w, http.StatusOK, respDB)
	}
}

func (h *Handler) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := h.Create()
	handler(response.ResponseWriter, request.Request)
}
