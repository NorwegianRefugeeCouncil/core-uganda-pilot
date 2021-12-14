package form

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/api/types/validation"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	uuid "github.com/satori/go.uuid"
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
		if errs := validation.ValidateForm(&form); !errs.IsEmpty() {
			l.Error("failed to validate form", zap.Error(errs.ToAggregate()))
			utils.ErrorResponse(w, meta.NewInvalid(meta.GroupResource{
				Group:    "core.nrc.no/v1",
				Resource: "forms",
			}, "", errs))
			return
		}

		f := &form
		newFormIDs(f)

		l.Debug("storing form")
		respForm, err := h.store.Create(ctx, f)
		if err != nil {
			l.Error("failed to store form", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully created form")
		utils.JSONResponse(w, http.StatusOK, respForm)
	}
}

func newFormIDs(form *types.FormDefinition) {
	form.ID = uuid.NewV4().String()
	newFieldIDs(form.Fields)
}

func newFieldIDs(fields []*types.FieldDefinition) {
	for _, field := range fields {
		field.ID = uuid.NewV4().String()
		if field.FieldType.SubForm != nil {
			newFieldIDs(field.FieldType.SubForm.Fields)
		}
		if field.FieldType.SingleSelect != nil {
			for _, option := range field.FieldType.SingleSelect.Options {
				option.ID = uuid.NewV4().String()
			}
		}
		if field.FieldType.MultiSelect != nil {
			for _, option := range field.FieldType.MultiSelect.Options {
				option.ID = uuid.NewV4().String()
			}
		}
	}
}

func (h *Handler) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := h.Create()
	handler(response.ResponseWriter, request.Request)
}
