package handler

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/server/data/api"
)

const (
	stringDataType          = "string"
	intDataType             = "int"
	booleanDataType         = "boolean"
	pathParamTable          = "table"
	pathParamId             = "id"
	queryParamRev           = "rev"
	queryParamSince         = "since"
	queryParamIsReplication = "replication"
)

// Handler is the HTTP Handler for the database API
// It is used to serve the following endpoints:
//   CreateTable: PUT    /api/v1/<table>
//   GetTable:    GET    /api/v1/tables/{table}
//   GetTables:   GET    /api/v1/tables
//   PutRow:      PUT    /api/v1/tables/{table}/rows/{row}
//   GetRecord:      GET    /api/v1/tables/{table}/rows/{row}?revision={revision}
//   GetRows:     GET    /api/v1/tables/{table}/rows
//   GetChanges:  GET    /api/v1/changes?since={seq}
type Handler struct {
	engine api.Engine
	ws     *restful.WebService
}

func (h *Handler) WebService() *restful.WebService {
	return h.ws
}

func NewHandler(engine api.Engine) *Handler {

	ws := new(restful.WebService).
		Path("/apis/data.nrc.no/v1").
		Doc("data.nrc.no API")

	ws.Route(ws.PUT(fmt.Sprintf("/tables/{%s}", pathParamTable)).
		Operation("PutTable").
		Doc("Creates or Updates a table").
		Reads(api.Table{}).
		Writes(api.Table{}).
		Consumes("application/json").
		Produces("application/json").
		Param(ws.
			PathParameter(pathParamTable, "table name").
			DataType(stringDataType).
			Required(true)).
		To(restfulPutTable(engine)).
		Returns(http.StatusOK, "OK", api.Table{}))

	ws.Route(ws.GET(fmt.Sprintf("/tables/{%s}/records/{%s}", pathParamTable, pathParamId)).
		Operation("GetRecord").
		Doc("Gets a record").
		Writes(api.Record{}).
		Produces("application/json").
		Param(ws.
			PathParameter(pathParamTable, "table name").
			DataType(stringDataType).
			Required(true)).
		Param(ws.
			PathParameter(pathParamId, "record id").
			DataType(stringDataType).
			Required(true)).
		Param(ws.
			QueryParameter(queryParamRev, "revision").
			DataType(stringDataType).
			Required(false)).
		To(restfulGetRow(engine)).
		Returns(http.StatusOK, "OK", api.Record{}))

	ws.Route(ws.PUT(fmt.Sprintf(`/tables/{%s}/records/{%s}`, pathParamTable, pathParamId)).
		Operation("PutRow").
		Doc("Puts a record in a table").
		Reads(api.PutRecordRequest{}).
		Writes(api.Record{}).
		Consumes("application/json").
		Produces("application/json").
		Param(ws.
			PathParameter(pathParamTable, "table name").
			DataType(stringDataType).
			Required(true)).
		Param(ws.
			PathParameter(pathParamId, "row id").
			DataType(stringDataType).
			Required(true)).
		Param(ws.
			QueryParameter(queryParamIsReplication, "is this a new record?").
			DataType(booleanDataType).
			Required(false)).
		To(restfulPutRow(engine)).
		Returns(http.StatusOK, "OK", api.Record{}))

	ws.Route(ws.GET("/changes").
		Operation("GetChanges").
		Doc("Get changes").
		Writes(api.Changes{}).
		Produces("application/json").
		Param(ws.
			PathParameter(queryParamSince, "checkpoint").
			DataType(intDataType).
			Required(true)).
		To(restfulGetChanges(engine)).
		Returns(http.StatusOK, "OK", api.Changes{}))

	return &Handler{
		engine: engine,
		ws:     ws,
	}
}
