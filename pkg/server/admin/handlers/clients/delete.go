package clients

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
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
		_, err := hydraAdmin.DeleteOAuth2Client(&admin.DeleteOAuth2ClientParams{
			ID:         clientID,
			Context:    req.Context(),
			HTTPClient: nil,
		})
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
