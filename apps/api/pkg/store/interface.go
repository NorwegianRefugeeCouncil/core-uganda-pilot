package store

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
)

type Interface interface {
	Create(ctx context.Context, key string, obj, out runtime.Object) error
	Get(ctx context.Context, key string, getOptions GetOptions, objPtr runtime.Object) error
	Update(ctx context.Context, key string, objType runtime.Object, updateFunc UpdateFunc, updateOptions UpdateOptions) error
}

type UpdateFunc func(obj runtime.Object) (runtime.Object, error)

type UpdateOptions struct {
  IgnoreNotFound bool
}

type GetOptions struct {
	IgnoreNotFound  bool
	ResourceVersion string
}
