package store

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
)

type Interface interface {
	Create(ctx context.Context, key string, obj, out runtime.Object) error
	Get(ctx context.Context, key string, getOptions GetOptions, objPtr runtime.Object) error
}

type GetOptions struct {
	IgnoreNotFound  bool
	ResourceVersion string
}
