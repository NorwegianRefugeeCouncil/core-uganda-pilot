package clients

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"net/http"
)

func restfulUpdate(hydraAdmin admin.ClientService) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		clientID := req.PathParameter("clientId")
		handleUpdate(hydraAdmin, clientID)(res.ResponseWriter, req.Request)
	}
}

type Client struct {
	ID                      string   `json:"id"`
	ClientName              string   `json:"clientName"`
	URI                     string   `json:"uri"`
	GrantTypes              []string `json:"grantTypes"`
	ResponseTypes           []string `json:"responseTypes"`
	Scope                   string   `json:"scope"`
	RedirectURIs            []string `json:"redirectUris"`
	AllowedCorsOrigins      []string `json:"allowedCorsOrigins"`
	TokenEndpointAuthMethod string   `json:"tokenEndpointAuthMethod"`
}

type ClientList struct {
	Items []*Client `json:"items"`
}

func handleUpdate(hydraAdmin admin.ClientService, clientID string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var client Client
		if err := utils.BindJSON(req, &client); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		resp, err := hydraAdmin.UpdateOAuth2Client(&admin.UpdateOAuth2ClientParams{
			ID:         clientID,
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

func mapToHydraClient(client Client) *models.OAuth2Client {
	return &models.OAuth2Client{
		ClientID:                client.ID,
		AllowedCorsOrigins:      client.AllowedCorsOrigins,
		ClientName:              client.ClientName,
		ClientURI:               client.URI,
		GrantTypes:              client.GrantTypes,
		RedirectUris:            client.RedirectURIs,
		ResponseTypes:           client.ResponseTypes,
		Scope:                   client.Scope,
		TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
	}
}

func mapFromHydraClient(client *models.OAuth2Client) *Client {
	return &Client{
		ID:                      client.ClientID,
		AllowedCorsOrigins:      client.AllowedCorsOrigins,
		ClientName:              client.ClientName,
		URI:                     client.ClientURI,
		GrantTypes:              client.GrantTypes,
		RedirectURIs:            client.RedirectUris,
		ResponseTypes:           client.ResponseTypes,
		Scope:                   client.Scope,
		TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
	}
}
