package login

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	store2 "github.com/nrc-no/core/pkg/server/login/store"
	"github.com/nrc-no/core/pkg/store"
	"go.uber.org/zap"
	"net/http"
)

func handleAuthCodeExchangeSucceeded(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func()), idpStore store.IdentityProviderStore, ctx context.Context, loginStore store2.Interface) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateAuthCodeExchangeSucceeded))
		l.Debug("entered state")

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		l.Debug("getting identity provider")
		idp, err := idpStore.Get(ctx, authRequest.IdentityProviderId, store.IdentityProviderGetOptions{ReturnClientSecret: false})
		if err != nil {
			l.Error("failed to get identity provider", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
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
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
					return
				}
				identifier = newIdentity.Credentials[0].Identifiers[0]
			} else {
				l.Error("failed to get user identifier for oidc provider", zap.Error(err))
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
				return
			}
		}

		authRequest.Identity = identifier.Credential.Identity

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		enqueue(func() {
			l.Debug("accepting login request")
			if err := authRequest.AcceptLoginRequest(); err != nil {
				l.Error("failed to accept login request")
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
				return
			}
		})

	}
}
