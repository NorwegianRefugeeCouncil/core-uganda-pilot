package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/ory/hydra-client-go/models"
	"github.com/sirupsen/logrus"
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

		logger := logrus.
			WithField("server", "core-login").
			WithField("handler", "login").
			WithField("state", authrequest.StateLoginRequested)

		logger.Trace("handling login requested")

		logger.Trace("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			logger.WithError(err).Errorf("failed to save auth request")
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
		}

		logger.Tracef("getting login challenge")
		loginChallenge := req.URL.Query().Get("login_challenge")
		loginRequest, err := getLoginRequest(loginChallenge)
		if err != nil {
			logger.WithError(err).Tracef("failed to get login challenge")
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		authRequest.LoginChallenge = loginChallenge

		logger.Trace("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			logger.WithError(err).Errorf("failed to save auth request")
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		logger.Tracef("checking if should skip login request or not")
		if authRequest.Identity != nil && loginRequest.Skip != nil && *loginRequest.Skip {

			logrus.Tracef("skipping login request")
			enqueue(func() {
				err := authRequest.SkipLoginRequest()
				logrus.WithError(err).Tracef("failed to invoke skipping login request")
				if err != nil {
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
				}
			})
			return
		} else {
			logrus.Tracef("performing login request")
			enqueue(func() {
				err := authRequest.PerformLogin()
				if err != nil {
					logrus.WithError(err).Tracef("failed to invoke perform login request")
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
				}
			})
			return
		}
	}
}
