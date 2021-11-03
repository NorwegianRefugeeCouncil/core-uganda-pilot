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

func handleRefreshingIdentity(
	w http.ResponseWriter,
	req *http.Request,
	userSession *sessions.Session,
	idpStore store.IdentityProviderStore,
	selfURL string,
	enqueue func(fn func()),
	getLoginRequest func(loginChallenge string) (*models.LoginRequest, error),
) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {

		ctx := req.Context()

		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		if len(authRequest.IdentityProviderId) > 0 {
			// refreshing identity using openid provider

			loginRequest, err := getLoginRequest(authRequest.LoginChallenge)
			if err != nil {
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
				return
			}

			// getting identity provider
			idp, err := idpStore.Get(ctx, authRequest.IdentityProviderId, store.IdentityProviderGetOptions{ReturnClientSecret: true})
			if err != nil {
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

			previousToken := &oauth2.Token{
				AccessToken:  authRequest.AccessToken,
				RefreshToken: authRequest.RefreshToken,
				Expiry:       authRequest.TokenExpiry,
				TokenType:    authRequest.TokenType,
			}
			restoredToken := restoreToken(ctx, previousToken, oauth2Config)
			newToken, err := restoredToken.Token()
			if err != nil {
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
				return
			}

			// save token in session
			authRequest.AccessToken = newToken.AccessToken
			authRequest.RefreshToken = newToken.RefreshToken
			authRequest.TokenExpiry = newToken.Expiry
			if err := authRequest.Save(w, req, userSession); err != nil {
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
				return
			}

			// mark identity as refreshed
			enqueue(func() {
				if err := authRequest.SetRefreshedIdentity(); err != nil {
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
				}
			})
			return

		} else {

			// refreshing identity using password credential

		}

		// noop
	}
}

func restoreToken(ctx context.Context, previousToken *oauth2.Token, oauth2Config *oauth2.Config) oauth2.TokenSource {
	return oauth2Config.TokenSource(ctx, previousToken)
}
