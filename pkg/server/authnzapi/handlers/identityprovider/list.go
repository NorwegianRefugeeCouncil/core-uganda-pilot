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

func (h *Handler) List(organizationID string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		l := logging.NewLogger(req.Context()).With(zap.String("organization_id", organizationID))

		l.Debug("validating organization id")
		if _, err := uuid.FromString(organizationID); err != nil {
			l.Error("failed to validate organization id", zap.Error(err))
			utils.ErrorResponse(w, meta.NewBadRequest("invalid organization id"))
			return
		}

		l.Debug("retrieving identity providers from store")
		identityProviders, err := h.store.List(req.Context(), organizationID, store.IdentityProviderListOptions{
			ReturnClientSecret: false,
		})
		if err != nil {
			l.Error("failed to retrieve identity providers from store", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully listed identity providers", zap.Int("count", len(identityProviders)))
		list := &types.IdentityProviderList{Items: identityProviders}
		utils.JSONResponse(w, http.StatusOK, list)
	}
}

func (h *Handler) RestfulList(request *restful.Request, response *restful.Response) {
	organizationID := request.QueryParameter(constants.ParamOrganizationID)
	handler := h.List(organizationID)
	handler(response.ResponseWriter, request.Request)
}
