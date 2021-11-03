package login

import (
	"context"
	"errors"
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/server/login/templates"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/utils/sets"
	"net/http"
)

func handlePromptingForIdentityProvider(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func()), idpStore store.IdentityProviderStore, ctx context.Context, orgStore store.OrganizationStore) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		if err := authRequest.Save(w, req, userSession); err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		// Retrieving suitable identity provider for given identifier
		identityProviders, err := idpStore.FindForEmailDomain(ctx, authRequest.EmailDomain, store.IdentityProviderListOptions{
			ReturnClientSecret: false,
		})
		if err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}
		if len(identityProviders) == 0 {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		// ensuring a single organization match the given email domain
		organizationIDs := sets.NewString()
		for _, idp := range identityProviders {
			organizationIDs.Insert(idp.OrganizationID)
		}
		if len(organizationIDs) > 1 {
			enqueue(func() {
				_ = authRequest.Fail(errors.New("email address domain conflict"))
			})
			return
		}

		// get organization
		organization, err := orgStore.Get(ctx, organizationIDs.List()[0])
		if err != nil {
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		// prompt choosing identity provider
		err = templates.Template.ExecuteTemplate(w, "login_idp", map[string]interface{}{
			"OrganizationName":  organization.Name,
			"IdentityProviders": identityProviders,
		})

	}
}
