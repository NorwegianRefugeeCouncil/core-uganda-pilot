package record

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/store"
	"net/http"
)

type Handler struct {
	store      store.RecordStore
	formStore      store.FormStore
	webService *restful.WebService
}

func NewHandler(store store.RecordStore, formStore store.FormStore) *Handler {
	h := &Handler{store: store, formStore: formStore}

	ws := new(restful.WebService).Path("/records")
	h.webService = ws

	recordPath := fmt.Sprintf("/{%s}", constants.ParamRecordID)
	ws.Route(ws.PUT(recordPath).To(h.RestfulUpdate).
		Doc("update a record").
		Operation("updateRecord").
		Param(restful.PathParameter(constants.ParamRecordID, "id of the record")).
		Reads(types.Record{}).
		Writes(types.Record{}).
		Returns(http.StatusOK, "OK", types.Record{}))

	ws.Route(ws.POST("/").To(h.RestfulCreate).
		Doc("create a record").
		Operation("createRecord").
		Reads(types.Record{}).
		Writes(types.Record{}).
		Returns(http.StatusOK, "OK", types.Record{}))

	ws.Route(ws.GET("/").To(h.RestfulList).
		Doc("list records").
		Operation("listRecords").
		Param(restful.QueryParameter(constants.ParamDatabaseID, "id of the database").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Param(restful.QueryParameter(constants.ParamFormID, "id of the form").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Writes(types.RecordList{}).
		Returns(http.StatusOK, "OK", types.RecordList{}))

	ws.Route(ws.GET(fmt.Sprintf("/{%s}", constants.ParamRecordID)).To(h.RestfulGet).
		Doc("get record").
		Operation("getRecord").
		Param(restful.QueryParameter(constants.ParamDatabaseID, "id of the database").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Param(restful.QueryParameter(constants.ParamFormID, "id of the form").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Writes(types.Record{}).
		Returns(http.StatusOK, "OK", types.Record{}))

	return h
}

func (h *Handler) WebService() *restful.WebService {
	return h.webService
}
