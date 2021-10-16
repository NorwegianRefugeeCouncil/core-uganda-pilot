package identityprovider

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/bla/constants"
	"github.com/nrc-no/core/pkg/bla/store"
	"github.com/nrc-no/core/pkg/bla/types"
	"github.com/nrc-no/core/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (h *Handler) Update(identityProviderID string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		if _, err := uuid.FromString(identityProviderID); err != nil {
			utils.ErrorResponse(w, meta.NewBadRequest("invalid identity provider id"))
			return
		}

		var identityProvider types.IdentityProvider
		if err := utils.BindJSON(req, &identityProvider); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		if len(identityProvider.ID) > 0 && identityProvider.ID != identityProviderID {
			utils.ErrorResponse(w, meta.NewBadRequest("identity provider id mismatch"))
			return
		}

		res, err := h.store.Update(req.Context(), &identityProvider, store.IdentityProviderUpdateOptions{
			ReturnClientSecret: false,
		})
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		utils.JSONResponse(w, http.StatusOK, res)
	}
}

func (h *Handler) RestfulUpdate(request *restful.Request, response *restful.Response) {
	recordID := request.PathParameter(constants.ParamRecordID)
	handler := h.Update(recordID)
	handler(response.ResponseWriter, request.Request)
}
