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

	ws.Route(ws.POST("/").To(h.RestfulCreate).
		Doc("create a form").
		Operation("createForm").
		Reads(types.FormDefinition{}).
		Writes(types.FormDefinition{}).
		Returns(http.StatusOK, "OK", types.FormDefinition{}))

	ws.Route(ws.GET("/").To(h.RestfulList).
		Doc("list forms").
		Operation("listForms").
		Writes(types.FormDefinitionList{}).
		Returns(http.StatusOK, "OK", types.FormDefinitionList{}))

	ws.Route(ws.GET(fmt.Sprintf("/{%s}", constants.ParamFormID)).To(h.RestfulGet).
		Doc("get form").
		Operation("getForm").
		Param(
			ws.
				PathParameter(constants.ParamFormID, "id of the form").
				Required(true).
				DataType("string").
				DataFormat("uuid"),
			).
		Writes(types.FormDefinition{}).
		Returns(http.StatusOK, "OK", types.FormDefinition{}))

	return h
}

func (h *Handler) WebService() *restful.WebService {
	return h.webService
}
