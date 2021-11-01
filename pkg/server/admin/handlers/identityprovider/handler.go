package identityprovider

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/store"
	"net/http"
)

type Handler struct {
	store      store.IdentityProviderStore
	webService *restful.WebService
}

func (h *Handler) WebService() *restful.WebService {
	return h.webService
}

func NewHandler(store store.IdentityProviderStore) *Handler {
	h := &Handler{store: store}

	ws := new(restful.WebService).Path("/identityproviders").
		Consumes("application/json").
		Produces("application/json")
	h.webService = ws

	identityProviderPath := fmt.Sprintf("/{%s}", constants.ParamOrganizationID)

	ws.Route(ws.GET("/").To(h.RestfulList).
		Doc("lists all identity providers").
		Operation("listIdentityProviders").
		Param(restful.PathParameter(constants.ParamOrganizationID, "organization id").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Writes(types.IdentityProviderList{}).
		Returns(http.StatusOK, "OK", types.IdentityProviderList{}))

	ws.Route(ws.POST("/").To(h.RestfulCreate).
		Doc("create an identity provider").
		Operation("createIdentityProvider").
		Reads(types.IdentityProvider{}).
		Writes(types.IdentityProvider{}).
		Returns(http.StatusOK, "OK", types.IdentityProvider{}))

	ws.Route(ws.POST(identityProviderPath).To(h.RestfulUpdate).
		Doc("update an identity provider").
		Operation("updateIdentityProvider").
		Param(restful.PathParameter(constants.ParamIdentityProviderID, "identity provider id").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Reads(types.IdentityProvider{}).
		Writes(types.IdentityProvider{}).
		Returns(http.StatusOK, "OK", types.IdentityProvider{}))

	return h
}
