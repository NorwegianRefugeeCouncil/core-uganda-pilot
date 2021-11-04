package login

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"go.uber.org/zap"
	"net/http"
)

func handleConsentRequestApproved(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func()), getConsentRequest func(consentChallenge string) (*models.ConsentRequest, error), hydraAdmin admin.ClientService, ctx context.Context) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateConsentRequestApproved))
		l.Debug("entered state")

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		l.Debug("getting consent request")
		consentRequest, err := getConsentRequest(authRequest.ConsentChallenge)
		if err != nil {
			l.Error("failed to get consent request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		l.Debug("accepting consent request")
		resp, err := hydraAdmin.AcceptConsentRequest(&admin.AcceptConsentRequestParams{
			Context:          ctx,
			ConsentChallenge: authRequest.ConsentChallenge,
			Body: &models.AcceptConsentRequest{

				GrantAccessTokenAudience: consentRequest.RequestedAccessTokenAudience,
				GrantScope:               consentRequest.RequestedScope,
				Remember:                 true,
				RememberFor:              0,
			},
		})
		if err != nil {
			l.Error("failed to accept consent request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		authRequest.PostConsentURL = *resp.Payload.RedirectTo

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		enqueue(func() {
			l.Debug("accepting auth request")
			if err := authRequest.Accept(); err != nil {
				l.Error("failed to accept auth request", zap.Error(err))
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
			}
		})

	}
}
