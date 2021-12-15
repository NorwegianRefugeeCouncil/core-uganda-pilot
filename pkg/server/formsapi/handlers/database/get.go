package database

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) Get(databaseId string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("database_id", databaseId))

		l.Debug("getting databases from store")
		db, err := h.store.Get(ctx, databaseId)
		if err != nil {
			l.Error("failed to get database from store", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully got database")
		utils.JSONResponse(w, http.StatusOK, db)
	}
}

func (h *Handler) RestfulGet(request *restful.Request, response *restful.Response) {
	handler := h.Get(request.PathParameter("databaseId"))
	handler(response.ResponseWriter, request.Request)
}
