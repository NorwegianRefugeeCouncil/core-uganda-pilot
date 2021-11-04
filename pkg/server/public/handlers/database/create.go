package database

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
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
