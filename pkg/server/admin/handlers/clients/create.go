package clients

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/server/admin/templates"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"net/http"
	"strings"
)

func restfulCreate(hydraAdmin admin.ClientService) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handleCreate(hydraAdmin)(res.ResponseWriter, req.Request)
	}
}

func handleCreate(hydraAdmin admin.ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		w.Header().Set("Content-Type", "text/html")

		if err := req.ParseForm(); err != nil {
			templates.Template.ExecuteTemplate(w, "client_add", map[string]interface{}{
				"Error": err.Error(),
			})
			return
		}

		q := req.Form

		clientID := q.Get("client_id")
		clientName := q.Get("client_name")
		clientURI := q.Get("client_uri")
		grantTypes := strings.Split(q.Get("grant_types"), ",")
		responseTypes := strings.Split(q.Get("response_types"), ",")
		scope := q.Get("scope")
		redirectUris := strings.Split(q.Get("redirect_uris"), ",")
		tokenEndpointAuthMethod := q.Get("token_endpoint_auth_method")
		allowedCorsOrigins := strings.Split(q.Get("allowed_cors_origins"), ",")

		if len(allowedCorsOrigins) == 1 && allowedCorsOrigins[0] == "" {
			allowedCorsOrigins = nil
		}
		if len(grantTypes) == 1 && grantTypes[0] == "" {
			grantTypes = nil
		}
		if len(redirectUris) == 1 && redirectUris[0] == "" {
			redirectUris = nil
		}
		if len(responseTypes) == 1 && responseTypes[0] == "" {
			responseTypes = nil
		}

		resp, err := hydraAdmin.CreateOAuth2Client(&admin.CreateOAuth2ClientParams{
			Body: &models.OAuth2Client{
				AllowedCorsOrigins:      allowedCorsOrigins,
				ClientID:                clientID,
				ClientName:              clientName,
				ClientURI:               clientURI,
				GrantTypes:              grantTypes,
				RedirectUris:            redirectUris,
				ResponseTypes:           responseTypes,
				Scope:                   scope,
				TokenEndpointAuthMethod: tokenEndpointAuthMethod,
			},
			Context:    req.Context(),
			HTTPClient: nil,
		})
		if err != nil {
			templates.Template.ExecuteTemplate(w, "client_add", map[string]interface{}{
				"Error":                   err.Error(),
				"ClientID":                clientID,
				"ClientName":              clientName,
				"ClientURI":               clientURI,
				"GrantTypes":              strings.Join(grantTypes, ","),
				"ResponseTypes":           strings.Join(responseTypes, ","),
				"Scope":                   scope,
				"RedirectURIs":            strings.Join(redirectUris, ","),
				"AllowedCORSOrigins":      strings.Join(allowedCorsOrigins, ","),
				"TokenEndpointAuthMethod": tokenEndpointAuthMethod,
			})
			return
		}

		payload := resp.Payload
		clientSecret := payload.ClientSecret

		renderClient(w, payload, clientSecret, false, "")
	}
}

func renderClient(
	w http.ResponseWriter,
	payload *models.OAuth2Client,
	clientSecret string,
	isNew bool,
	error string,
) error {
	return templates.Template.ExecuteTemplate(w, "client_add", map[string]interface{}{
		"ClientID":                payload.ClientID,
		"ClientName":              payload.ClientName,
		"ClientURI":               payload.ClientURI,
		"GrantTypes":              strings.Join(payload.GrantTypes, ","),
		"ResponseTypes":           strings.Join(payload.ResponseTypes, ","),
		"Scope":                   payload.Scope,
		"RedirectURIs":            strings.Join(payload.RedirectUris, ","),
		"AllowedCORSOrigins":      strings.Join(payload.AllowedCorsOrigins, ","),
		"TokenEndpointAuthMethod": payload.TokenEndpointAuthMethod,
		"ClientSecret":            clientSecret,
		"IsNew":                   isNew,
		"Error":                   error,
	})
}
