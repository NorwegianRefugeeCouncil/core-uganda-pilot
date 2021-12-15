package identityprovider

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) Update(identityProviderID string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		l := logging.NewLogger(req.Context()).With(zap.String("identity_provider_id", identityProviderID))

		l.Debug("validating identity provider id")
		if _, err := uuid.FromString(identityProviderID); err != nil {
			l.Error("failed to validate identity provider id", zap.Error(err))
			utils.ErrorResponse(w, meta.NewBadRequest("invalid identity provider id"))
			return
		}

		l.Debug("unmarshaling identity provider")
		var identityProvider types.IdentityProvider
		if err := utils.BindJSON(req, &identityProvider); err != nil {
			l.Error("failed to unmarshal identity provider", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("ensuring identity provider id from path and from body do match")
		if len(identityProvider.ID) > 0 && identityProvider.ID != identityProviderID {
			l.Error("failed to verify match on identity provider id")
			utils.ErrorResponse(w, meta.NewBadRequest("identity provider id mismatch"))
			return
		}

		l.Debug("storing new identity provider")
		res, err := h.store.Update(req.Context(), &identityProvider, store.IdentityProviderUpdateOptions{
			ReturnClientSecret: false,
		})
		if err != nil {
			l.Error("failed to store new identity provider", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully updated identity provider")
		utils.JSONResponse(w, http.StatusOK, res)
	}
}

func (h *Handler) RestfulUpdate(request *restful.Request, response *restful.Response) {
	identityProviderID := request.PathParameter(constants.ParamIdentityProviderID)
	handler := h.Update(identityProviderID)
	handler(response.ResponseWriter, request.Request)
}
