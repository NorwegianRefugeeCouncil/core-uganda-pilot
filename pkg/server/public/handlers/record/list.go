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

func (h *Handler) List(databaseID, formID string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		_, err := uuid.FromString(databaseID)
		if err != nil {
			utils.ErrorResponse(w, meta.NewBadRequest("invalid database id"))
			return
		}

		_, err = uuid.FromString(formID)
		if err != nil {
			utils.ErrorResponse(w, meta.NewBadRequest("invalid form id"))
			return
		}

		records, err := h.store.List(req.Context(), types.RecordListOptions{
			DatabaseID: databaseID,
			FormID:     formID,
		})

		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		utils.JSONResponse(w, http.StatusOK, records)
	}
}

func (h *Handler) RestfulList(request *restful.Request, response *restful.Response) {
	databaseID := request.QueryParameter(constants.ParamDatabaseID)
	formID := request.QueryParameter(constants.ParamFormID)
	handler := h.List(databaseID, formID)
	handler(response.ResponseWriter, request.Request)
}
