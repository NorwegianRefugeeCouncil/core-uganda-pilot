package login

import (
	"context"
	"errors"
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/server/login/templates"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/utils/sets"
	"go.uber.org/zap"
	"net/http"
)

func handlePromptingForIdentityProvider(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func()), idpStore store.IdentityProviderStore, ctx context.Context, orgStore store.OrganizationStore) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StatePromptingForIdentityProvider))
		l.Debug("entered state")

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		l.Debug("retrieving suitable identity provider for identifier")
		identityProviders, err := idpStore.FindForEmailDomain(ctx, authRequest.EmailDomain, store.IdentityProviderListOptions{
			ReturnClientSecret: false,
		})
		if err != nil {
			l.Error("failed to retrieve suitable identity providers", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}
		if len(identityProviders) == 0 {
			l.Error("no suitable identity providers found")
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		l.Debug("ensuring a single organization match the email domain")
		organizationIDs := sets.NewString()
		for _, idp := range identityProviders {
			organizationIDs.Insert(idp.OrganizationID)
		}
		if len(organizationIDs) > 1 {
			l.Error("multiple organizations matched the email domain")
			enqueue(func() {
				_ = authRequest.Fail(errors.New("email address domain conflict"))
			})
			return
		}

		l.Debug("getting organization", zap.String("organization_id", organizationIDs.List()[0]))
		// get organization
		organization, err := orgStore.Get(ctx, organizationIDs.List()[0])
		if err != nil {
			l.Error("failed to get organization", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		l.Debug("prompting user for identity provider")
		err = templates.Template.ExecuteTemplate(w, "login_idp", map[string]interface{}{
			"OrganizationName":  organization.Name,
			"IdentityProviders": identityProviders,
		})
		if err != nil {
			l.Error("failed to prompt user for identity provider", zap.Error(err))
		}

	}
}
