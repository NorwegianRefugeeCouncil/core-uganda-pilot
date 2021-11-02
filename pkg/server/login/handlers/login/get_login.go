package login

import (
	"errors"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/server/login/templates"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/utils/sets"
	"github.com/ory/hydra-client-go/client/admin"
	"net/http"
	"strings"
)

func restfulGetLogin(
	hydraAdmin admin.ClientService,
	orgStore store.OrganizationStore,
	idpStore store.IdentityProviderStore,
) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		getLogin(hydraAdmin, orgStore, idpStore)(res.ResponseWriter, req.Request)
	}
}

func getLogin(
	hydraAdmin admin.ClientService,
	orgStore store.OrganizationStore,
	idpStore store.IdentityProviderStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		if err := req.ParseForm(); err != nil {
			renderSubjectLogin(w, "", err)
			return
		}

		q := req.Form

		loginChallenge := q.Get("login_challenge")
		loginRequest, err := hydraAdmin.GetLoginRequest(&admin.GetLoginRequestParams{
			Context:        ctx,
			LoginChallenge: loginChallenge,
		})

		if err != nil {
			renderSubjectLogin(w, loginChallenge, err)
			return
		}

		emailDomain, err := getEmailDomain(q.Get("email"))
		if err != nil {
			renderSubjectLogin(w, loginChallenge, err)
			return
		}

		identityProviders, err := idpStore.FindForEmailDomain(ctx, emailDomain, store.IdentityProviderListOptions{
			ReturnClientSecret: false,
		})
		if err != nil {
			renderSubjectLogin(w, loginChallenge, err)
			return
		}
		if len(identityProviders) == 0 {
			renderSubjectLogin(w, loginChallenge, errors.New("unrecognized email address"))
			return
		}

		organizationIDs := sets.NewString()
		for _, idp := range identityProviders {
			organizationIDs.Insert(idp.OrganizationID)
		}
		if len(organizationIDs) > 1 {
			renderSubjectLogin(w, loginChallenge, errors.New("email address domain conflict"))
			return
		}

		organization, err := orgStore.Get(ctx, organizationIDs.List()[0])
		if err != nil {
			renderSubjectLogin(w, loginChallenge, err)
			return
		}

		renderIDPLogin(w, *loginRequest.Payload.Challenge, organization, identityProviders, nil)

	}
}

func renderIDPLogin(
	w http.ResponseWriter,
	challenge string,
	organization *types.Organization,
	identityProviders []*types.IdentityProvider,
	err error,
) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	_ = templates.Template.ExecuteTemplate(w, "login_idp", map[string]interface{}{
		"Error":             errMsg,
		"LoginChallenge":    challenge,
		"OrganizationName":  organization.Name,
		"IdentityProviders": identityProviders,
	})
	return
}

func renderSubjectLogin(w http.ResponseWriter, challenge string, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	templates.Template.ExecuteTemplate(w, "login_subject", map[string]interface{}{
		"Error":          errMsg,
		"LoginChallenge": challenge,
	})
	return
}

func getEmailDomain(val string) (string, error) {
	components := strings.Split(val, "@")
	if len(components) != 2 {
		return "", errors.New("invalid email address")
	}
	if len(components[0]) == 0 {
		return "", errors.New("invalid email address")
	}
	return components[1], nil
}
