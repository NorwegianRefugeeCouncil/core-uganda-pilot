package record

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/types"
	"github.com/nrc-no/core/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (h *Handler) Update(recordId string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		_, err := uuid.FromString(recordId)
		if err != nil {
			utils.ErrorResponse(w, meta.NewBadRequest("invalid record id"))
			return
		}

		var input types.Record
		if err := utils.BindJSON(req, &input); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		if len(input.ID) > 0 && input.ID != recordId {
			utils.ErrorResponse(w, meta.NewBadRequest("record id mismatch"))
			return
		}

		result, err := h.store.Update(req.Context(), &input)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		utils.JSONResponse(w, http.StatusOK, result)
	}
}

func (h *Handler) RestfulUpdate(request *restful.Request, response *restful.Response) {
	recordID := request.PathParameter(constants.ParamRecordID)
	handler := h.Update(recordID)
	handler(response.ResponseWriter, request.Request)
}
