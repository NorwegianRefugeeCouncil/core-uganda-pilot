package record

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/api/types/validation"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		l := logging.NewLogger(ctx)

		l.Debug("unmarshaling record")
		var input types.Record
		if err := utils.BindJSON(req, &input); err != nil {
			l.Error("failed to unmarshal record", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("getting form for record", zap.String("form_id", input.FormID))
		form, err := h.formStore.Get(req.Context(), input.FormID)
		if err != nil {
			l.Error("failed to get record form", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		input.DatabaseID = form.DatabaseID

		l.Debug("validating record")
		if errList := validation.ValidateRecord(&input, form); !errList.IsEmpty() {
			err := meta.NewInvalid(types.RecordGR, "", errList)
			l.Warn("record is invalid", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("storing record")
		resultRecord, err := h.store.Create(req.Context(), &input)
		if err != nil {
			l.Error("failed to store record", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully stored record")
		utils.JSONResponse(w, http.StatusOK, resultRecord)
	}
}

func (h *Handler) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := h.Create()
	handler(response.ResponseWriter, request.Request)
}
