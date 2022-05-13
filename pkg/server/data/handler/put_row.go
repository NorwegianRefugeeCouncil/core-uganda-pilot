package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/server/data/api"
)

func putRow(e api.Engine) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var request api.PutRecordRequest
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(bodyBytes, &request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ret, err := e.PutRecord(req.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseBytes, err := json.Marshal(ret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		header := w.Header()
		header.Set("Content-Type", "application/json")
		header.Set("ETag", fmt.Sprintf("%s", ret.Revision.String()))
		w.Write(responseBytes)
	})
}

func restfulPutRow(engine api.Engine) func(req *restful.Request, resp *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		putRow(engine).ServeHTTP(resp.ResponseWriter, req.Request)
	}
}
