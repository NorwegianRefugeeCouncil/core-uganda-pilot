package clients

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
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
		oauth2Client, err := hydraAdmin.GetOAuth2Client(&admin.GetOAuth2ClientParams{
			Context: req.Context(),
			ID:      clientID,
		})
		if err != nil {
			utils.ErrorResponse(w, meta.NewInternalServerError(err))
			return
		}

		if err != nil {
			utils.ErrorResponse(w, meta.NewInternalServerError(err))
			return
		}

		w.Header().Set("Content-Type", "text/html")
		renderClient(w, oauth2Client.Payload, "", false, "")
	}
}
