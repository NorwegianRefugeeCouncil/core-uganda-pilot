package login

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/store"
	"github.com/ory/hydra-client-go/models"
	"golang.org/x/oauth2"
	"net/http"
)

func handleUseIdentityProvider(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func()), getLoginRequest func(loginChallenge string) (*models.LoginRequest, error), pathParameters map[string]string, idpStore store.IdentityProviderStore, ctx context.Context, selfURL string) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		loginRequest, err := getLoginRequest(authRequest.LoginChallenge)
		if err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		identityProviderID := pathParameters["identityProviderId"]
		// getting identity provider
		idp, err := idpStore.Get(ctx, identityProviderID, store.IdentityProviderGetOptions{ReturnClientSecret: true})
		if err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}
		authRequest.IdentityProviderId = identityProviderID

		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		// creating state variable
		stateVar, err := createStateVariable()
		if err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}
		authRequest.StateVariable = stateVar

		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		// Getting Identity Provider Client Config
		oauth2Config, _, _, err := getOauthProvider(ctx, idp, selfURL, loginRequest)
		if err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		authCodeURL := oauth2Config.AuthCodeURL(stateVar, oauth2.SetAuthURLParam("login_hint", authRequest.Identifier))
		http.Redirect(w, req, authCodeURL, http.StatusSeeOther)

	}
}
