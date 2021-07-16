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

func (s *Server) AuthClients(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	if req.Method == "POST" {
		s.PostAuthClient(w, req, true, &models.OAuth2Client{})
		return
	}

	clients, err := s.HydraAdmin.ListOAuth2Clients(&admin.ListOAuth2ClientsParams{
		Context: ctx,
	})
	if err != nil {
		s.Error(w, err)
		return
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "auth_clients", map[string]interface{}{
		"Clients": clients.Payload,
	}); err != nil {
		s.Error(w, err)
		return
	}

}

func (s *Server) AuthClientNewSecret(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		s.Error(w, err)
		return
	}

	cli, err := s.HydraAdmin.GetOAuth2Client(&admin.GetOAuth2ClientParams{
		ID:      id,
		Context: ctx,
	})
	if err != nil {
		s.Error(w, err)
		return
	}

	clientSecret, err := GenerateRandomStringURLSafe(32)
	if err != nil {
		s.Error(w, err)
		return
	}

	c := cli.Payload
	c.ClientSecret = clientSecret

	_, err = s.HydraAdmin.UpdateOAuth2Client(&admin.UpdateOAuth2ClientParams{
		Context: ctx,
		Body:    c,
		ID:      c.ClientID,
	})
	if err != nil {
		s.Error(w, err)
		return
	}
	s.sessionManager.Put(ctx, "client_secret", clientSecret)

	http.Redirect(w, req, "/settings/authclients/"+id, http.StatusSeeOther)

}

func (s *Server) AuthClient(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		s.Error(w, err)
		return
	}

	if id == "new" {
		if err := s.renderFactory.New(req).ExecuteTemplate(w, "auth_client", map[string]interface{}{
			"GrantTypes":    map[string]bool{},
			"ResponseTypes": map[string]bool{},
			"AuthMethod":    "none",
		}); err != nil {
			s.Error(w, err)
			return
		}
		return
	}

	cli, err := s.HydraAdmin.GetOAuth2Client(&admin.GetOAuth2ClientParams{
		ID:      id,
		Context: ctx,
	})
	if err != nil {
		s.Error(w, err)
		return
	}

	if req.Method == "POST" {
		s.PostAuthClient(w, req, false, cli.Payload)
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

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "auth_client", map[string]interface{}{
		"Client":        cli.Payload,
		"ClientSecret":  s.sessionManager.PopString(ctx, "client_secret"),
		"GrantTypes":    grantTypes,
		"ResponseTypes": responseTypes,
		"AuthMethod":    cli.Payload.TokenEndpointAuthMethod,
	}); err != nil {
		s.Error(w, err)
		return
	}

}

func (s *Server) DeleteAuthClient(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		s.Error(w, err)
		return
	}

	_, err := s.HydraAdmin.DeleteOAuth2Client(&admin.DeleteOAuth2ClientParams{
		ID:      id,
		Context: ctx,
	})
	if err != nil {
		s.Error(w, err)
		return
	}

	http.Redirect(w, req, "/settings/authclients", http.StatusSeeOther)

}

func (s *Server) PostAuthClient(w http.ResponseWriter, req *http.Request, isNew bool, cli *models.OAuth2Client) {

	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		s.Error(w, err)
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
			s.Error(w, err)
			return
		}
		s.sessionManager.Put(ctx, "client_secret", clientSecret)
		cli.ClientSecret = clientSecret
		cli.ClientID = uuid.NewV4().String()
		response, err := s.HydraAdmin.CreateOAuth2Client(&admin.CreateOAuth2ClientParams{
			Context: ctx,
			Body:    cli,
		})
		if err != nil {
			s.Error(w, err)
			return
		}
		http.Redirect(w, req, "/settings/authclients/"+response.Payload.ClientID, http.StatusSeeOther)
		return
	} else {
		response, err := s.HydraAdmin.UpdateOAuth2Client(&admin.UpdateOAuth2ClientParams{
			Body:    cli,
			ID:      cli.ClientID,
			Context: ctx,
		})
		if err != nil {
			s.Error(w, err)
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
