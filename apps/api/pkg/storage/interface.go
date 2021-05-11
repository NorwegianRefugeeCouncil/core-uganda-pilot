package storage

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/watch"
)

type ResponseMeta struct {
	ResourceVersion uint64
}

type UpdateFunc func(input runtime.Object, res ResponseMeta) (output runtime.Object, ttl *uint64, err error)

type GetOptions struct {
	IgnoreNotFound  bool
	ResourceVersion string
}

type ListOptions struct {
	ResourceVersion string
}

type Interface interface {
	Versioner() Versioner
	Get(ctx context.Context, key string, getOptions GetOptions, out runtime.Object) error
	List(ctx context.Context, listOptions ListOptions, out runtime.Object) error
	Create(ctx context.Context, obj runtime.Object) error
	Update(ctx context.Context, key string, out runtime.Object, update UpdateFunc) error
	Watch(ctx context.Context, key string, opts ListOptions) (watch.Interface, error)
}

type Versioner interface {
	UpdateObject(obj runtime.Object, resourceVersion uint64) error
	UpdateList(obj runtime.Object, resourceVersion uint64, continueValue string, remainingItemCount *int64) error
	PrepareObjectForStorage(obj runtime.Object) error
	ObjectResourceVersion(obj runtime.Object) (uint64, error)
	ParseResourceVersion(resourceVersion string) (uint64, error)
}
