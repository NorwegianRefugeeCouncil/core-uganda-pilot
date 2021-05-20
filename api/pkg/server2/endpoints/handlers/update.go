package handlers

import (
	"fmt"
	"github.com/nrc-no/core/api/pkg/server2/registry/rest"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"net/http"
)

func UpdateResource(scope *RequestScope, updater rest.Updater) http.HandlerFunc {
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

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		s, err := negotiation.NegotiateInputSerializer(req, false, scope.Serializer)
		if err != nil {
			scope.err(err, w, req)
			return
		}
		defaultGVK := scope.Kind
		original := updater.New()

		decoder := scope.Serializer.DecoderToVersion(s.Serializer, scope.HubGroupVersion)
		obj, gvk, err := decoder.Decode(body, &defaultGVK, original)
		if err != nil {
			err = transformDecodeError(scope.Typer, err, original, gvk, body)
			scope.err(err, w, req)
			return
		}

		if !scope.AcceptsGroupVersion(gvk.GroupVersion()) {
			err = errors.NewBadRequest(fmt.Sprintf("the API version in the data (%s) does not match the expected API version (%s)", gvk.GroupVersion(), defaultGVK.GroupVersion()))
			scope.err(err, w, req)
			return
		}

		obj, err = updater.Update(ctx, name, rest.DefaultUpdatedObjectInfo(obj))
		if err != nil {
			scope.err(err, w, req)
			return
		}

		status := http.StatusOK
		transformResponseObject(ctx, scope, req, w, status, outputMediaType, obj)

	}
}
