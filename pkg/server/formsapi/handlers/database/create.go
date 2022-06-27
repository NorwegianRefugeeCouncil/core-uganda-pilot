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
	"log"
	"net/http"
)

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		l := logging.NewLogger(ctx)

		var db types.Database
		if err := utils.BindJSON(req, &db); err != nil {
			l.Error("failed to unmarshal database", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		if validationErrs := validation.ValidateDatabase(&db); !validationErrs.IsEmpty() {
			err := meta.NewInvalid(meta.GroupResource{
				Group:    "core.nrc.no",
				Resource: "databases",
			}, "", validationErrs)
			l.Warn("database is invalid", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		db.ID = uuid.NewV4().String()

		respDB, err := h.store.Create(ctx, &db)
		if err != nil {
			l.Error("failed to store database", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		resp, err := h.zanzibarClient.WriteDB2UserRel(ctx, db.ID, "userId")
		if err != nil {
			log.Fatalf("failed to create relationship between database and creator: %s, %s", err, resp)
			return
		}

		utils.JSONResponse(w, http.StatusOK, respDB)
	}
}

func (h *Handler) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := h.Create()
	handler(response.ResponseWriter, request.Request)
}
