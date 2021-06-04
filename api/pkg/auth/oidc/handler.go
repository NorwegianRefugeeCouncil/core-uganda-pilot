package oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/emicklei/go-restful"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gopkg.in/boj/redistore.v1"
	"io"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
	"net/url"
)

func init() {
	gob.Register(map[string]interface{}{})
}

type OIDCHandler struct {
	codec      runtime.NegotiatedSerializer
	getConfig  func(req *http.Request) oauth2.Config
	provider   *oidc.Provider
	oidcConfig *oidc.Config
	verifier   *oidc.IDTokenVerifier
	store      sessions.Store
	issuerUrl  string
	clientId   string
}

func NewOIDCHandler(
	codec runtime.NegotiatedSerializer,
	clientId string,
	clientSecret string,
	issuerURL string,
) *OIDCHandler {

	provider, err := oidc.NewProvider(context.Background(), issuerURL)
	if err != nil {
		logrus.WithError(err).Fatal("could not create oidc provider")
	}

	oidcConfig := &oidc.Config{
		ClientID: clientId,
	}

	verifier := provider.Verifier(oidcConfig)

	// TODO: allow users to specify key using flag
	cookieKey, err := randString(32)
	if err != nil {
		panic(err)
	}

	store, err := redistore.NewRediStore(10, "tcp", ":6379", "", []byte(cookieKey))
	if err != nil {
		panic(err)
	}

	return &OIDCHandler{
		store:      store,
		codec:      codec,
		provider:   provider,
		oidcConfig: oidcConfig,
		verifier:   verifier,
		issuerUrl:  issuerURL,
		clientId:   clientId,
		getConfig: func(req *http.Request) oauth2.Config {

			scheme := "http://"
			if req.TLS != nil {
				scheme = "https://"
			}

			return oauth2.Config{
				ClientID:     clientId,
				ClientSecret: clientSecret,
				Endpoint:     provider.Endpoint(),
				RedirectURL:  fmt.Sprintf("%s%s/auth/callback", scheme, req.Host),
				Scopes: []string{
					oidc.ScopeOpenID,
					"profile",
					"email",
					"roles",
				},
			}
		},
	}
}

func (h *OIDCHandler) WebService() *restful.WebService {
	ws := new(restful.WebService)

	loginRoute := ws.GET("/auth/login").To(h.restfulLogin()).
		Doc("login using oidc provider").
		Operation("oidcLogin").
		Param(ws.QueryParameter("redirectUrl", "redirects to this url after successful login").Required(false)).
		Returns(http.StatusTemporaryRedirect, "OK", nil)

	ws.Route(loginRoute)

	logoutRoute := ws.GET("/auth/logout").To(h.restfulLogout()).
		Doc("logout using oidc provider").
		Operation("oidcLogout").
		Returns(http.StatusSeeOther, "OK", nil)

	ws.Route(logoutRoute)

	callbackRoute := ws.GET("/auth/callback").To(h.restfulCallback()).
		Doc("logout using oidc provider").
		Operation("oidcLogout").
		Param(ws.QueryParameter("code", "oidc authorization code").Required(true)).
		Param(ws.QueryParameter("state", "oidc state").Required(true)).
		Returns(http.StatusTemporaryRedirect, "OK", nil)

	ws.Route(callbackRoute)

	return ws
}

func (h *OIDCHandler) restfulLogin() restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		h.ServeLogin(response.ResponseWriter, request.Request)
	}
}
func (h *OIDCHandler) restfulLogout() restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		h.ServeLogout(response.ResponseWriter, request.Request)
	}
}
func (h *OIDCHandler) restfulCallback() restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		h.ServeCallback(response.ResponseWriter, request.Request)
	}
}

func (h *OIDCHandler) ServeLogin(w http.ResponseWriter, req *http.Request) {

	// Retrieve the session
	session, _ := h.store.Get(req, "auth-session")

	// Create a random state string
	stateBytes, err := randString(32)
	if err != nil {
		logrus.WithError(err).Error("unable to get random state")
		h.err(apierrors.NewInternalError(errors.New("internal error")), w, req)
		return
	}
	state := base64.StdEncoding.EncodeToString(stateBytes)

	// Create a random nonce string
	nonceBytes, err := randString(32)
	if err != nil {
		logrus.WithError(err).Error("unable to get random state")
		h.err(apierrors.NewInternalError(errors.New("internal error")), w, req)
		return
	}
	nonce := base64.StdEncoding.EncodeToString(nonceBytes)

	// Perhaps we want to redirect to a given url
	// after a successful login
	// First, clear the redirectUrl flash by requesting it
	// That would clear things up if there was unsuccessful
	// attemps to login
	session.Flashes("redirectUrl")
	redirectUrl := req.URL.Query().Get("redirectUrl")
	if len(redirectUrl) != 0 {
		session.AddFlash(redirectUrl, "redirectUrl")
	}

	// Save nonce and state into session store
	session.Values["state"] = state
	session.Values["nonce"] = nonce
	if err := session.Save(req, w); err != nil {
		logrus.WithError(err).Error("oidc: could not save session")
		h.err(apierrors.NewInternalError(errors.New("unable to save session")), w, req)
		return
	}

	// redirect to oidc issuer for login
	cfg := h.getConfig(req)
	http.Redirect(w, req, cfg.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusTemporaryRedirect)
}

func (h *OIDCHandler) ServeLogout(w http.ResponseWriter, req *http.Request) {

	// retrieve the session
	session, err := h.store.Get(req, "auth-session")
	if err != nil {
		logrus.WithError(err).Error("oidc: unable to get session")
		h.err(apierrors.NewUnauthorized("cannot find session"), w, req)
		return
	}

	// clear the session
	session.Values = map[interface{}]interface{}{}
	if err := session.Save(req, w); err != nil {
		logrus.WithError(err).Error("oidc: could not save session")
		h.err(apierrors.NewInternalError(errors.New("unable to save session")), w, req)
		return
	}

	// build the returnTo URL
	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}
	returnTo, err := url.Parse(fmt.Sprintf("%s://%s", scheme, req.Host))
	if err != nil {
		logrus.WithError(err).Error("oidc: could not parse returnTo url")
		h.err(apierrors.NewInternalError(errors.New("unable to parse returnTo url")), w, req)
		return
	}

	// build the oidc-compliant url for redirecting to logout endpoint
	logoutUrl, err := url.Parse(fmt.Sprintf("%s/protocol/openid-connect/logout", h.issuerUrl))
	if err != nil {
		logrus.WithError(err).Error("oidc: could not parse logout url")
		h.err(apierrors.NewInternalError(errors.New("unable to parse logout url")), w, req)
		return
	}
	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", h.clientId)
	logoutUrl.RawQuery = parameters.Encode()

	// redirect user to oidc logout endpoint
	http.Redirect(w, req, logoutUrl.String(), http.StatusTemporaryRedirect)
}

func (h *OIDCHandler) err(statusError *apierrors.StatusError, w http.ResponseWriter, req *http.Request) {
	responsewriters.ErrorNegotiated(apierrors.FromObject(&statusError.ErrStatus), h.codec, metav1.SchemeGroupVersion, w, req)
}

func (h *OIDCHandler) ServeCallback(w http.ResponseWriter, req *http.Request) {

	// retrieve state value from session
	session, err := h.store.Get(req, "auth-session")
	if err != nil {
		logrus.WithError(err).Error("oidc: unable to get session")
		h.err(apierrors.NewUnauthorized("cannot find session"), w, req)
		return
	}

	// check that the url state and the cookie state match
	if req.URL.Query().Get("state") != session.Values["state"] {
		logrus.Warnf("oidc: state mismatch")
		h.err(apierrors.NewUnauthorized("state mismatch"), w, req)
		return
	}

	cfg := h.getConfig(req)
	ctx := context.Background()

	oauth2Token, err := cfg.Exchange(ctx, req.URL.Query().Get("code"))
	if err != nil {
		logrus.WithError(err).Warnf("oidc: exchange token failure")
		h.err(apierrors.NewUnauthorized("unable to exchange tokens"), w, req)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		logrus.Warnf("oidc: unable to retrieve id token: no id token in oauth2 token")
		h.err(apierrors.NewUnauthorized("no id_token field in oauth2 token"), w, req)
		return
	}

	idToken, err := h.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		logrus.WithError(err).Warn("oidc: unable to verify raw id token")
		h.err(apierrors.NewUnauthorized("failed to verify id token"), w, req)
		return
	}

	if idToken.Nonce != session.Values["nonce"] {
		logrus.WithError(err).Error("oidc: nonce mismatch")
		h.err(apierrors.NewUnauthorized("nonce mismatch"), w, req)
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		logrus.WithError(err).Error("oidc: could not decode id_token profile")
		h.err(apierrors.NewInternalError(errors.New("unable to retrieve id token profile")), w, req)
		return
	}

	session.Values["id_token"] = rawIDToken
	session.Values["access_token"] = oauth2Token.AccessToken
	session.Values["profile"] = profile

	encodedProfile, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		logrus.WithError(err).Error("oidc: could not encode profile")
		h.err(apierrors.NewInternalError(errors.New("unable to encode profile")), w, req)
		return
	}

	logrus.Infof(string(encodedProfile))
	logrus.Infof(rawIDToken)
	logrus.Infof(oauth2Token.AccessToken)

	delete(session.Values, "state")
	delete(session.Values, "nonce")

	if err := session.Save(req, w); err != nil {
		logrus.WithError(err).Error("oidc: could not save session")
		h.err(apierrors.NewInternalError(errors.New("unable to save session")), w, req)
		return
	}

	var redirectUri = "/apis"
	redirectUris := session.Flashes("redirectUrl")
	if len(redirectUris) > 0 {
		redirectUri, ok = redirectUris[0].(string)
		if !ok {
			logrus.Warnf("invalid redirectUri type: %#v", redirectUris[0])
		}
	}

	http.Redirect(w, req, redirectUri, http.StatusSeeOther)
}

func randString(nByte int) ([]byte, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, err
	}
	return b, nil
}
