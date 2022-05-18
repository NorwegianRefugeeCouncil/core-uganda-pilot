package handler

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/server/data/api"
)

const (
	stringDataType           = "string"
	intDataType              = "int"
	booleanDataType          = "boolean"
	pathParamTableName       = "table"
	pathParamRecordId        = "id"
	queryParamRecordRevision = "rev"
	queryParamCheckpoint     = "since"
	queryParamIsReplication  = "replication"
	mimeTypeApplicationJson  = "application/json"
)

// Handler is the HTTP Handler for the database API
type Handler struct {
	engine api.Engine
	ws     *restful.WebService
}

// WebService returns the restful.WebService for the database API
func (h *Handler) WebService() *restful.WebService {
	return h.ws
}

// NewHandler creates a new Handler
func NewHandler(engine api.Engine) *Handler {

	ws := new(restful.WebService).
		Path("/apis/data.nrc.no/v1").
		Doc("data.nrc.no API")

	ws.Route(ws.PUT(fmt.Sprintf("/tables/{%s}", pathParamTableName)).
		Operation("PutTable").
		Doc("Creates or Updates a table").
		Reads(api.Table{}).
		Writes(api.Table{}).
		Consumes(mimeTypeApplicationJson).
		Produces(mimeTypeApplicationJson).
		Param(tableNameParam(ws).Required(true)).
		To(restfulPutTable(engine)).
		Returns(http.StatusOK, "OK", api.Table{}))

	ws.Route(ws.GET("/tables").
		Operation("GetTables").
		Doc("Gets tables").
		Writes(api.TableList{}).
		Produces(mimeTypeApplicationJson).
		To(restfulGetTables(engine)).
		Returns(http.StatusOK, "OK", api.TableList{}))

	ws.Route(ws.GET("/tables/{"+pathParamTableName+"}").
		Operation("GetTable").
		Doc("Get Table").
		Writes(api.Table{}).
		Produces(mimeTypeApplicationJson).
		To(restfulGetTable(engine)).
		Returns(http.StatusOK, "OK", api.Table{}))

	ws.Route(ws.GET("/tables/{"+pathParamTableName+"}/records").
		Operation("GetRecords").
		Doc("Gets Records").
		Writes(api.RecordList{}).
		Produces(mimeTypeApplicationJson).
		To(restfulGetRecords(engine)).
		Returns(http.StatusOK, "OK", api.RecordList{}))

	ws.Route(ws.GET(fmt.Sprintf("/tables/{%s}/records/{%s}", pathParamTableName, pathParamRecordId)).
		Operation("GetRecord").
		Doc("Gets a record").
		Writes(api.Record{}).
		Produces(mimeTypeApplicationJson).
		Param(tableNameParam(ws).Required(true)).
		Param(recordIdParam(ws).Required(true)).
		Param(revisionParam(ws).Required(false)).
		To(restfulGetRecord(engine)).
		Returns(http.StatusOK, "OK", api.Record{}))

	ws.Route(ws.PUT(fmt.Sprintf(`/tables/{%s}/records/{%s}`, pathParamTableName, pathParamRecordId)).
		Operation("PutRow").
		Doc("Puts a record in a table").
		Reads(api.PutRecordRequest{}).
		Writes(api.Record{}).
		Consumes(mimeTypeApplicationJson).
		Produces(mimeTypeApplicationJson).
		Param(tableNameParam(ws).Required(true)).
		Param(recordIdParam(ws).Required(true)).
		Param(ws.
			QueryParameter(queryParamIsReplication, "is this a new record?").
			DataType(booleanDataType).
			Required(false)).
		To(restfulPutRecord(engine)).
		Returns(http.StatusOK, "OK", api.Record{}))

	ws.Route(ws.GET("/changes").
		Operation("GetChanges").
		Doc("Get changes").
		Writes(api.Changes{}).
		Produces(mimeTypeApplicationJson).
		Param(checkpointParam(ws).Required(true)).
		To(restfulGetChanges(engine)).
		Returns(http.StatusOK, "OK", api.Changes{}))

	return &Handler{
		engine: engine,
		ws:     ws,
	}
}

func checkpointParam(ws *restful.WebService) *restful.Parameter {
	return ws.
		PathParameter(queryParamCheckpoint, "Checkpoint").
		DataType(intDataType)
}

func revisionParam(ws *restful.WebService) *restful.Parameter {
	return ws.
		QueryParameter(queryParamRecordRevision, "Record Revision").
		DataType(stringDataType)
}

func recordIdParam(ws *restful.WebService) *restful.Parameter {
	return ws.
		PathParameter(pathParamRecordId, "Record ID").
		DataType(stringDataType)
}

func tableNameParam(ws *restful.WebService) *restful.Parameter {
	return ws.
		PathParameter(pathParamTableName, "Table Name").
		DataType(stringDataType)
}
