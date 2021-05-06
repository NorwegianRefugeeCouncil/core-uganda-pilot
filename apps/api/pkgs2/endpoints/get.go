package endpoints

import (
	"context"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/nrc-no/core/apps/api/pkgs2/runtime"
	"net/http"
)

func restfulGetResource(r Getter, scope RequestScope) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		GetResource(r, &scope)(response.ResponseWriter, request.Request)
	}
}

func GetResource(getter Getter, scope *RequestScope) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		result, err := getter.Get(ctx, "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bytes, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	}
}

type getterFunc func(ctx context.Context, id string, req *http.Request) (runtime.Object, error)

func getResourceHandler(scope *RequestScope, getter getterFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
