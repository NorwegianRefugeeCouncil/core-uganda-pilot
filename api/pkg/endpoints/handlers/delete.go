package handlers

import (
	rest2 "github.com/nrc-no/core/api/pkg/registry/rest"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"net/http"
)

// DeleteResource is a generic REST handler to delete resoures
func DeleteResource(scope *RequestScope, deleter rest2.Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		name, err := scope.Namer.Name(req)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		outputMediaType, _, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		result, _, err := deleter.Delete(ctx, name)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		transformResponseObject(scope, req, w, http.StatusOK, outputMediaType, result)

	}
}
