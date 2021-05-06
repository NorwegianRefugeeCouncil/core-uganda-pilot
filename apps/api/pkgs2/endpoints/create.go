package endpoints

import (
	"github.com/emicklei/go-restful"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"io/ioutil"
	"net/http"
)

func restfulCreateResource(r Creater, scope RequestScope) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {

	}
}

func CreateResource(r rest.Creater, scope *RequestScope) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		gv := scope.Kind.GroupVersion()
		bodyBytes, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
