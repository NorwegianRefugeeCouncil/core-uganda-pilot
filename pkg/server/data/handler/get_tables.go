package handler

import (
	"encoding/json"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/server/data/api"
)

func getTables(e api.Engine) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		tables, err := e.GetTables(req.Context(), api.GetTablesRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseBytes, err := json.Marshal(tables)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		header := w.Header()
		header.Set("Content-Type", "application/json")
		w.Write(responseBytes)
	})
}

func restfulGetTables(engine api.Engine) func(req *restful.Request, resp *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		getTables(engine).ServeHTTP(resp.ResponseWriter, req.Request)
	}
}
