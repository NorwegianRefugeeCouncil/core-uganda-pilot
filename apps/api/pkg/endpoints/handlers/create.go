package handlers

import (
  "context"
  "fmt"
  "github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/negotiation"
  "github.com/nrc-no/core/apps/api/pkg/registry/rest"
  "k8s.io/apimachinery/pkg/api/errors"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/apimachinery/pkg/runtime"
  "net/http"
)

func CreateResource(r rest.Creater, scope *RequestScope) http.HandlerFunc {
  return createHandler(&namedCreaterAdapter{r}, scope)
}

type namedCreaterAdapter struct {
  rest.Creater
}

func (c *namedCreaterAdapter) Create(ctx context.Context, obj runtime.Object, createValidatingAdmission rest.ValidateObjectFunc) (runtime.Object, error) {
  return c.Creater.Create(ctx, obj, createValidatingAdmission)
}

func createHandler(r rest.Creater, scope *RequestScope) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {

    ctx, cancel := context.WithTimeout(req.Context(), requestTimeoutUpperBound)
    defer cancel()
    outputMediaType, _, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
    if err != nil {
      scope.err(err, w, req)
    }

    gv := scope.Kind.GroupVersion()
    s, err := negotiation.NegotiateInputSerializer(req, false, scope.Serializer)
    if err != nil {
      scope.err(err, w, req)
      return
    }

    decoder := scope.Serializer.DecoderToVersion(s.Serializer, scope.HubGroupVersion)
    body, err := limitedReadBody(req, scope.MaxRequestBodyBytes)
    if err != nil {
      scope.err(err, w, req)
      return
    }

    defaultGVK := scope.Kind
    original := r.New()
    obj, gvk, err := decoder.Decode(body, &defaultGVK, original)
    if err != nil {
      err = transformDecodeError(scope.Typer, err, original, gvk, body)
      scope.err(err, w, req)
      return
    }
    if !scope.AcceptsGroupVersion(gvk.GroupVersion()) {
      err = errors.NewBadRequest(fmt.Sprintf("the API version in the data (%s) does not match the expected API version (%v)", gvk.GroupVersion().String(), gv.String()))
      scope.err(err, w, req)
      return
    }

    requestFunc := func() (runtime.Object, error) {
      return r.Create(
        ctx,
        obj,
        func(ctx context.Context, obj runtime.Object) error {
          return nil
        },
      )
    }

    result, err := requestFunc()
    if err != nil {
      scope.err(err, w, req)
      return
    }

    code := http.StatusCreated
    status, ok := result.(*metav1.Status)
    if ok && status.Code == 0 {
      status.Code = int32(code)
    }

    transformResponseObject(ctx, scope, req, w, code, outputMediaType, result)

  }
}
