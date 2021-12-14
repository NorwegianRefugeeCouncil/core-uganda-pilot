package foo

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
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

		var reads types.FooReads
		if err := utils.BindJSON(req, &reads); err != nil {
			l.Error("Failed to bind request body", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		if errList := validation.ValidateFooReads(&reads); !errList.IsEmpty() {
			err := meta.NewInvalid(types.FooGR, "", errList)
			l.Error("Failed to validate request body", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		// I don't like how there's no validation that each field is set.
		// I get that they will default to a zero value
		// But that isn't the same as forgetting to set it.
		// Eg. here the validation will pass if I forget to set "valid"
		// because it will default to false.
		foo := types.Foo{
			ID:         uuid.NewV4(),
			Name:       *reads.Name,
			OtherField: *reads.OtherField,
			Valid:      *reads.Valid,
		}

		createdFoo, err := h.store.Create(ctx, &foo)
		if err != nil {
			l.Error("Failed to create foo", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		utils.JSONResponse(w, http.StatusCreated, createdFoo)
	}
}

func (h *Handler) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := h.Create()
	handler(response.ResponseWriter, request.Request)
}
