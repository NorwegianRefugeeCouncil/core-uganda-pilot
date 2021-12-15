package organization

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) Get(organizationID string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		l := logging.NewLogger(req.Context()).With(zap.String("organization_id", organizationID))

		l.Debug("retrieving organization from store")
		result, err := h.store.Get(req.Context(), organizationID)
		if err != nil {
			l.Error("failed to retrieve organization from store", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully retrieved organization from store")
		utils.JSONResponse(w, http.StatusOK, result)

	}
}

func (h *Handler) RestfulGet(request *restful.Request, response *restful.Response) {
	organizationID := request.PathParameter(constants.ParamOrganizationID)
	handler := h.Get(organizationID)
	handler(response.ResponseWriter, request.Request)
}
