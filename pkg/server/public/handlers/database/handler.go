package database

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/types"
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

	databasesPath := fmt.Sprintf("/{%s}", constants.ParamDatabaseID)

	ws.Route(ws.DELETE(databasesPath).To(h.RestfulDelete).
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
		Writes(types.DatabaseList{}).
		Returns(http.StatusOK, "OK", types.DatabaseList{}))

	ws.Route(ws.POST("/").To(h.RestfulCreate).
		Doc("create a database").
		Operation("createDatabase").
		Reads(types.Database{}).
		Writes(types.Database{}).
		Returns(http.StatusOK, "OK", types.Database{}))

	return h
}

func (h *Handler) WebService() *restful.WebService {
	return h.webService
}
