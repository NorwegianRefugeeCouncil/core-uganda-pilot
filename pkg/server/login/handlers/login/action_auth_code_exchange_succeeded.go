package login

import (
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	store2 "github.com/nrc-no/core/pkg/server/login/store"
	"github.com/nrc-no/core/pkg/store"
	"go.uber.org/zap"
	"net/http"
)

func handleAuthCodeExchangeSucceeded(
	req *http.Request,
	dispatch func(evt string),
	idpStore store.IdentityProviderStore,
	loginStore store2.Interface,
) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateAuthCodeExchangeSucceeded))

		l.Debug("getting identity provider")
		idp, err := idpStore.Get(ctx, authRequest.IdentityProviderId, store.IdentityProviderGetOptions{ReturnClientSecret: false})
		if err != nil {
			l.Error("failed to get identity provider", zap.Error(err))
			return err
		}

		l.Debug("finding user identifier for oidc provider", zap.String("subject", authRequest.Claims.Subject), zap.String("domain", idp.Domain))
		identifier, err := loginStore.FindOidcIdentifier(authRequest.Claims.Subject, idp.Domain)
		if err != nil {
			if meta.ReasonForError(err) == meta.StatusReasonNotFound {
				l.Info("user identifier not found. creating new identity")
				newIdentity, err := loginStore.CreateOidcIdentity(
					idp.Domain,
					authRequest.Claims.Subject,
					authRequest.AccessToken,
					authRequest.RefreshToken,
					authRequest.IDToken)
				if err != nil {
					l.Error("failed to create new identity", zap.Error(err))
					return err
				}
				identifier = newIdentity.Credentials[0].Identifiers[0]
			} else {
				l.Error("failed to get user identifier for oidc provider", zap.Error(err))
				return err
			}
		}

		authRequest.Identity = identifier.Credential.Identity

		dispatch(authrequest.EventAcceptLoginRequest)
		return nil

	}
}
