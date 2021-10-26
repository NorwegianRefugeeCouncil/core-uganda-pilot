package form

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/types"
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

	return h
}

func (h *Handler) WebService() *restful.WebService {
	return h.webService
}
