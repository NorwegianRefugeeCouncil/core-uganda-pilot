package clients

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"net/http"
)

func restfulList(hydraAdmin admin.ClientService) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handleList(hydraAdmin)(res.ResponseWriter, req.Request)
	}
}

func handleList(hydraAdmin admin.ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp, err := hydraAdmin.ListOAuth2Clients(&admin.ListOAuth2ClientsParams{
			Context: req.Context(),
		})
		if err != nil {
			utils.ErrorResponse(w, meta.NewInternalServerError(err))
			return
		}
		utils.JSONResponse(w, http.StatusOK, mapFromHydraClients(resp.Payload))
	}
}
