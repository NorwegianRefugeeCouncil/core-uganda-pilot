package handlers

import (
	"github.com/nrc-no/core/api/pkg/server2/registry/rest"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"net/http"
)

// GetResource is a generic REST handler to list resources (multiple result)
func ListResource(scope *RequestScope, getter rest.Lister) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		outputMediaType, _, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		result, err := getter.List(ctx)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		transformResponseObject(ctx, scope, req, w, http.StatusOK, outputMediaType, result)

	}
}
