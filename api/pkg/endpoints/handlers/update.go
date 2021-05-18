package handlers

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/apps/api/pkg/apis/meta"
	"github.com/nrc-no/core/apps/api/pkg/apis/meta/internalversion/scheme"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/negotiation"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
	"io/ioutil"
	"net/http"
)

func UpdateResource(r rest.Updater, scope *RequestScope) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		outputMediaType, _, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		options := &metav1.UpdateOptions{}
		if err := scheme.ParameterCodec.DecodeParameters(req.URL.Query(), metav1.SchemeGroupVersion, options); err != nil {
			err = exceptions.NewBadRequest(err.Error())
			scope.Error(err, w, req)
			return
		}
		options.TypeMeta.SetGroupVersionKind(metav1.SchemeGroupVersion.WithKind("UpdateOptions"))

		s, err := negotiation.NegotiateInputSerializer(req, false, scope.Serializer)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		defaultGVK := scope.Kind
		original := r.New()

		decoder := scope.Serializer.DecoderToVersion(s.Serializer, scope.HubGroupVersion)
		obj, gvk, err := decoder.Decode(bodyBytes, &defaultGVK, original)
		if err != nil {
			err = transformDecodeError(scope.Typer, err, original, gvk, bodyBytes)
			scope.Error(err, w, req)
			return
		}

		if !scope.AcceptsGroupVersion(gvk.GroupVersion()) {
			err = exceptions.NewBadRequest(fmt.Sprintf("the API version in the data (%s) does not match the expected API version (%s)", gvk.GroupVersion(), defaultGVK.GroupVersion()))
			scope.Error(err, w, req)
			return
		}

		ctx := req.Context()

		accessor, err := meta.Accessor(obj)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		out, created, err := r.Update(
			ctx,
			accessor.GetUID(),
			DefaultUpdatedObjectInfo(obj),
			func(ctx context.Context, obj runtime.Object) error {
				return nil
			}, func(ctx context.Context, obj, old runtime.Object) error {
				return nil
			},
			true,
			options)

		if err != nil {
			scope.Error(err, w, req)
			return
		}

		statusCode := http.StatusOK
		if created {
			statusCode = http.StatusCreated
		}

		transformResponseObject(ctx, scope, req, w, statusCode, outputMediaType, out)

	}
}

// TransformFunc is a function to transform and return newObj
type TransformFunc func(ctx context.Context, newObj runtime.Object, oldObj runtime.Object) (transformedNewObj runtime.Object, err error)

// defaultUpdatedObjectInfo implements UpdatedObjectInfo
type defaultUpdatedObjectInfo struct {
	// obj is the updated object
	obj runtime.Object

	// transformers is an optional list of transforming functions that modify or
	// replace obj using information from the context, old object, or other sources.
	transformers []TransformFunc
}

// DefaultUpdatedObjectInfo returns an UpdatedObjectInfo impl based on the specified object.
func DefaultUpdatedObjectInfo(obj runtime.Object, transformers ...TransformFunc) rest.UpdatedObjectInfo {
	return &defaultUpdatedObjectInfo{obj, transformers}
}

// UpdatedObject satisfies the UpdatedObjectInfo interface.
// It returns a copy of the held obj, passed through any configured transformers.
func (i *defaultUpdatedObjectInfo) UpdatedObject(ctx context.Context, oldObj runtime.Object) (runtime.Object, error) {
	var err error
	// Start with the configured object
	newObj := i.obj

	// If the original is non-nil (might be nil if the first transformer builds the object from the oldObj), make a copy,
	// so we don't return the original. BeforeUpdate can mutate the returned object, doing things like clearing ResourceVersion.
	// If we're re-called, we need to be able to return the pristine version.
	if newObj != nil {
		newObj = newObj.DeepCopyObject()
	}

	// Allow any configured transformers to update the new object
	for _, transformer := range i.transformers {
		newObj, err = transformer(ctx, newObj, oldObj)
		if err != nil {
			return nil, err
		}
	}

	return newObj, nil
}
