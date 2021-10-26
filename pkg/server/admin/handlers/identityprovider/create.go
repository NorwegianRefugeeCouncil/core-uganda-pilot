package identityprovider

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/types"
	"github.com/nrc-no/core/pkg/utils"
	"net/http"
)

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var idp types.IdentityProvider
		if err := utils.BindJSON(req, &idp); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		res, err := h.store.Create(req.Context(), &idp, store.IdentityProviderCreateOptions{
			ReturnClientSecret: false,
		})
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		utils.JSONResponse(w, http.StatusOK, res)
	}
}

func (h *Handler) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := h.Create()
	handler(response.ResponseWriter, request.Request)
}
