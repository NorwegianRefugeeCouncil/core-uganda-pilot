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

func ReadFoo(req *http.Request, l *zap.Logger) (*types.FooReads, error) {
	var foo types.FooReads
	if err := utils.BindJSON(req, &foo); err != nil {
		l.Error("Failed to bind request body", zap.Error(err))
		return nil, err
	}

	if errList := validation.ValidateFooReads(&foo); !errList.IsEmpty() {
		err := meta.NewInvalid(types.FooGR, "", errList)
		l.Error("Failed to validate request body", zap.Error(err))
		return nil, err
	}
	return &foo, nil
}

func NewFoo(id uuid.UUID, name string, otherField int, uuidField uuid.UUID, valid bool, l *zap.Logger) (*types.Foo, error) {
	foo := types.Foo{
		ID:         id,
		Name:       name,
		OtherField: otherField,
		UUIDField:  uuidField,
		Valid:      valid,
	}

	// This is probably excessive
	if errList := validation.ValidateFoo(&foo); !errList.IsEmpty() {
		err := meta.NewInvalid(types.FooGR, "", errList)
		l.Error("Failed to validate foo", zap.Error(err))
		return nil, err
	}

	return &foo, nil
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		l := logging.NewLogger(ctx)

		fooReads, err := ReadFoo(req, l)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		foo, err := NewFoo(
			uuid.NewV4(),
			*fooReads.Name,
			*fooReads.OtherField,
			*fooReads.UUIDField,
			*fooReads.Valid,
			l,
		)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		createdFoo, err := h.store.Create(ctx, foo)
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
