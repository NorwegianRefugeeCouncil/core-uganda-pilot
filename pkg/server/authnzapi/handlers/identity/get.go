package identity

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

func (h *Handler) Get(identityId string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		l := logging.NewLogger(req.Context()).With(zap.String("identity_provider_id", identityId))

		l.Debug("validating identity provider id")
		if _, err := uuid.FromString(identityId); err != nil {
			l.Error("failed to validate identity provider id", zap.Error(err))
			utils.ErrorResponse(w, meta.NewBadRequest("invalid identity provider id"))
			return
		}

		l.Debug("getting identity provider")
		identity, err := h.store.Get(req.Context(), identityId)
		if err != nil {
			l.Error("failed to get identity provider", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully got identity provider")
		utils.JSONResponse(w, http.StatusOK, identity)
	}
}

func (h *Handler) RestfulGet(request *restful.Request, response *restful.Response) {
	identityId := request.PathParameter(constants.ParamIdentityID)
	handler := h.Get(identityId)
	handler(response.ResponseWriter, request.Request)
}
