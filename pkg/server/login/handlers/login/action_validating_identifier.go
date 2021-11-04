package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"go.uber.org/zap"
	"net/http"
)

func handleValidatingIdentifier(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func())) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("event", authrequest.EventUseIdentityProvider))
		l.Debug("entered state")

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		l.Debug("parsing form data")
		if err := req.ParseForm(); err != nil {
			l.Error("failed to parse form data", zap.Error(err))
			enqueue(func() {
				l.Debug("failing identifier validation")
				if err := authRequest.FailIdentifierValidation(err); err != nil {
					l.Error("failed to fail identifier validation", zap.Error(err))
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
				}
			})
			return
		}
		q := req.Form

		l.Debug("getting email domain")
		email := q.Get("email")
		emailDomain, err := getEmailDomain(email)
		if err != nil {
			l.Error("failed to get email domain", zap.Error(err))
			enqueue(func() {
				if err := authRequest.FailIdentifierValidation(err); err != nil {
					l.Error("failed to fail identifier validation", zap.Error(err))
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
				}
			})
			return
		}

		authRequest.Identifier = email
		authRequest.EmailDomain = emailDomain

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		enqueue(func() {
			l.Debug("succeeding identifier validation")
			if err := authRequest.SucceedIdentifierValidation(); err != nil {
				enqueue(func() {
					l.Error("failed to succeed identifier validation", zap.Error(err))
					if err := authRequest.FailIdentifierValidation(err); err != nil {
						l.Error("failed to fail identifier validation", zap.Error(err))
						enqueue(func() {
							_ = authRequest.Fail(err)
						})
					}
				})
			}
		})

	}
}
