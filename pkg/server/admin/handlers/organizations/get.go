package organizations

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/store"
	"net/http"
)

func restfulGet(orgStore store.OrganizationStore, idpStore store.IdentityProviderStore) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handleGet(req.PathParameter("organizationId"), orgStore, idpStore)(res.ResponseWriter, req.Request)
	}
}

func handleGet(id string, orgStore store.OrganizationStore, idpStore store.IdentityProviderStore) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		org, err := orgStore.Get(req.Context(), id)
		idps, err := idpStore.List(req.Context(), id, store.IdentityProviderListOptions{
			ReturnClientSecret: false,
		})
		_ = renderOrganization(w, org, idps, false, err)
	}
}
