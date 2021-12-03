package database

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

		l.Debug("unmarshaling database")
		var db types.Database
		if err := utils.BindJSON(req, &db); err != nil {
			l.Error("failed to unmarshal database", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("validating database")
		if validationErrs := validation.ValidateDatabase(&db); validationErrs.HasAny() {
			err := meta.NewInvalid(meta.GroupResource{
				Group:    "core.nrc.no",
				Resource: "databases",
			}, "", validationErrs)
			l.Warn("database is invalid", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		db.ID = uuid.NewV4().String()

		l.Debug("storing database")
		respDB, err := h.store.Create(ctx, &db)
		if err != nil {
			l.Error("failed to store database", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully created database")
		utils.JSONResponse(w, http.StatusOK, respDB)
	}
}

func (h *Handler) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := h.Create()
	handler(response.ResponseWriter, request.Request)
}
