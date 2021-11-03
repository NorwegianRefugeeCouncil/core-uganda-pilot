package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/ory/hydra-client-go/models"
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

		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
		}

		// getting login request
		loginChallenge := req.URL.Query().Get("login_challenge")
		loginRequest, err := getLoginRequest(loginChallenge)
		if err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		authRequest.LoginChallenge = loginChallenge

		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		if loginRequest.Skip != nil && *loginRequest.Skip {
			enqueue(func() {
				err := authRequest.SkipLoginRequest()
				if err != nil {
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
				}
			})
			return
		} else {
			enqueue(func() {
				err := authRequest.PerformLogin()
				if err != nil {
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
				}
			})
			return
		}
	}
}
