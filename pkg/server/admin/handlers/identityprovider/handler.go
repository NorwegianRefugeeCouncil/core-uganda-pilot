package identityprovider

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/mimetypes"
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

	ws := new(restful.WebService).
		Path("/apis/admin.nrc.no/v1/identityproviders").
		Doc("identityproviders.admin.nrc.no API")

	h.webService = ws

	ws.Route(ws.GET("/").To(h.RestfulList).
		Doc("lists all identity providers").
		Operation("listIdentityProviders").
		Param(restful.QueryParameter(constants.ParamOrganizationID, "organization id").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Produces(mimetypes.ApplicationJson).
		Writes(types.IdentityProviderList{}).
		Returns(http.StatusOK, "OK", types.IdentityProviderList{}))

	ws.Route(ws.GET(fmt.Sprintf("/{%s}", constants.ParamIdentityProviderID)).To(h.RestfulGet).
		Doc("gets an identity provider").
		Operation("getIdentityProvider").
		Param(restful.PathParameter(constants.ParamIdentityProviderID, "identity provider id").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Produces(mimetypes.ApplicationJson).
		Writes(types.IdentityProvider{}).
		Returns(http.StatusOK, "OK", types.IdentityProvider{}))

	ws.Route(ws.POST("/").To(h.RestfulCreate).
		Doc("create an identity provider").
		Operation("createIdentityProvider").
		Consumes(mimetypes.ApplicationJson).
		Produces(mimetypes.ApplicationJson).
		Reads(types.IdentityProvider{}).
		Writes(types.IdentityProvider{}).
		Returns(http.StatusOK, "OK", types.IdentityProvider{}))

	ws.Route(ws.PUT(fmt.Sprintf("/{%s}", constants.ParamIdentityProviderID)).To(h.RestfulUpdate).
		Doc("update an identity provider").
		Operation("updateIdentityProvider").
		Param(restful.PathParameter(constants.ParamIdentityProviderID, "identity provider id").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Consumes(mimetypes.ApplicationJson).
		Produces(mimetypes.ApplicationJson).
		Reads(types.IdentityProvider{}).
		Writes(types.IdentityProvider{}).
		Returns(http.StatusOK, "OK", types.IdentityProvider{}))

	return h
}
