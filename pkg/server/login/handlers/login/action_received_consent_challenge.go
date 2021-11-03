package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/ory/hydra-client-go/models"
	"net/http"
	"net/url"
)

func handleReceivedConsentChallenge(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func()), requestParameters url.Values, getConsentRequest func(consentChallenge string) (*models.ConsentRequest, error)) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {

		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		// getting consent request
		consentChallenge := requestParameters.Get("consent_challenge")
		consentRequest, err := getConsentRequest(consentChallenge)
		if err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}
		authRequest.ConsentChallenge = consentChallenge

		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		if consentRequest.Skip {
			enqueue(func() {
				if err := authRequest.SkipConsentRequest(); err != nil {
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
					return
				}
			})
		} else {
			enqueue(func() {
				if err := authRequest.PresentConsentChallenge(); err != nil {
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
					return
				}
			})
		}
	}
}
