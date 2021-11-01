package clients

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/server/admin/templates"
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

		oauth2List, err := hydraAdmin.ListOAuth2Clients(&admin.ListOAuth2ClientsParams{
			Context: req.Context(),
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
		templates.Template.ExecuteTemplate(w, "client_list", map[string]interface{}{
			"Clients": oauth2List.Payload,
		})
	}
}
