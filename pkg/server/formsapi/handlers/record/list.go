package record

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) List(databaseID, formID string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("database_id", databaseID), zap.String("form_id", formID))

		l.Debug("validating database id")
		_, err := uuid.FromString(databaseID)
		if err != nil {
			l.Error("failed to validate database id", zap.Error(err))
			utils.ErrorResponse(w, meta.NewBadRequest("invalid database id"))
			return
		}

		l.Debug("validating form id")
		_, err = uuid.FromString(formID)
		if err != nil {
			l.Error("failed to validate form id", zap.Error(err))
			utils.ErrorResponse(w, meta.NewBadRequest("invalid form id"))
			return
		}

		l.Debug("listing records in store")
		records, err := h.store.List(ctx, types.RecordListOptions{
			DatabaseID: databaseID,
			FormID:     formID,
		})
		if err != nil {
			l.Error("failed to list records in store", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully listed records", zap.Int("count", len(records.Items)))
		utils.JSONResponse(w, http.StatusOK, records)
	}
}

func (h *Handler) RestfulList(request *restful.Request, response *restful.Response) {
	databaseID := request.QueryParameter(constants.ParamDatabaseID)
	formID := request.QueryParameter(constants.ParamFormID)
	handler := h.List(databaseID, formID)
	handler(response.ResponseWriter, request.Request)
}
