package organizations

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/store"
	"net/http"
)

func restfulAddPost(orgStore store.OrganizationStore) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handleAddPost(orgStore)(res.ResponseWriter, req.Request)
	}
}

func handleAddPost(orgStore store.OrganizationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		org, err := parseOrganization(req)
		if err != nil {
			_ = renderOrganization(w, org, nil, true, err)
			return
		}

		org, err = orgStore.Create(req.Context(), org)
		if err != nil {
			_ = renderOrganization(w, org, nil, true, err)
			return
		}

		w.Header().Set("Location", "/admin/organizations/"+org.ID)
		w.WriteHeader(http.StatusSeeOther)

		return
	}
}
