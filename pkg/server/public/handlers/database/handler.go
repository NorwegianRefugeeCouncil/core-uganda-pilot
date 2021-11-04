package database

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/store"
	"net/http"
)

type Handler struct {
	store      store.DatabaseStore
	webService *restful.WebService
}

func NewHandler(store store.DatabaseStore) *Handler {
	h := &Handler{store: store}

	ws := new(restful.WebService).Path("/databases").
		Consumes("application/json").
		Produces("application/json")
	h.webService = ws

	dbPath := fmt.Sprintf("/{%s}", constants.ParamDatabaseID)

	ws.Route(ws.DELETE(dbPath).To(h.RestfulDelete).
		Param(restful.PathParameter(constants.ParamDatabaseID, "id of the database").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Doc("deletes databases").
		Operation("deleteDatabase").
		Writes(nil).
		Returns(http.StatusNoContent, "OK", nil))

	ws.Route(ws.GET("/").To(h.RestfulList).
		Doc("lists all databases").
		Operation("listDatabases").
		Produces("application/json").
		Writes(types.DatabaseList{}).
		Returns(http.StatusOK, "OK", types.DatabaseList{}))

	ws.Route(ws.GET(dbPath).To(h.RestfulGet).
		Doc("gets a databases").
		Operation("getDatabase").
		Param(restful.PathParameter(constants.ParamDatabaseID, "id of the database").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Produces("application/json").
		Writes(types.Database{}).
		Returns(http.StatusOK, "OK", types.Database{}))

	ws.Route(ws.POST("/").To(h.RestfulCreate).
		Doc("create a database").
		Operation("createDatabase").
		Consumes("application/json").
		Produces("application/json").
		Reads(types.Database{}).
		Writes(types.Database{}).
		Returns(http.StatusOK, "OK", types.Database{}))

	return h
}

func (h *Handler) WebService() *restful.WebService {
	return h.webService
}
