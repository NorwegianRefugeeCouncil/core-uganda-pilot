package handlers

import (
	"fmt"
	rest2 "github.com/nrc-no/core/api/pkg/registry/rest"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"net/http"
)

// UpdateResource is a generic REST handler for updating (PUTting) resources
func UpdateResource(scope *RequestScope, updater rest2.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		// Retrieve the resource name
		name, err := scope.Namer.Name(req)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		// Negotiate the output media type
		outputMediaType, _, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		// Read request body
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		// Negotiate output media type
		s, err := negotiation.NegotiateInputSerializer(req, false, scope.Serializer)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		// Decode the body to the Hub api version
		decoder := scope.Serializer.DecoderToVersion(s.Serializer, scope.HubGroupVersion)
		defaultGVK := scope.Kind
		original := updater.New()
		inputObj, inputGvk, err := decoder.Decode(body, &defaultGVK, original)
		if err != nil {
			err = transformDecodeError(scope.Typer, err, original, inputGvk, body)
			scope.err(err, w, req)
			return
		}

		// Make sure we accept the input group version
		if !scope.AcceptsGroupVersion(inputGvk.GroupVersion()) {
			err = errors.NewBadRequest(fmt.Sprintf("the API version in the data (%s) does not match the expected API version (%s)", inputGvk.GroupVersion(), defaultGVK.GroupVersion()))
			scope.err(err, w, req)
			return
		}

		// Update the resource
		inputObj, err = updater.Update(ctx, name, rest2.DefaultUpdatedObjectInfo(inputObj))
		if err != nil {
			scope.err(err, w, req)
			return
		}

		// Respond
		status := http.StatusOK
		transformResponseObject(scope, req, w, status, outputMediaType, inputObj)

	}
}
