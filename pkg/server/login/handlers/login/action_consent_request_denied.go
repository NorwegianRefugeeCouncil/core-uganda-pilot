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

func handleConsentRequestDenied(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func()), hydraAdmin admin.ClientService, ctx context.Context) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		resp, err := hydraAdmin.RejectConsentRequest(&admin.RejectConsentRequestParams{
			Context:          ctx,
			ConsentChallenge: authRequest.ConsentChallenge,
			Body:             &models.RejectRequest{},
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
			if err := authRequest.Decline(); err != nil {
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
			}
		})

	}
}
