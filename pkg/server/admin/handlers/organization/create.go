package organization

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
		l := logging.NewLogger(req.Context())

		l.Debug("unmarshaling organization")
		var org types.Organization
		if err := utils.BindJSON(req, &org); err != nil {
			l.Error("failed to unmarshal organization", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("storing organization")
		result, err := h.store.Create(req.Context(), &org)
		if err != nil {
			l.Error("failed to store organization", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully stored organization")
		utils.JSONResponse(w, http.StatusOK, result)
	}
}

func (h *Handler) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := h.Create()
	handler(response.ResponseWriter, request.Request)
}
