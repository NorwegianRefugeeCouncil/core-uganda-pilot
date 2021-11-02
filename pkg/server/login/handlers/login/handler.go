package login

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/sessions"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	loginstore "github.com/nrc-no/core/pkg/server/login/store"
	"github.com/nrc-no/core/pkg/store"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"net/url"
)

type Handler struct {
	hydraAdmin admin.ClientService
	ws         *restful.WebService
}

func NewHandler(
	sessionStore sessions.Store,
	orgStore store.OrganizationStore,
	idpStore store.IdentityProviderStore,
	loginStore loginstore.Interface,
	selfURL string,
) (*Handler, error) {
	h := &Handler{}
	adminURL, err := url.Parse("http://localhost:4445")
	if err != nil {
		return nil, err
	}
	hydraAdmin := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{adminURL.Scheme}, Host: adminURL.Host, BasePath: adminURL.Path}).Admin
	h.hydraAdmin = hydraAdmin

	ws := new(restful.WebService)

	requestActionHandler := handleAuthRequestAction(
		sessionStore,
		idpStore,
		orgStore,
		loginStore,
		hydraAdmin,
		selfURL,
	)

	ws.Route(ws.GET("/login").To(func(req *restful.Request, res *restful.Response) {
		requestActionHandler(authrequest.EventRequestLogin, req.PathParameters(), req.Request.URL.Query())(res.ResponseWriter, req.Request)
	}))

	ws.Route(ws.POST("/login").To(func(req *restful.Request, res *restful.Response) {
		requestActionHandler(authrequest.EventProvideIdentifier, req.PathParameters(), req.Request.URL.Query())(res.ResponseWriter, req.Request)
	}))

	ws.Route(ws.POST("/login/oidc/{identityProviderId}").To(func(req *restful.Request, res *restful.Response) {
		requestActionHandler(authrequest.EventUseIdentityProvider, req.PathParameters(), req.Request.URL.Query())(res.ResponseWriter, req.Request)
	}))

	ws.Route(ws.GET("/oidc/callback").To(func(req *restful.Request, res *restful.Response) {
		requestActionHandler(authrequest.EventCallOidcCallback, req.PathParameters(), req.Request.URL.Query())(res.ResponseWriter, req.Request)
	}))

	ws.Route(ws.POST("/consent/approve").To(func(req *restful.Request, res *restful.Response) {
		requestActionHandler(authrequest.EventApproveConsentChallenge, req.PathParameters(), req.Request.URL.Query())(res.ResponseWriter, req.Request)
	}))

	ws.Route(ws.POST("/consent/decline").To(func(req *restful.Request, res *restful.Response) {
		requestActionHandler(authrequest.EventDeclineConsentChallenge, req.PathParameters(), req.Request.URL.Query())(res.ResponseWriter, req.Request)
	}))

	ws.Route(ws.GET("/consent").To(func(req *restful.Request, res *restful.Response) {
		requestActionHandler(authrequest.EventReceiveConsentChallenge, req.PathParameters(), req.Request.URL.Query())(res.ResponseWriter, req.Request)
	}))

	h.ws = ws
	return h, nil
}

func (h *Handler) WebService() *restful.WebService {
	return h.ws
}
