package login

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"net/http"
)

func handleConsentRequestApproved(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func()), getConsentRequest func(consentChallenge string) (*models.ConsentRequest, error), hydraAdmin admin.ClientService, ctx context.Context) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
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
		if err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}
		resp, err := hydraAdmin.AcceptConsentRequest(&admin.AcceptConsentRequestParams{
			Context:          ctx,
			ConsentChallenge: authRequest.ConsentChallenge,
			Body: &models.AcceptConsentRequest{
				GrantAccessTokenAudience: consentRequest.RequestedAccessTokenAudience,
				GrantScope:               consentRequest.RequestedScope,
			},
		})
		if err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		authRequest.PostConsentURL = *resp.Payload.RedirectTo

		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		enqueue(func() {
			if err := authRequest.Accept(); err != nil {
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
			}
		})

	}
}
