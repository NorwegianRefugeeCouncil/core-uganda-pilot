package form

import (
	"context"
	"net/http"

	validation2 "github.com/nrc-no/core/pkg/validation"

	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/api/types/defaults"
	"github.com/nrc-no/core/pkg/api/types/validation"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
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

		// apply defaults onto FormDefinition
		form = defaults.FormDefinitionDefaults(form)

		l.Debug("validating form")
		if errs := validation.ValidateForm(&form); !errs.IsEmpty() {
			l.Error("failed to validate form", zap.Error(errs.ToAggregate()))
			utils.ErrorResponse(w, meta.NewInvalid(types.FormGR, "", errs))
			return
		}

		if err := validateRecipientFormParent(ctx, form, h.store.Get); err != nil {
			l.Error("failed to validate recipient parent", zap.Error(err))
			utils.ErrorResponse(w, err)
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

func validateRecipientFormParent(
	ctx context.Context,
	form types.FormDefinition,
	getForm func(ctx context.Context, formID string) (*types.FormDefinition, error),
) error {

	l := logging.NewLogger(ctx)

	if form.Type != types.RecipientFormType {
		return nil
	}

	l.Debug("validation recipient form")
	for _, field := range form.Fields {

		if !field.Key {
			continue
		}

		// At this point, the validation has already asserted that the form has either no key fields,
		// or a single "reference" key field. So we can safely assume that this fieldType is
		// a Reference field
		referencedRecipient := field.FieldType.Reference
		parentRecipientForm, err := getForm(ctx, referencedRecipient.FormID)
		if err != nil {
			l.Error("failed to get parent recipient form", zap.Error(err))
			return err
		}

		// The rule is that a recipient form must either have no parent at all, or another
		// recipient form.
		if parentRecipientForm.Type != types.RecipientFormType {
			return meta.NewInvalid(types.FormGR, "", validation2.ErrorList{
				validation2.Forbidden(validation2.NewPath("type"), "The parent form of a recipient form must also be a recipient form"),
			})
		}
	}

	return nil
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
