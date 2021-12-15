package clients

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"go.uber.org/zap"
	"net/http"
)

func restfulGet(hydraAdmin admin.ClientService) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		clientID := req.PathParameter("clientId")
		handleGet(hydraAdmin, clientID)(res.ResponseWriter, req.Request)
	}
}

func handleGet(hydraAdmin admin.ClientService, clientID string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		l := logging.NewLogger(req.Context()).With(zap.String("client_id", clientID))

		l.Debug("getting oauth2 client")
		oauth2Client, err := hydraAdmin.GetOAuth2Client(&admin.GetOAuth2ClientParams{
			Context: req.Context(),
			ID:      clientID,
		})
		if err != nil {
			l.Error("failed to get hydra client", zap.Error(err))
			utils.ErrorResponse(w, meta.NewInternalServerError(err))
			return
		}

		l.Debug("successfully got hydra client")
		utils.JSONResponse(w, http.StatusOK, mapFromHydraClient(oauth2Client.Payload))
	}
}
