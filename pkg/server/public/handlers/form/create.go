package form

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

		l.Debug("unmarshaling form")
		var form types.FormDefinition
		if err := utils.BindJSON(req, &form); err != nil {
			l.Error("failed to unmarshal form", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("validating form")
		if errs := validation.ValidateForm(&form); errs.HasAny() {
			l.Error("failed to validate form", zap.Error(errs.ToAggregate()))
			utils.ErrorResponse(w, meta.NewInvalid(meta.GroupResource{
				Group:    "core.nrc.no/v1",
				Resource: "forms",
			}, "", errs))
			return
		}

		l.Debug("storing form")
		respForm, err := h.store.Create(ctx, &form)
		if err != nil {
			l.Error("failed to store form", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully created form")
		utils.JSONResponse(w, http.StatusOK, respForm)
	}
}

func (h *Handler) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := h.Create()
	handler(response.ResponseWriter, request.Request)
}
