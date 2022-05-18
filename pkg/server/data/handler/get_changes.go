package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/server/data/api"
)

func getChanges(e api.Engine, request api.GetChangesRequest) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		changeStream, err := e.GetChanges(req.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonBytes, err := json.Marshal(changeStream)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)
	})
}

func restfulGetChanges(e api.Engine) func(req *restful.Request, resp *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		var request api.GetChangesRequest
		sinceStr := req.QueryParameter(queryParamCheckpoint)
		since, err := strconv.ParseInt(sinceStr, 10, 64)
		if err != nil {
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}
		request.Since = since
		getChanges(e, request).ServeHTTP(resp.ResponseWriter, req.Request)
	}
}
