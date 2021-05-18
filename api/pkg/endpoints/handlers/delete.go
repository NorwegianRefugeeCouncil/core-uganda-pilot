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

func DeleteResource(r rest.Deleter, scope *RequestScope) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		outputMediaType, _, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		options := &metav1.DeleteOptions{}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		if len(body) > 0 {

			s, err := negotiation.NegotiateInputSerializer(req, false, scheme.Codecs)
			if err != nil {
				scope.Error(err, w, req)
				return
			}

			gvk := metav1.SchemeGroupVersion.WithKind("DeleteOptions")
			obj, _, err := scheme.Codecs.DecoderToVersion(s.Serializer, gvk.GroupVersion()).Decode(body, &gvk, options)
			if err != nil {
				scope.Error(err, w, req)
				return
			}
			if obj != options {
				scope.Error(fmt.Errorf("decoded object cannot be converted to DeleteOptions"), w, req)
				return
			}

		} else {
			if err := scheme.ParameterCodec.DecodeParameters(req.URL.Query(), metav1.SchemeGroupVersion, options); err != nil {
				err = exceptions.NewBadRequest(err.Error())
				scope.Error(err, w, req)
				return
			}
		}

		ctx := req.Context()

		options.TypeMeta.SetGroupVersionKind(metav1.SchemeGroupVersion.WithKind("DeleteOptions"))
		wasDeleted := true

		_, key, err := scope.Namer.Name(req)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		obj, deleted, err := r.Delete(ctx, key, func(ctx context.Context, obj runtime.Object) error {
			return nil
		})
		wasDeleted = deleted

		if err != nil {
			scope.Error(err, w, req)
			return
		}

		status := http.StatusOK
		if !wasDeleted {
			status = http.StatusAccepted
		}

		if obj == nil {
			obj = &metav1.Status{
				Status: metav1.StatusSuccess,
				Code:   int32(status),
				Details: &metav1.StatusDetails{
					UID:   key,
					Group: scope.Kind.Group,
					Kind:  scope.Kind.Kind,
				},
			}
		}

		transformResponseObject(ctx, scope, req, w, status, outputMediaType, obj)

	}
}
