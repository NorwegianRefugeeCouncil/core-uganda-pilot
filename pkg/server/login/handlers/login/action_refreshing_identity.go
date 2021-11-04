package login

import (
	"context"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/store"
	"github.com/ory/hydra-client-go/models"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func handleRefreshingIdentity(
	ctx context.Context,
	idpStore store.IdentityProviderStore,
	selfURL string,
	dispatch func(evt string),
	getLoginRequest func(loginChallenge string) (*models.LoginRequest, error),
) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateRefreshingIdentity))

		if len(authRequest.IdentityProviderId) > 0 {

			l.Debug("refreshing identity using identity provider", zap.String("identity_provider_id", authRequest.IdentityProviderId))

			l.Debug("getting login request")
			loginRequest, err := getLoginRequest(authRequest.LoginChallenge)
			if err != nil {
				l.Error("failed to get login request", zap.Error(err))
				return err
			}

			l.Debug("getting identity provider")
			idp, err := idpStore.Get(ctx, authRequest.IdentityProviderId, store.IdentityProviderGetOptions{ReturnClientSecret: true})
			if err != nil {
				l.Error("failed to get identity provider", zap.Error(err))
				return err
			}

			l.Debug("getting identity provider oauth config")
			oauth2Config, _, _, err := getOauthProvider(ctx, idp, selfURL, loginRequest)
			if err != nil {
				l.Error("failed to get identity provider oauth config")
				return err
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
				return err
			}

			l.Debug("saving token in session")
			authRequest.AccessToken = newToken.AccessToken
			authRequest.RefreshToken = newToken.RefreshToken
			authRequest.TokenExpiry = newToken.Expiry

			l.Debug("set refreshed identity")
			dispatch(authrequest.EventSetRefreshedIdentity)

			return nil

		} else {
			// refreshing identity using password credential
			return nil
		}
	}
}

func restoreToken(ctx context.Context, previousToken *oauth2.Token, oauth2Config *oauth2.Config) oauth2.TokenSource {
	return oauth2Config.TokenSource(ctx, previousToken)
}
