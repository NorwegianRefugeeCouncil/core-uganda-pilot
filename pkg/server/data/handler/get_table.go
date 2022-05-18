package handler

import (
	"encoding/json"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/server/data/api"
)

func getTable(e api.Engine, tableName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		table, err := e.GetTable(req.Context(), api.GetTableRequest{TableName: tableName})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseBytes, err := json.Marshal(table)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		header := w.Header()
		header.Set("Content-Type", "application/json")
		w.Write(responseBytes)
	})
}

func restfulGetTable(engine api.Engine) func(req *restful.Request, resp *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		tableName := req.PathParameter(pathParamTableName)
		getTable(engine, tableName).ServeHTTP(resp.ResponseWriter, req.Request)
	}
}
