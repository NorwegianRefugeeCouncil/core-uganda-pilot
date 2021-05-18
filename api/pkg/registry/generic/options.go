package generic

import (
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/storage/backend"
)

type RESTOptions struct {
	StorageConfig  *backend.Config
	ResourcePrefix string
	Decorator      StorageDecorator
}

func (opts RESTOptions) GetRESTOptions(resource schema.GroupResource) (RESTOptions, error) {
	return opts, nil
}

type RESTOptionsGetter interface {
	GetRESTOptions(resource schema.GroupResource) (RESTOptions, error)
}

type StoreOptions struct {
	RESTOptions RESTOptionsGetter
}
