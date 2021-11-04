package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/ory/hydra-client-go/models"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

func handleReceivedConsentChallenge(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func()), requestParameters url.Values, getConsentRequest func(consentChallenge string) (*models.ConsentRequest, error)) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateReceivedConsentChallenge))
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
		consentChallenge := requestParameters.Get("consent_challenge")
		consentRequest, err := getConsentRequest(consentChallenge)
		if err != nil {
			l.Error("failed to get consent request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}
		authRequest.ConsentChallenge = consentChallenge

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		if consentRequest.Skip {
			l.Debug("skipping consent request")
			enqueue(func() {
				if err := authRequest.SkipConsentRequest(); err != nil {
					l.Error("failed to skip consent request", zap.Error(err))
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
					return
				}
			})
		} else {
			l.Debug("presenting consent challenge")
			enqueue(func() {
				if err := authRequest.PresentConsentChallenge(); err != nil {
					l.Error("failed to present consent challenge", zap.Error(err))
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
					return
				}
			})
		}
	}
}
