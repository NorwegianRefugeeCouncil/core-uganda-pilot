package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nrc-no/core-kafka/pkg/keycloak"
	"github.com/nrc-no/core-kafka/pkg/sessionmanager"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
)

const IdTokenKey = "id_token"
const AccessTokenKey = "access_token"

type Profile struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type Handler struct {
	Store          sessionmanager.Store
	Provider       *oidc.Provider
	Config         oauth2.Config
	Verifier       *oidc.IDTokenVerifier
	KeycloakClient *keycloak.Client
}

func init() {
	// needed so that we can store these types of values
	// in the session store.
	gob.Register(map[string]interface{}{})
	gob.Register(&Profile{})
}

func NewHandler(
	ctx context.Context,
	issuerURL, clientID, clientSecret, redirectURL string,
	store sessionmanager.Store,
	keycloakClient *keycloak.Client,
) (*Handler, error) {
	l := logrus.WithField("logger", "auth.NewHandler")
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		err := fmt.Errorf("failed to get oidc provider: %v", err)
		l.WithError(err).Errorf("")
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	verifier := provider.Verifier(oidcConfig)

	return &Handler{
		Provider:       provider,
		Config:         oauth2Config,
		Verifier:       verifier,
		Store:          store,
		KeycloakClient: keycloakClient,
	}, nil
}

func (h *Handler) ClearSession(ctx context.Context) {
	h.Store.Destroy(ctx)
}

func (h *Handler) Callback(w http.ResponseWriter, req *http.Request) {

	l := logrus.WithField("logger", "auth.Callback")

	ctx := req.Context()

	// retrieve the state query parameter, compare it to the session
	// query parameter. The two should match

	if req.URL.Query().Get("state") != h.Store.PopString(ctx, "state") {
		err := fmt.Errorf("invalid state parameter")
		l.WithError(err).Warnf("")
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.ClearSession(ctx)
		return
	}

	// Perform code exchange to get oauth token
	code := req.URL.Query().Get("code")
	oauth2Token, err := h.Config.Exchange(ctx, code)
	if err != nil {
		err = fmt.Errorf("failed to exchange token: %v", err)
		l.WithError(err).Warnf("")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		h.ClearSession(ctx)
		return
	}

	// Retrieve ID Token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		err := fmt.Errorf("missing id token")
		l.WithError(err).Warnf("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.ClearSession(ctx)
		return
	}

	// Verify ID token
	idToken, err := h.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		err := fmt.Errorf("failed to verify ID token: %v", err)
		l.WithError(err).Warnf("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.ClearSession(ctx)
		return
	}

	// Extract profile from claims
	var profile Profile
	if err := idToken.Claims(&profile); err != nil {
		err := fmt.Errorf("failed to unmarshal id token claims: %v", err)
		l.WithError(err).Warnf("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.ClearSession(ctx)
		return
	}

	// Populate session with ID & Access tokens, profile
	h.Store.Put(ctx, IdTokenKey, rawIDToken)
	h.Store.Put(ctx, AccessTokenKey, oauth2Token.AccessToken)

	// Redirect to some page
	// TODO: perhaps redirect to the initial requested page?
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Header().Set("Location", "/individuals")

}

func (h *Handler) Login(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	l := logrus.WithField("logger", "auth.Login")

	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		err := fmt.Errorf("failed to create random state: %v", err)
		l.WithError(err).Errorf("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.ClearSession(ctx)
		return
	}

	state := base64.StdEncoding.EncodeToString(b)

	h.Store.Put(ctx, "state", state)

	http.Redirect(w, req, h.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)

}

func (h *Handler) Logout(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	tokenURL := h.Provider.Endpoint().TokenURL
	logoutURLStr := tokenURL[0:len(tokenURL)-6] + "/logout"

	l := logrus.WithField("logger", "logout")

	if err := h.Store.Destroy(ctx); err != nil {
		l.WithError(err).Warnf("failed to clear session")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logoutURL, err := url.Parse(logoutURLStr)
	if err != nil {
		l.WithError(err).Warnf("failed to parse logout url")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.ClearSession(ctx)
		return
	}

	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}
	returnTo, err := url.Parse(fmt.Sprintf("%s://%s", scheme, req.Host))
	if err != nil {
		l.WithError(err).Warnf("failed to parse return url")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.ClearSession(ctx)
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", h.Config.ClientID)
	parameters.Add("client_secret", h.Config.ClientSecret)
	parameters.Add("redirect_uri", returnTo.String())
	logoutURL.RawQuery = parameters.Encode()

	http.Redirect(w, req, logoutURL.String(), http.StatusTemporaryRedirect)

}
