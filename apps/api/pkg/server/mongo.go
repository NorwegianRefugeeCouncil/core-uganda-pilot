package server

import (
	"github.com/nrc-no/core/apps/api/pkg/registry/generic"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/storage/backend"
)

type MongoOptions struct {
	StorageConfig backend.Config
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
