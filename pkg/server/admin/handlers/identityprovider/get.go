package identityprovider

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) Get(identityProviderId string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		l := logging.NewLogger(req.Context()).With(zap.String("identity_provider_id", identityProviderId))

		l.Debug("validating identity provider id")
		if _, err := uuid.FromString(identityProviderId); err != nil {
			l.Error("failed to validate identity provider id", zap.Error(err))
			utils.ErrorResponse(w, meta.NewBadRequest("invalid identity provider id"))
			return
		}

		l.Debug("getting identity provider")
		identityProvider, err := h.store.Get(req.Context(), identityProviderId, store.IdentityProviderGetOptions{
			ReturnClientSecret: false,
		})
		if err != nil {
			l.Error("failed to get identity provider", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully got identity provider")
		utils.JSONResponse(w, http.StatusOK, identityProvider)
	}
}

func (h *Handler) RestfulGet(request *restful.Request, response *restful.Response) {
	identityProviderId := request.PathParameter(constants.ParamIdentityProviderID)
	handler := h.Get(identityProviderId)
	handler(response.ResponseWriter, request.Request)
}
