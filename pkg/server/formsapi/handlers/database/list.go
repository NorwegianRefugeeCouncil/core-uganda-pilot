package database

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		l := logging.NewLogger(ctx)

		l.Debug("listing databases in store")
		dbs, err := h.store.List(ctx)
		if err != nil {
			l.Error("failed to list databases in store", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully listed databases", zap.Int("count", len(dbs.Items)))
		utils.JSONResponse(w, http.StatusOK, dbs)
	}
}

func (h *Handler) RestfulList(request *restful.Request, response *restful.Response) {
	handler := h.List()
	handler(response.ResponseWriter, request.Request)
}
