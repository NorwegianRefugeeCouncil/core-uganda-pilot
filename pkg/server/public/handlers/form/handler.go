package form

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/store"
	"net/http"
)

type Handler struct {
	store      store.FormStore
	webService *restful.WebService
}

func NewHandler(store store.FormStore) *Handler {
	h := &Handler{store: store}

	ws := new(restful.WebService).Path("/forms")
	h.webService = ws

	formRoute := fmt.Sprintf("/{%s}", constants.ParamFormID)

	ws.Route(ws.POST("/").To(h.RestfulCreate).
		Doc("create a form").
		Operation("createForm").
		Consumes("application/json").
		Produces("application/json").
		Reads(types.FormDefinition{}).
		Writes(types.FormDefinition{}).
		Returns(http.StatusOK, "OK", types.FormDefinition{}))

	ws.Route(ws.GET("/").To(h.RestfulList).
		Doc("list forms").
		Operation("listForms").
		Produces("application/json").
		Writes(types.FormDefinitionList{}).
		Returns(http.StatusOK, "OK", types.FormDefinitionList{}))

	ws.Route(ws.DELETE(formRoute).To(h.RestfulDelete).
		Doc("deletes a form").
		Operation("deleteForm").
		Param(restful.PathParameter(constants.ParamFormID, "id of the form").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Returns(http.StatusNoContent, "OK", nil))

	return h
}

func (h *Handler) WebService() *restful.WebService {
	return h.webService
}
