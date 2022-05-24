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
		l := logging.NewLogger(req.Context()).With(zap.String("identity_id", identityId))

		l.Debug("validating identity id")
		if _, err := uuid.FromString(identityId); err != nil {
			l.Error("failed to validate identity id", zap.Error(err))
			utils.ErrorResponse(w, meta.NewBadRequest("invalid identity id"))
			return
		}

		l.Debug("getting identity")
		identity, err := h.store.Get(req.Context(), identityId)
		if err != nil {
			l.Error("failed to get identity", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully got identity")
		utils.JSONResponse(w, http.StatusOK, identity)
	}
}

func (h *Handler) RestfulGet(request *restful.Request, response *restful.Response) {
	identityId := request.PathParameter(constants.ParamIdentityID)
	handler := h.Get(identityId)
	handler(response.ResponseWriter, request.Request)
}
