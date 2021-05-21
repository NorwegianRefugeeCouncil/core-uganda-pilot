package generic

import (
	store2 "github.com/nrc-no/core/api/pkg/store"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type RESTOptions struct {
	StorageConfig *store2.Config
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
