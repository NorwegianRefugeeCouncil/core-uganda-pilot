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

func (h *Handler) Update(recordId string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("record_id", recordId))

		l.Debug("validating record id")
		_, err := uuid.FromString(recordId)
		if err != nil {
			l.Error("failed to validate record id", zap.Error(err))
			utils.ErrorResponse(w, meta.NewBadRequest("invalid record id"))
			return
		}

		l.Debug("unmarshaling record")
		var input types.Record
		if err := utils.BindJSON(req, &input); err != nil {
			l.Error("failed to unmarshal record", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("ensuring record id from body and from path match")
		if len(input.ID) > 0 && input.ID != recordId {
			l.Error("failed to match record id in path and in body")
			utils.ErrorResponse(w, meta.NewBadRequest("record id mismatch"))
			return
		}

		l.Debug("updating record in store")
		result, err := h.store.Update(ctx, &input)
		if err != nil {
			l.Error("failed to update record in store", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully updated record")
		utils.JSONResponse(w, http.StatusOK, result)
	}
}

func (h *Handler) RestfulUpdate(request *restful.Request, response *restful.Response) {
	recordID := request.PathParameter(constants.ParamRecordID)
	handler := h.Update(recordID)
	handler(response.ResponseWriter, request.Request)
}
