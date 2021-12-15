package clients

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"go.uber.org/zap"
	"net/http"
)

func restfulDelete(hydraAdmin admin.ClientService) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		clientID := req.PathParameter("clientId")
		handleDelete(hydraAdmin, clientID)(res.ResponseWriter, req.Request)
	}
}

func handleDelete(hydraAdmin admin.ClientService, clientID string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		l := logging.NewLogger(req.Context()).With(zap.String("client_id", clientID))

		l.Debug("deleting hydra client")
		_, err := hydraAdmin.DeleteOAuth2Client(&admin.DeleteOAuth2ClientParams{
			ID:         clientID,
			Context:    req.Context(),
			HTTPClient: nil,
		})
		if err != nil {
			l.Error("failed to delete hydra client", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully deleted hydra client")
		w.WriteHeader(http.StatusNoContent)
	}
}
