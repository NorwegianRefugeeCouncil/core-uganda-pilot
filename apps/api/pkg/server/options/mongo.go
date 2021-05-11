package options

import (
	"fmt"
	"github.com/nrc-no/core/apps/api/pkg/registry/generic"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/server"
	serverstorage "github.com/nrc-no/core/apps/api/pkg/server/storage"
	storagebackend "github.com/nrc-no/core/apps/api/pkg/storage/backend"
)

type MongoOptions struct {
	StorageConfig storagebackend.Config
}

func (o *MongoOptions) ApplyWithStorageFactoryTo(factory *serverstorage.DefaultStorageFactory, config *server.Config) error {
	config.RESTOptionsGetter = &StorageFactoryRestOptionsFactory{Options: *o, StorageFactory: factory}
	return nil
}

func NewMongoOptions(backendConfig *storagebackend.Config) *MongoOptions {
	options := &MongoOptions{
		StorageConfig: *backendConfig,
	}
	return options
}

type MongoRestOptionsFactory struct {
	Options MongoOptions
}

func (f *MongoRestOptionsFactory) GetRESTOptions(resource schema.GroupResource) (generic.RESTOptions, error) {
	ret := generic.RESTOptions{
		StorageConfig:  &f.Options.StorageConfig,
		ResourcePrefix: resource.Group + "/" + resource.Resource,
		Decorator:      generic.UndecoratedStorage,
	}
	return ret, nil
}

type StorageFactoryRestOptionsFactory struct {
	Options        MongoOptions
	StorageFactory serverstorage.StorageFactory
}

func (f *StorageFactoryRestOptionsFactory) GetRESTOptions(resource schema.GroupResource) (generic.RESTOptions, error) {
	storageConfig, err := f.StorageFactory.NewConfig(resource)
	if err != nil {
		return generic.RESTOptions{}, fmt.Errorf("unable to find storage destination for %v, due to %v", resource, err.Error())
	}
	ret := generic.RESTOptions{
		StorageConfig: storageConfig,
		Decorator:     generic.UndecoratedStorage,
		//DeleteCollectionWorkers: f.Options.DeleteCollectionWorkers,
		//EnableGarbageCollection: f.Options.EnableGarbageCollection,
		ResourcePrefix: f.StorageFactory.ResourcePrefix(resource),
		//CountMetricPollPeriod:   f.Options.StorageConfig.CountMetricPollPeriod,
	}
	//if f.Options.EnableWatchCache {
	//  sizes, err := ParseWatchCacheSizes(f.Options.WatchCacheSizes)
	//  if err != nil {
	//    return generic.RESTOptions{}, err
	//  }
	//  size, ok := sizes[resource]
	//  if ok && size > 0 {
	//    klog.Warningf("Dropping watch-cache-size for %v - watchCache size is now dynamic", resource)
	//  }
	//  if ok && size <= 0 {
	//    ret.Decorator = generic.UndecoratedStorage
	//  } else {
	//    ret.Decorator = genericregistry.StorageWithCacher()
	//  }
	//}

	return ret, nil
}
