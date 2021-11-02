package login

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/server/login/templates"
	"net/http"
)

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
	err = templates.Template.ExecuteTemplate(w, "login_idp", map[string]interface{}{
		"Error":             errMsg,
		"LoginChallenge":    challenge,
		"OrganizationName":  organization.Name,
		"IdentityProviders": identityProviders,
	})
	if err != nil {
		fmt.Println(err)
	}
	return
}

func promptUserForIdentifier(w http.ResponseWriter, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	templates.Template.ExecuteTemplate(w, "login_subject", map[string]interface{}{
		"Error": errMsg,
	})
	return
}
