package handlers

import (
	"context"
	metainternalscheme "github.com/nrc-no/core/apps/api/pkg/apis/meta/internalversion/scheme"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/negotiation"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
	"net/http"
)

type getterFunc func(ctx context.Context, name string, req *http.Request) (runtime.Object, error)

func getResourceHandler(scope *RequestScope, getter getterFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		_, name, err := scope.Namer.Name(req)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		outputMediaType, _, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		ctx := req.Context()
		result, err := getter(ctx, name, req)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		transformResponseObject(ctx, scope, req, w, http.StatusOK, outputMediaType, result)

	}
}

func GetResource(r rest.Getter, scope *RequestScope) http.HandlerFunc {
	return getResourceHandler(scope, func(ctx context.Context, name string, req *http.Request) (runtime.Object, error) {
		options := metav1.GetOptions{}
		if values := req.URL.Query(); len(values) > 0 {
			if err := metainternalscheme.ParameterCodec.DecodeParameters(values, metav1.SchemeGroupVersion, &options); err != nil {
				err = exceptions.NewBadRequest(err.Error())
				return nil, err
			}
		}
		return r.Get(ctx, name, &options)
	})
}
