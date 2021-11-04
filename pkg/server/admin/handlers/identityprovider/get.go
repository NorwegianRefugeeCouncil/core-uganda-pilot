package identityprovider

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (h *Handler) Get(identityProviderId string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		if _, err := uuid.FromString(identityProviderId); err != nil {
			utils.ErrorResponse(w, meta.NewBadRequest("invalid identity provider id"))
			return
		}

		identityProvider, err := h.store.Get(req.Context(), identityProviderId, store.IdentityProviderGetOptions{
			ReturnClientSecret: false,
		})

		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		utils.JSONResponse(w, http.StatusOK, identityProvider)
	}
}

func (h *Handler) RestfulGet(request *restful.Request, response *restful.Response) {
	identityProviderId := request.PathParameter(constants.ParamIdentityProviderID)
	handler := h.Get(identityProviderId)
	handler(response.ResponseWriter, request.Request)
}
