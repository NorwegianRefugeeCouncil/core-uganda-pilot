package organizations

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/server/admin/templates"
	"github.com/nrc-no/core/pkg/store"
	"net/http"
)

func restfulList(orgStore store.OrganizationStore) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handleList(orgStore)(res.ResponseWriter, req.Request)
	}
}

func handleList(orgStore store.OrganizationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		orgs, err := orgStore.List(req.Context())
		if err != nil {
			_ = templates.Template.ExecuteTemplate(w, "organization_list", map[string]interface{}{
				"Error": err.Error(),
			})
			return
		}
		_ = templates.Template.ExecuteTemplate(w, "organization_list", map[string]interface{}{
			"Organizations": orgs,
		})
		return
	}
}
