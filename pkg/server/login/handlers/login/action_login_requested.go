package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/ory/hydra-client-go/models"
	"go.uber.org/zap"
	"net/http"
)

func handleLoginRequested(
	w http.ResponseWriter,
	req *http.Request,
	userSession *sessions.Session,
	enqueue func(fn func()),
	getLoginRequest func(loginChallenge string) (*models.LoginRequest, error),
) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {

		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateLoginRequested))
		l.Debug("entered state")

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
		}

		l.Debug("getting login challenge")
		loginChallenge := req.URL.Query().Get("login_challenge")
		loginRequest, err := getLoginRequest(loginChallenge)
		if err != nil {
			l.Error("failed to get login challenge", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		authRequest.LoginChallenge = loginChallenge

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		l.Debug("checking if should skip login request or not")
		if authRequest.Identity != nil && loginRequest.Skip != nil && *loginRequest.Skip {

			l.Debug("skipping login request")
			enqueue(func() {
				err := authRequest.SkipLoginRequest()
				if err != nil {
					l.Debug("failed to skip login request", zap.Error(err))
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
				}
			})
			return
		} else {
			l.Debug("performing login request")
			enqueue(func() {
				err := authRequest.PerformLogin()
				if err != nil {
					l.Error("failed to perform login request", zap.Error(err))
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
				}
			})
			return
		}
	}
}
