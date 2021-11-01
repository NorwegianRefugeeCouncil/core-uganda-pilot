package organizations

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/server/admin/templates"
	"github.com/nrc-no/core/pkg/store"
	"net/http"
)

func restfulUpdate(orgStore store.OrganizationStore) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handleUpdate(req.PathParameter("organizationId"), orgStore)(res.ResponseWriter, req.Request)
	}
}

func handleUpdate(id string, orgStore store.OrganizationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		org, err := parseOrganization(req)
		if err != nil {
			_ = renderOrganization(w, org, nil, false, err)
			return
		}
		org.ID = id

		org, err = orgStore.Update(req.Context(), org)
		if err != nil {
			_ = renderOrganization(w, org, nil, false, err)
			return
		}

		w.Header().Set("Location", "/admin/organizations/"+org.ID)
		w.WriteHeader(http.StatusSeeOther)

		return
	}
}

func renderOrganization(w http.ResponseWriter, org *types.Organization, identityProviders []*types.IdentityProvider, isNew bool, err error) error {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	orgName := ""
	orgID := ""
	if org != nil {
		orgID = org.ID
		orgName = org.Name
	}

	return templates.Template.ExecuteTemplate(w, "organization_add", map[string]interface{}{
		"ID":                orgID,
		"IsNew":             isNew,
		"Error":             errMsg,
		"Name":              orgName,
		"IdentityProviders": identityProviders,
	})
}

func parseOrganization(req *http.Request) (*types.Organization, error) {
	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	q := req.Form
	orgName := q.Get("organization_name")
	org := &types.Organization{
		Name: orgName,
	}
	return org, nil
}
