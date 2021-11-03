package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"net/http"
)

func handleDeclined(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func())) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}
		http.Redirect(w, req, authRequest.PostConsentURL, http.StatusTemporaryRedirect)
	}
}
