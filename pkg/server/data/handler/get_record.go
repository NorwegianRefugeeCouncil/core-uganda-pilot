package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/server/data/api"
)

func getRecord(e api.Engine, request api.GetRecordRequest) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rec, err := e.GetRecord(req.Context(), request)
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
		header.Set("ETag", fmt.Sprintf("%s", rec.Revision.String()))
		w.Write(responseBytes)
	})
}

func restfulGetRecord(e api.Engine) func(req *restful.Request, resp *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		rev, err := api.ParseRevision(req.QueryParameter(queryParamRev))
		if err != nil {
			http.Error(resp.ResponseWriter, err.Error(), http.StatusBadRequest)
			return
		}
		var request = api.GetRecordRequest{
			RecordID:  req.PathParameter(pathParamId),
			TableName: req.PathParameter(pathParamTable),
			Revision:  rev,
		}
		getRecord(e, request).ServeHTTP(resp.ResponseWriter, req.Request)
	}
}
