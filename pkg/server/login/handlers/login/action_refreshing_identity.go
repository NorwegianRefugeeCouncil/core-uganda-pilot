package login

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/store"
	"github.com/ory/hydra-client-go/models"
	"go.uber.org/zap"
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
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateRefreshingIdentity))
		l.Debug("entered state")

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		if len(authRequest.IdentityProviderId) > 0 {

			l.Debug("refreshing identity using identity provider", zap.String("identity_provider_id", authRequest.IdentityProviderId))

			l.Debug("getting login request")
			loginRequest, err := getLoginRequest(authRequest.LoginChallenge)
			if err != nil {
				l.Error("failed to get login request", zap.Error(err))
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
				return
			}

			l.Debug("getting identity provider")
			idp, err := idpStore.Get(ctx, authRequest.IdentityProviderId, store.IdentityProviderGetOptions{ReturnClientSecret: true})
			if err != nil {
				l.Error("failed to get identity provider", zap.Error(err))
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
				return
			}

			// Getting Identity Provider Client Config
			l.Debug("getting identity provider oauth config")
			oauth2Config, _, _, err := getOauthProvider(ctx, idp, selfURL, loginRequest)
			if err != nil {
				l.Error("failed to get identity provider oauth config")
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

			l.Debug("refreshing token if necessary")
			newToken, err := restoredToken.Token()
			if err != nil {
				l.Error("failed to refresh token", zap.Error(err))
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
				return
			}

			l.Debug("saving token in session")
			authRequest.AccessToken = newToken.AccessToken
			authRequest.RefreshToken = newToken.RefreshToken
			authRequest.TokenExpiry = newToken.Expiry
			if err := authRequest.Save(w, req, userSession); err != nil {
				l.Error("failed to save token in session", zap.Error(err))
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
				return
			}

			l.Debug("set refreshed identity")
			enqueue(func() {
				if err := authRequest.SetRefreshedIdentity(); err != nil {
					l.Error("failed to set refreshed identity", zap.Error(err))
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
