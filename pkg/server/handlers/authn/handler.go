package authn

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type Handler struct {
	sessionStore  sessions.Store
	oauth2Config  *oauth2.Config
	tokenVerifier *oidc.IDTokenVerifier
	webService    *restful.WebService
}

func NewHandler(
	sessionKey string,
	redirectURL string,
	sessionStore sessions.Store,
	oauth2Config *oauth2.Config,
	tokenVerifier *oidc.IDTokenVerifier,
) *Handler {
	h := &Handler{
		sessionStore:  sessionStore,
		oauth2Config:  oauth2Config,
		tokenVerifier: tokenVerifier,
	}

	ws := new(restful.WebService).Path("/oidc")
	h.webService = ws

	ws.Route(ws.GET("/login").
		Doc("initiates login flow").
		To(h.RestfulLogin(sessionKey, false)))

	ws.Route(ws.GET("/renew").
		Doc("renews session").
		To(h.RestfulLogin(sessionKey, true)))

	ws.Route(ws.GET("/callback").To(h.RestfulCallback(sessionKey, redirectURL)).
		Doc("oauth2 callback").
		Param(ws.QueryParameter("redirect_uri", "redirection uri after successful authentication").Required(false)))

	ws.Route(ws.GET("/session").To(h.RestfulSession(sessionKey)).
		Doc("gets session").
		Param(ws.QueryParameter("redirect_uri", "redirection uri after successful authentication").Required(false)))

	return h
}

func (h *Handler) WebService() *restful.WebService {
	return h.webService
}
