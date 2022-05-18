package handler

import (
	"encoding/json"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/server/data/api"
)

func getRecords(e api.Engine, request api.GetRecordsRequest) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rec, err := e.GetRecords(req.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseBytes, err := json.Marshal(rec)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		header := w.Header()
		header.Set("Content-Type", "application/json")
		w.Write(responseBytes)
	})
}

func restfulGetRecords(e api.Engine) func(req *restful.Request, resp *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		revStr := req.QueryParameter("revisions")
		var request = api.GetRecordsRequest{
			TableName: req.PathParameter(pathParamTableName),
			Revisions: revStr == "true",
		}
		getRecords(e, request).ServeHTTP(resp.ResponseWriter, req.Request)
	}
}
