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

func restfulList(hydraAdmin admin.ClientService) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handleList(hydraAdmin)(res.ResponseWriter, req.Request)
	}
}

func handleList(hydraAdmin admin.ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		l := logging.NewLogger(req.Context())

		l.Debug("listing clients")
		resp, err := hydraAdmin.ListOAuth2Clients(&admin.ListOAuth2ClientsParams{
			Context: req.Context(),
		})
		if err != nil {
			l.Error("failed to list clients", zap.Error(err))
			utils.ErrorResponse(w, meta.NewInternalServerError(err))
			return
		}

		l.Debug("successfully listed clients")
		utils.JSONResponse(w, http.StatusOK, mapFromHydraClients(resp.Payload))
	}
}
