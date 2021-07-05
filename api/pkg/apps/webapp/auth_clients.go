package webapp

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
)

func (h *Server) AuthClients(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	if req.Method == "POST" {
		h.PostAuthClient(w, req, true, &models.OAuth2Client{})
		return
	}

	clients, err := h.HydraAdmin.ListOAuth2Clients(&admin.ListOAuth2ClientsParams{
		Context: ctx,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "auth_clients", map[string]interface{}{
		"Clients": clients.Payload,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Server) AuthClientNewSecret(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cli, err := h.HydraAdmin.GetOAuth2Client(&admin.GetOAuth2ClientParams{
		ID:      id,
		Context: ctx,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clientSecret, err := GenerateRandomStringURLSafe(32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := cli.Payload
	c.ClientSecret = clientSecret

	_, err = h.HydraAdmin.UpdateOAuth2Client(&admin.UpdateOAuth2ClientParams{
		Context: ctx,
		Body:    c,
		ID:      c.ClientID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.sessionManager.Put(ctx, "client_secret", clientSecret)

	http.Redirect(w, req, "/settings/authclients/"+id, http.StatusSeeOther)

}

func (h *Server) AuthClient(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if id == "new" {
		if err := h.renderFactory.New(req).ExecuteTemplate(w, "auth_client", map[string]interface{}{
			"GrantTypes":    map[string]bool{},
			"ResponseTypes": map[string]bool{},
			"AuthMethod":    "none",
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	cli, err := h.HydraAdmin.GetOAuth2Client(&admin.GetOAuth2ClientParams{
		ID:      id,
		Context: ctx,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostAuthClient(w, req, false, cli.Payload)
		return
	}

	grantTypes := map[string]bool{}
	for _, grantType := range cli.Payload.GrantTypes {
		grantTypes[grantType] = true
	}

	responseTypes := map[string]bool{}
	for _, responseType := range cli.Payload.ResponseTypes {
		responseTypes[responseType] = true
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "auth_client", map[string]interface{}{
		"Client":        cli.Payload,
		"ClientSecret":  h.sessionManager.PopString(ctx, "client_secret"),
		"GrantTypes":    grantTypes,
		"ResponseTypes": responseTypes,
		"AuthMethod":    cli.Payload.TokenEndpointAuthMethod,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Server) DeleteAuthClient(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err := h.HydraAdmin.DeleteOAuth2Client(&admin.DeleteOAuth2ClientParams{
		ID:      id,
		Context: ctx,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/settings/authclients", http.StatusSeeOther)

}

func (h *Server) PostAuthClient(w http.ResponseWriter, req *http.Request, isNew bool, cli *models.OAuth2Client) {

	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	values := req.Form

	grantTypes, responseTypes := getFormValues(values)

	cli.ClientName = values.Get("client_name")
	cli.GrantTypes = grantTypes
	cli.ResponseTypes = responseTypes
	cli.RedirectUris = values["redirect_uris"]
	cli.TokenEndpointAuthMethod = values.Get("auth_method")

	if isNew {
		cli = &*cli
		clientSecret, err := GenerateRandomStringURLSafe(32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.sessionManager.Put(ctx, "client_secret", clientSecret)
		cli.ClientSecret = clientSecret
		cli.ClientID = uuid.NewV4().String()
		response, err := h.HydraAdmin.CreateOAuth2Client(&admin.CreateOAuth2ClientParams{
			Context: ctx,
			Body:    cli,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, "/settings/authclients/"+response.Payload.ClientID, http.StatusSeeOther)
		return
	} else {
		response, err := h.HydraAdmin.UpdateOAuth2Client(&admin.UpdateOAuth2ClientParams{
			Body:    cli,
			ID:      cli.ClientID,
			Context: ctx,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, "/settings/authclients/"+response.Payload.ClientID, http.StatusSeeOther)
		return
	}

}

func getFormValues(values url.Values) ([]string, []string) {

	var grantTypesFieldNames = []string{
		"authorization_code",
		"code",
		"refresh_token",
		"id_token",
		"implicit",
		"client_credentials",
	}

	var grantTypes []string

	for _, fieldName := range grantTypesFieldNames {
		if values.Get(fmt.Sprintf("grant_types[%s]", fieldName)) == "true" {
			grantTypes = append(grantTypes, "fieldName")
		}
	}

	var responseTypesFieldNames = []string{
		"token",
		"code",
		"id_token",
	}

	var responseTypes []string

	for _, fieldName := range responseTypesFieldNames {
		if values.Get(fmt.Sprintf("response_types[%s]", fieldName)) == "true" {
			responseTypes = append(responseTypes, "fieldName")
		}
	}

	return grantTypes, responseTypes
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomStringURLSafe returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomStringURLSafe(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}
