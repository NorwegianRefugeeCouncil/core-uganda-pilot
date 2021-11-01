package identityprovider

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (h *Handler) List(organizationID string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		if _, err := uuid.FromString(organizationID); err != nil {
			utils.ErrorResponse(w, meta.NewBadRequest("invalid organization id"))
			return
		}

		identityProviders, err := h.store.List(req.Context(), organizationID, store.IdentityProviderListOptions{
			ReturnClientSecret: false,
		})

		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		list := &types.IdentityProviderList{Items: identityProviders}
		utils.JSONResponse(w, http.StatusOK, list)
	}
}

func (h *Handler) RestfulList(request *restful.Request, response *restful.Response) {
	organizationID := request.QueryParameter(constants.ParamOrganizationID)
	handler := h.List(organizationID)
	handler(response.ResponseWriter, request.Request)
}
