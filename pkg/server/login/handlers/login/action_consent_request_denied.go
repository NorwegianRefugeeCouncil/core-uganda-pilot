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

func handleConsentRequestDenied(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func()), hydraAdmin admin.ClientService, ctx context.Context) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateConsentRequestDeclined))
		l.Debug("entered state")

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		l.Debug("rejecting consent request")
		resp, err := hydraAdmin.RejectConsentRequest(&admin.RejectConsentRequestParams{
			Context:          ctx,
			ConsentChallenge: authRequest.ConsentChallenge,
			Body:             &models.RejectRequest{},
		})
		if err != nil {
			l.Error("failed to reject consent request", zap.Error(err))
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

		l.Debug("declining auth request")
		enqueue(func() {
			if err := authRequest.Decline(); err != nil {
				l.Error("failed to decline auth request", zap.Error(err))
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
			}
		})

	}
}
