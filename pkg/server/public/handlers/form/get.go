package form

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/utils"
	"net/http"
)

func (h *Handler) Get(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		forms, err := h.store.Get(ctx, id)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		utils.JSONResponse(w, http.StatusOK, forms)
	}
}

func (h *Handler) RestfulGet(request *restful.Request, response *restful.Response) {
	handler := h.Get(request.PathParameter(constants.ParamFormID))
	handler(response.ResponseWriter, request.Request)
}
