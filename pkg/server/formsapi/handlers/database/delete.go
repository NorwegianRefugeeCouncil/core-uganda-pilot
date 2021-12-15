package database

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) Delete(databaseId string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("database_id", databaseId))

		l.Debug("validating database id")
		_, err := uuid.FromString(databaseId)
		if err != nil {
			l.Error("failed to validate database id", zap.Error(err))
			utils.ErrorResponse(w, meta.NewBadRequest("invalid database id"))
			return
		}

		l.Debug("deleting database from store")
		if err := h.store.Delete(ctx, databaseId); err != nil {
			l.Error("failed to delete database from store", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully deleted database")
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func (h *Handler) RestfulDelete(request *restful.Request, response *restful.Response) {
	databaseID := request.PathParameter(constants.ParamDatabaseID)
	h.Delete(databaseID)(response.ResponseWriter, request.Request)
}
