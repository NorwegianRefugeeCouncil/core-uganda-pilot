package login

import (
	"context"
	"errors"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/server/login/templates"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/utils/sets"
	"go.uber.org/zap"
	"net/http"
)

func handlePromptingForIdentityProvider(
	ctx context.Context,
	w http.ResponseWriter,
	idpStore store.IdentityProviderStore,
	orgStore store.OrganizationStore,
) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StatePromptingForIdentityProvider))

		l.Debug("retrieving suitable identity provider for identifier")
		identityProviders, err := idpStore.FindForEmailDomain(ctx, authRequest.EmailDomain, store.IdentityProviderListOptions{
			ReturnClientSecret: false,
		})
		if err != nil {
			l.Error("failed to retrieve suitable identity providers", zap.Error(err))
			return err
		}

		if len(identityProviders) == 0 {
			l.Error("no suitable identity providers found")
			return errors.New("no suitable identity provider")
		}

		l.Debug("ensuring a single organization match the email domain")
		organizationIDs := sets.NewString()
		for _, idp := range identityProviders {
			organizationIDs.Insert(idp.OrganizationID)
		}
		if len(organizationIDs) > 1 {
			l.Error("multiple organizations matched the email domain")
			return errors.New("multiple organizations share the same email domain")
		}

		l.Debug("getting organization", zap.String("organization_id", organizationIDs.List()[0]))
		// get organization
		organization, err := orgStore.Get(ctx, organizationIDs.List()[0])
		if err != nil {
			l.Error("failed to get organization", zap.Error(err))
			return err
		}

		l.Debug("prompting user for identity provider")
		err = templates.Template.ExecuteTemplate(w, "login_idp", map[string]interface{}{
			"OrganizationName":  organization.Name,
			"IdentityProviders": identityProviders,
		})
		if err != nil {
			l.Error("failed to prompt user for identity provider", zap.Error(err))
			return err
		}

		return nil

	}
}
