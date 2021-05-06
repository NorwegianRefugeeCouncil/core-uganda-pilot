package storage

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkgs2/runtime"
)

type UpdateFunc func(obj interface{}) (interface{}, error)

type Interface interface {
	Create(ctx context.Context, key string, obj, out runtime.Object) error
	Get(ctx context.Context, key string, out runtime.Object) error
	Update(ctx context.Context, key string, objType runtime.Object, updateFunc UpdateFunc) error
}
