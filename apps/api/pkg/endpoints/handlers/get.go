package handlers

import (
  "context"
  "github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/negotiation"
  "github.com/nrc-no/core/apps/api/pkg/registry/rest"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/apimachinery/pkg/runtime"
  "net/http"
)

type getterFunc func(ctx context.Context, req *http.Request) (runtime.Object, error)

func getResourceHandler(scope *RequestScope, getter getterFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {

    ctx := req.Context()

    outputMediaType, _, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
    if err != nil {
      scope.err(err, w, req)
      return
    }

    result, err := getter(ctx, req)
    if err != nil {
      scope.err(err, w, req)
      return
    }

    transformResponseObject(ctx, scope, req, w, http.StatusOK, outputMediaType, result)

  }
}

func GetResource(r rest.Getter, scope *RequestScope) http.HandlerFunc {
  return getResourceHandler(scope, func(ctx context.Context, req *http.Request) (runtime.Object, error) {
    options := metav1.GetOptions{}
    return r.Get(ctx, "", &options)
  })
}
