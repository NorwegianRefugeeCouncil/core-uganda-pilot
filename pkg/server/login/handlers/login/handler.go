package login

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/sessions"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	loginstore "github.com/nrc-no/core/pkg/server/login/store"
	"github.com/nrc-no/core/pkg/store"
	"github.com/ory/hydra-client-go/client/admin"
)

type Handler struct {
	loginWs *restful.WebService
}

func NewHandler(
	sessionStore sessions.Store,
	orgStore store.OrganizationStore,
	idpStore store.IdentityProviderStore,
	loginStore loginstore.Interface,
	selfURL string,
	hydraAdmin admin.ClientService,
) (*Handler, error) {
	h := &Handler{}

	requestActionHandler := handleAuthRequestAction(
		sessionStore,
		idpStore,
		orgStore,
		loginStore,
		hydraAdmin,
		selfURL,
	)

	loginWs := new(restful.WebService).Path("/login")
	h.loginWs = loginWs

	loginWs.Route(loginWs.GET("/identify").
		Operation("login").
		To(func(req *restful.Request, res *restful.Response) {
			requestActionHandler(authrequest.EventRequestLogin, req.PathParameters(), req.Request.URL.Query())(res.ResponseWriter, req.Request)
		}))

	loginWs.Route(loginWs.POST("/identify").
		Operation("provide_credentials").
		To(func(req *restful.Request, res *restful.Response) {
			requestActionHandler(authrequest.EventProvideIdentifier, req.PathParameters(), req.Request.URL.Query())(res.ResponseWriter, req.Request)
		}))

	loginWs.Route(loginWs.POST("/oidc/{identityProviderId}").
		Operation("use_identity_provider").
		To(func(req *restful.Request, res *restful.Response) {
			requestActionHandler(authrequest.EventUseIdentityProvider, req.PathParameters(), req.Request.URL.Query())(res.ResponseWriter, req.Request)
		}))

	loginWs.Route(loginWs.GET("/callback").
		Operation("call_oidc_callback").
		To(func(req *restful.Request, res *restful.Response) {
			requestActionHandler(authrequest.EventCallOidcCallback, req.PathParameters(), req.Request.URL.Query())(res.ResponseWriter, req.Request)
		}))

	loginWs.Route(loginWs.POST("/consent/approve").
		Operation("approve_consent_request").
		To(func(req *restful.Request, res *restful.Response) {
			requestActionHandler(authrequest.EventApproveConsentChallenge, req.PathParameters(), req.Request.URL.Query())(res.ResponseWriter, req.Request)
		}))

	loginWs.Route(loginWs.POST("/consent/decline").
		Operation("decline_consent_request").
		To(func(req *restful.Request, res *restful.Response) {
			requestActionHandler(authrequest.EventDeclineConsentChallenge, req.PathParameters(), req.Request.URL.Query())(res.ResponseWriter, req.Request)
		}))

	loginWs.Route(loginWs.GET("/consent").
		Operation("receive_consent_request").
		To(func(req *restful.Request, res *restful.Response) {
			requestActionHandler(authrequest.EventReceiveConsentChallenge, req.PathParameters(), req.Request.URL.Query())(res.ResponseWriter, req.Request)
		}))

	return h, nil
}

func (h *Handler) WebServices() []*restful.WebService {
	return []*restful.WebService{
		h.loginWs,
	}
}
