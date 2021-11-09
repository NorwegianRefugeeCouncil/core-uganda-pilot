package organization

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
	store      store.OrganizationStore
	webService *restful.WebService
}

func (h *Handler) WebService() *restful.WebService {
	return h.webService
}

func NewHandler(store store.OrganizationStore) *Handler {
	h := &Handler{store: store}

	ws := new(restful.WebService).
		Path("/apis/admin.nrc.no/v1/organizations").
		Doc("organizations.admin.nrc.no api")

	h.webService = ws

	orgPath := fmt.Sprintf("/{%s}", constants.ParamOrganizationID)

	ws.Route(ws.GET("/").To(h.RestfulList).
		Doc("lists all organizations").
		Operation("listOrganizations").
		Param(restful.PathParameter(constants.ParamOrganizationID, "organization id").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Produces(mimetypes.ApplicationJson).
		Writes(types.OrganizationList{}).
		Returns(http.StatusOK, "OK", types.OrganizationList{}))

	ws.Route(ws.POST("/").To(h.RestfulCreate).
		Doc("create an organization").
		Operation("createOrganization").
		Consumes(mimetypes.ApplicationJson).
		Produces(mimetypes.ApplicationJson).
		Reads(types.Organization{}).
		Writes(types.Organization{}).
		Returns(http.StatusOK, "OK", types.Organization{}))

	ws.Route(ws.GET(orgPath).To(h.RestfulGet).
		Doc("get an organization").
		Operation("getOrganization").
		Param(restful.PathParameter(constants.ParamOrganizationID, "organization id").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Produces(mimetypes.ApplicationJson).
		Reads(types.Organization{}).
		Writes(types.Organization{}).
		Returns(http.StatusOK, "OK", types.Organization{}))

	return h
}
