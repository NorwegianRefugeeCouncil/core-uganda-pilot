package folder

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/bla/types"
	"github.com/nrc-no/core/pkg/utils"
	"net/http"
)

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		var folder types.Folder
		if err := utils.BindJSON(req, &folder); err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		respForm, err := h.store.Create(ctx, &folder)
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
