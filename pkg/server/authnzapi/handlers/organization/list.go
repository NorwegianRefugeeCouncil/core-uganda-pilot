package organization

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		l := logging.NewLogger(req.Context())

		l.Debug("retrieving organizations from store")
		result, err := h.store.List(req.Context())
		if err != nil {
			l.Error("failed to retrieve organizations from store", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully retrieved organizations from store", zap.Int("count", len(result)))
		utils.JSONResponse(w, http.StatusOK, &types.OrganizationList{
			Items: result,
		})
	}
}

func (h *Handler) RestfulList(request *restful.Request, response *restful.Response) {
	handler := h.List()
	handler(response.ResponseWriter, request.Request)
}
