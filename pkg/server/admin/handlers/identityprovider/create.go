package identityprovider

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		l := logging.NewLogger(req.Context())

		l.Debug("unmarshaling identity provider")
		var idp types.IdentityProvider
		if err := utils.BindJSON(req, &idp); err != nil {
			l.Error("failed to unmarshal identity provider", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("storing identity provider")
		res, err := h.store.Create(req.Context(), &idp, store.IdentityProviderCreateOptions{
			ReturnClientSecret: false,
		})
		if err != nil {
			l.Error("failed to store identity provider", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully created identity provider")
		utils.JSONResponse(w, http.StatusOK, res)
	}
}

func (h *Handler) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := h.Create()
	handler(response.ResponseWriter, request.Request)
}
