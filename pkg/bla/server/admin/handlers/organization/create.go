package organization

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/bla/types"
	"github.com/nrc-no/core/pkg/utils"
	"net/http"
)

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var org types.Organization
		if err := utils.BindJSON(req, &org); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		result, err := h.store.Create(req.Context(), &org)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		utils.JSONResponse(w, http.StatusOK, result)
	}
}

func (h *Handler) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := h.Create()
	handler(response.ResponseWriter, request.Request)
}
