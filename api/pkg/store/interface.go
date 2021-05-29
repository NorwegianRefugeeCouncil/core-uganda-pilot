package store

import (
	"context"
	"github.com/nrc-no/core/api/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

// Interface represents the generic store interface
type Interface interface {
	Get(ctx context.Context, key string, options GetOptions, into runtime.Object) error
	Create(ctx context.Context, key string, options CreateOptions, obj, out runtime.Object) error
	List(ctx context.Context, key string, options ListOptions, listObj runtime.Object) error
	Update(ctx context.Context, key string, options UpdateOptions, out runtime.Object, tryUpdate UpdateFunc) error
	Delete(ctx context.Context, key string, options DeleteOptions, out runtime.Object) error
	Watch(ctx context.Context, key string, options ListOptions) (watch.Interface, error)
}

// GetOptions placeholder to put store options for GET requests that return a single result
type GetOptions struct{}

// CreateOptions placeholder to put store options for PUT requests
type CreateOptions struct{}

// ListOptions placeholder to put store options for GET requests that return a list result
type ListOptions struct {
	Limit *int64

	Selector fields.Selector

	ResourceVersion string
	// SyncOnly will emit events since the given ResourceVersion (or since beginning if not provided)
	// and will close the channel
	SyncOnly bool
}

// UpdateOptions placeholder to put store options for PUT request
type UpdateOptions struct{}

// UpdateOptions placeholder to put store options for PUT request
type DeleteOptions struct{}

type UpdateFunc func(input runtime.Object) (output runtime.Object, err error)
