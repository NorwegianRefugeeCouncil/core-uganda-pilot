package rest

import (
  "context"
  "github.com/nrc-no/core/apps/api/apis/meta"
  "io"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/apimachinery/pkg/runtime"
  "k8s.io/apimachinery/pkg/runtime/schema"
)

type ValidateObjectFunc func(ctx context.Context, obj runtime.Object) error

type Storage interface {
  New() runtime.Object
}

type Creater interface {
  New() runtime.Object
  Create(ctx context.Context, obj runtime.Object, createValidation ValidateObjectFunc) (runtime.Object, error)
}

type Lister interface {
  NewList() runtime.Object
  List(ctx context.Context, options meta.ListOptions) (runtime.Object, error)
}

type Getter interface {
  Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error)
}

type Updater interface {
  New() runtime.Object
  Update(ctx context.Context, name string, updateValidation ValidateObjectFunc) (runtime.Object, error)
}

type TableConvertor interface {
  ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error)
}

type ResourceStreamer interface {
  // InputStream should return an io.ReadCloser if the provided object supports streaming. The desired
  // api version and an accept header (may be empty) are passed to the call. If no error occurs,
  // the caller may return a flag indicating whether the result should be flushed as writes occur
  // and a content type string that indicates the type of the stream.
  // If a null stream is returned, a StatusNoContent response wil be generated.
  InputStream(ctx context.Context, apiVersion, acceptHeader string) (stream io.ReadCloser, flush bool, mimeType string, err error)
}

type GroupVersionKindProvider interface {
  GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind
}
