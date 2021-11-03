package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"net/http"
)

func handleValidatingIdentifier(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func())) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		// Parsing form data
		if err := req.ParseForm(); err != nil {

			enqueue(func() {
				if err := authRequest.FailIdentifierValidation(err); err != nil {
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
				}
			})
			return
		}
		q := req.Form

		// Retrieving identifier from user form data
		email := q.Get("email")
		emailDomain, err := getEmailDomain(email)
		if err != nil {
			enqueue(func() {
				if err := authRequest.FailIdentifierValidation(err); err != nil {
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
				}
			})
			return
		}

		authRequest.Identifier = email
		authRequest.EmailDomain = emailDomain

		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		enqueue(func() {
			if err := authRequest.SucceedIdentifierValidation(); err != nil {
				enqueue(func() {
					if err := authRequest.FailIdentifierValidation(err); err != nil {
						enqueue(func() {
							_ = authRequest.Fail(err)
						})
					}
				})
			}
		})

	}
}
