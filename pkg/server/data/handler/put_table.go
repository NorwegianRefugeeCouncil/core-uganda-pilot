package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/server/data/api"
)

func putTable(e api.Engine) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var table api.Table
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(bodyBytes, &table); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, err := e.CreateTable(req.Context(), table); err != nil {
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

func restfulPutTable(engine api.Engine) func(req *restful.Request, resp *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		putTable(engine).ServeHTTP(resp.ResponseWriter, req.Request)
	}
}
