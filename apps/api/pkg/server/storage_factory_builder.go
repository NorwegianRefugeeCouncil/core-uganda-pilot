package server

import (
	"github.com/nrc-no/core/apps/api/pkg/api/defaultscheme"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/storage/backend"
)

type StorageFactoryConfig struct {
	StorageConfig backend.Config
	Serializer    runtime.StorageSerializer
}

func NewStorageFactoryConfig() *StorageFactoryConfig {
	return &StorageFactoryConfig{
		Serializer: defaultscheme.Codecs,
	}
}

func (c *StorageFactoryConfig) Complete(mongoOptions *MongoOptions) {}

type completedStorageFactoryConfig struct {
	*StorageFactoryConfig
}

func (c *completedStorageFactoryConfig) New() {

}
