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

func NewHandler(sessionStore sessions.Store, oauth2Config *oauth2.Config, tokenVerifier *oidc.IDTokenVerifier) *Handler {
	h := &Handler{
		sessionStore:  sessionStore,
		oauth2Config:  oauth2Config,
		tokenVerifier: tokenVerifier,
	}

	ws := new(restful.WebService).Path("/oidc")
	h.webService = ws

	ws.Route(ws.GET("/login").To(h.RestfulLogin).
		Doc("initiates login flow"))

	ws.Route(ws.GET("/callback").To(h.RestfulCallback).
		Doc("oauth2 callback"))

	return h
}

func (h *Handler) WebService() *restful.WebService {
	return h.webService
}
