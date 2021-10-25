package record

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/bla/types"
	"github.com/nrc-no/core/pkg/utils"
	"net/http"
)

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var input types.Record
		if err := utils.BindJSON(req, &input); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		resultRecord, err := h.store.Create(req.Context(), &input)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		utils.JSONResponse(w, http.StatusOK, resultRecord)
	}
}

func (h *Handler) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := h.Create()
	handler(response.ResponseWriter, request.Request)
}