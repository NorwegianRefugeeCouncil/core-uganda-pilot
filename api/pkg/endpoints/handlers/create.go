package handlers

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/apps/api/pkg/apis/meta/internalversion/scheme"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/negotiation"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
	"io/ioutil"
	"net/http"
)

func createHandler(r rest.Creater, scope *RequestScope) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		outputMediaType, _, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		gv := scope.Kind.GroupVersion()
		s, err := negotiation.NegotiateInputSerializer(req, false, scope.Serializer)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		decoder := scope.Serializer.DecoderToVersion(s.Serializer, scope.HubGroupVersion)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		options := &metav1.CreateOptions{}
		values := req.URL.Query()
		if err := scheme.ParameterCodec.DecodeParameters(values, metav1.SchemeGroupVersion, options); err != nil {
			err = exceptions.NewBadRequest(err.Error())
			scope.Error(err, w, req)
			return
		}
		options.TypeMeta.SetGroupVersionKind(metav1.SchemeGroupVersion.WithKind("CreateOptions"))
		defaultGvk := scope.Kind
		original := r.New()
		obj, gvk, err := decoder.Decode(body, &defaultGvk, original)
		if err != nil {
			err = transformDecodeError(scope.Typer, err, original, gvk, body)
			scope.Error(err, w, req)
			return
		}

		if !scope.AcceptsGroupVersion(gvk.GroupVersion()) {
			err = exceptions.NewBadRequest(fmt.Sprintf("the API version in the data (%s) does not match the expected API version (%s)", gvk.GroupVersion().String(), gv.String()))
			scope.Error(err, w, req)
			return
		}

		ctx := req.Context()

		result, err := r.Create(ctx, obj, func(ctx context.Context, obj runtime.Object) error {
			return nil
		}, options)

		if err != nil {
			scope.Error(err, w, req)
			return
		}

		code := http.StatusCreated
		status, ok := result.(*metav1.Status)
		if ok && err == nil && status.Code == 0 {
			status.Code = int32(code)
		}

		transformResponseObject(ctx, scope, req, w, code, outputMediaType, result)

	}
}

func CreateResource(r rest.Creater, scope *RequestScope) http.HandlerFunc {
	return createHandler(r, scope)
}
