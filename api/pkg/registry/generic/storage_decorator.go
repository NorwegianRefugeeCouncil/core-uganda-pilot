package generic

import (
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/storage"
	"github.com/nrc-no/core/apps/api/pkg/storage/backend"
)

type StorageDecorator func(
	config *backend.Config,
	resourcePrefix string,
	keyFunc func(obj runtime.Object) (string, error),
	newFunc func() runtime.Object,
	newListFunc func() runtime.Object,
	//getAttrsFunc AttrFunc,
	// trigger storage.IndexerFuncs,
	//indexers *cache.Indexers,
) (storage.Interface, backend.DestroyFunc, error)

// UndecoratedStorage returns the given a new storage from the given config
// without any decoration.
func UndecoratedStorage(
	config *backend.Config,
	resourcePrefix string,
	keyFunc func(obj runtime.Object) (string, error),
	newFunc func() runtime.Object,
	newListFunc func() runtime.Object,
	//getAttrsFunc storage.AttrFunc,
	//trigger storage.IndexerFuncs,
	//indexers *cache.Indexers,
) (storage.Interface, backend.DestroyFunc, error) {
	return NewRawStorage(config, newFunc)
}

// NewRawStorage creates the low level kv storage. This is a work-around for current
// two layer of same storage interface.
// TODO: Once cacher is enabled on all registries (event registry is special), we will remove this method.
func NewRawStorage(config *backend.Config, newFunc func() runtime.Object) (storage.Interface, backend.DestroyFunc, error) {
	return backend.Create(*config, newFunc)
}
