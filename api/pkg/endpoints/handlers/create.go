package handlers

import (
	"fmt"
	rest2 "github.com/nrc-no/core/api/pkg/registry/rest"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"net/http"
)

// CreateResource is a generic REST handler to create resources
func CreateResource(scope *RequestScope, creater rest2.Creater) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		// Retrieve the resource name
		name, err := scope.Namer.Name(req)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		// Negotiate output media type
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

		// Negotiate input media type
		s, err := negotiation.NegotiateInputSerializer(req, false, scope.Serializer)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		// Decode request body
		decoder := scope.Serializer.DecoderToVersion(s.Serializer, scope.HubGroupVersion)
		defaultGVK := scope.Kind
		original := creater.New()
		inputObj, inputGvk, err := decoder.Decode(body, &defaultGVK, original)
		if err != nil {
			err = transformDecodeError(scope.Typer, err, original, inputGvk, body)
			scope.err(err, w, req)
			return
		}

		// Make sure we accept the body api GroupVersion
		if !scope.AcceptsGroupVersion(inputGvk.GroupVersion()) {
			err = errors.NewBadRequest(fmt.Sprintf("the API version in the data (%s) does not match the expected API version (%s)", inputGvk.GroupVersion(), defaultGVK.GroupVersion()))
			scope.err(err, w, req)
			return
		}

		// Create the resource
		inputObj, err = creater.Create(ctx, name, inputObj)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		// Respond
		status := http.StatusCreated
		transformResponseObject(ctx, scope, req, w, status, outputMediaType, inputObj)

	}
}
