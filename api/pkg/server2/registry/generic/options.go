package generic

import (
	"github.com/nrc-no/core/api/pkg/server2/store"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type RESTOptions struct {
	StorageConfig *store.Config
}

func (opts RESTOptions) GetRESTOptions(schema.GroupResource) (RESTOptions, error) {
	return opts, nil
}

type RESTOptionsGetter interface {
	GetRESTOptions(resource schema.GroupResource) (RESTOptions, error)
}

type StoreOptions struct {
	RESTOptions RESTOptionsGetter
}
