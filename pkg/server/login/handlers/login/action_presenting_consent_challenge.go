package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/server/login/templates"
	"github.com/ory/hydra-client-go/models"
	"net/http"
)

func handlePresentingConsentChallenge(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func()), getConsentRequest func(consentChallenge string) (*models.ConsentRequest, error)) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		consentRequest, err := getConsentRequest(authRequest.ConsentChallenge)
		if err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		// prompt choosing identity provider
		err = templates.Template.ExecuteTemplate(w, "challenge", map[string]interface{}{
			"Scopes":     consentRequest.RequestedScope,
			"ClientName": consentRequest.Client.ClientName,
		})

	}
}
