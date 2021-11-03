package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/utils/pointers"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

func handleAcceptLoginRequest(
	w http.ResponseWriter,
	req *http.Request,
	userSession *sessions.Session,
	enqueue func(fn func()),
	hydraAdmin admin.ClientService,
) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {

	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {

		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		acceptResp, err := hydraAdmin.AcceptLoginRequest(&admin.AcceptLoginRequestParams{
			Body: &models.AcceptLoginRequest{
				Acr:         "",
				Context:     nil,
				Remember:    false,
				RememberFor: 0,
				Subject:     pointers.String(authRequest.Identity.ID),
			},
			LoginChallenge: authRequest.LoginChallenge,
			Context:        req.Context(),
		})
		if err != nil {
			logrus.Warnf("could not accept login request: %v", err)
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		http.Redirect(w, req, *acceptResp.Payload.RedirectTo, http.StatusTemporaryRedirect)

	}
}
