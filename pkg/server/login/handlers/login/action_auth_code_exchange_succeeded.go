package login

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	store2 "github.com/nrc-no/core/pkg/server/login/store"
	"github.com/nrc-no/core/pkg/store"
	"github.com/sirupsen/logrus"
	"net/http"
)

func handleAuthCodeExchangeSucceeded(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func()), idpStore store.IdentityProviderStore, ctx context.Context, loginStore store2.Interface) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {

		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		// getting identity provider
		idp, err := idpStore.Get(ctx, authRequest.IdentityProviderId, store.IdentityProviderGetOptions{ReturnClientSecret: false})
		if err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		identifier, err := loginStore.FindOidcIdentifier(authRequest.Claims.Subject, idp.Domain)
		if err != nil {
			if meta.ReasonForError(err) == meta.StatusReasonNotFound {
				newIdentity, err := loginStore.CreateOidcIdentity(
					idp.Domain,
					authRequest.Claims.Subject,
					authRequest.AccessToken,
					authRequest.RefreshToken,
					authRequest.IDToken)
				if err != nil {
					logrus.Warnf("failed to create new oidc identity: %v", err)
					enqueue(func() {
						_ = authRequest.Fail(err)
					})
					return
				}
				identifier = newIdentity.Credentials[0].Identifiers[0]
			} else {
				logrus.Warnf("could not retrieve oidc identifier for user: %v", err)
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
				return
			}
		}

		authRequest.Identity = identifier.Credential.Identity

		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		enqueue(func() {
			if err := authRequest.AcceptLoginRequest(); err != nil {
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
				return
			}
		})

	}
}
