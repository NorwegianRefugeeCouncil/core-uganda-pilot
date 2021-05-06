package storage

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg3/runtime"
)

type Interface interface {
	Get(ctx context.Context, key string, out runtime.Object) error
	Create(ctx context.Context, key string, in, out runtime.Object) error
	Update(ctx context.Context, key string, in, out runtime.Object) error
}
