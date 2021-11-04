package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/utils/pointers"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"go.uber.org/zap"
	"net/http"
)

func handleAwaitingConsentChallenge(w http.ResponseWriter, req *http.Request, hydraAdmin admin.ClientService, userSession *sessions.Session, enqueue func(fn func())) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateAwaitingConsentChallenge))
		l.Debug("entered state")

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		l.Debug("accepting login request")
		acceptResp, err := hydraAdmin.AcceptLoginRequest(&admin.AcceptLoginRequestParams{
			Body: &models.AcceptLoginRequest{
				Acr:         "",
				Context:     nil,
				Remember:    true,
				RememberFor: 0,
				Subject:     pointers.String(authRequest.Identity.ID),
			},
			LoginChallenge: authRequest.LoginChallenge,
			Context:        req.Context(),
		})
		if err != nil {
			l.Error("failed to accept login request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		redirectURI := *acceptResp.Payload.RedirectTo
		l.Debug("redirecting to post-login-request uri", zap.String("redirect_uri", redirectURI))
		http.Redirect(w, req, redirectURI, http.StatusTemporaryRedirect)

	}
}
