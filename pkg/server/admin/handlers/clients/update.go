package clients

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"go.uber.org/zap"
	"net/http"
)

func restfulUpdate(hydraAdmin admin.ClientService) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		clientID := req.PathParameter("clientId")
		handleUpdate(hydraAdmin, clientID)(res.ResponseWriter, req.Request)
	}
}

func handleUpdate(hydraAdmin admin.ClientService, clientID string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		l := logging.NewLogger(req.Context()).With(zap.String("client_id", clientID))

		l.Debug("unmarshaling client")
		var client types.Oauth2Client
		if err := utils.BindJSON(req, &client); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("updating hydra client")
		resp, err := hydraAdmin.UpdateOAuth2Client(&admin.UpdateOAuth2ClientParams{
			ID:         clientID,
			Body:       mapToHydraClient(client),
			Context:    req.Context(),
			HTTPClient: nil,
		})
		if err != nil {
			l.Error("failed to update hydra client", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully updated hydra client")
		utils.JSONResponse(w, http.StatusOK, mapFromHydraClient(resp.Payload))
	}
}
