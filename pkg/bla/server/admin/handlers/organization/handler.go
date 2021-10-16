package organization

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/bla/constants"
	"github.com/nrc-no/core/pkg/bla/store"
	"github.com/nrc-no/core/pkg/bla/types"
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

	ws := new(restful.WebService).Path("/organizations").
		Consumes("application/json").
		Produces("application/json")
	h.webService = ws

	OrganizationPath := fmt.Sprintf("/{%s}", constants.ParamOrganizationID)

	ws.Route(ws.GET("/").To(h.RestfulList).
		Doc("lists all organizations").
		Operation("listOrganizations").
		Param(restful.PathParameter(constants.ParamOrganizationID, "organization id").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Writes(types.OrganizationList{}).
		Returns(http.StatusOK, "OK", types.OrganizationList{}))

	ws.Route(ws.POST("/").To(h.RestfulCreate).
		Doc("create an organization").
		Operation("createOrganization").
		Reads(types.Organization{}).
		Writes(types.Organization{}).
		Returns(http.StatusOK, "OK", types.Organization{}))

	ws.Route(ws.POST(OrganizationPath).To(h.RestfulGet).
		Doc("get an organization").
		Operation("getOrganization").
		Param(restful.PathParameter(constants.ParamOrganizationID, "organization id").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Reads(types.Organization{}).
		Writes(types.Organization{}).
		Returns(http.StatusOK, "OK", types.Organization{}))

	return h
}
