package storage

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
)

type Interface interface {
	Get(ctx context.Context, key string, out runtime.Object) error
	List(ctx context.Context, out runtime.Object) error
	Create(ctx context.Context, in, out runtime.Object) error
	Update(ctx context.Context, key string, in, out runtime.Object) error
	Watch(ctx context.Context, objPtr runtime.Object, watchFunc func(eventType string, obj runtime.Object)) error
}
