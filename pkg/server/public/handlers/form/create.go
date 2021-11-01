package form

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/utils"
	"net/http"
)

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		var form types.FormDefinition
		if err := utils.BindJSON(req, &form); err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		respForm, err := h.store.Create(ctx, &form)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		utils.JSONResponse(w, http.StatusOK, respForm)
	}
}

func (h *Handler) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := h.Create()
	handler(response.ResponseWriter, request.Request)
}
