package clients

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"net/http"
)

func restfulCreate(hydraAdmin admin.ClientService) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handleCreate(hydraAdmin)(res.ResponseWriter, req.Request)
	}
}

func handleCreate(hydraAdmin admin.ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var client types.Oauth2Client
		if err := utils.BindJSON(req, &client); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		resp, err := hydraAdmin.CreateOAuth2Client(&admin.CreateOAuth2ClientParams{
			Body:       mapToHydraClient(client),
			Context:    req.Context(),
			HTTPClient: nil,
		})
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		utils.JSONResponse(w, http.StatusOK, mapFromHydraClient(resp.Payload))
	}
}
