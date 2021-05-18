package nrc_apiserver

import (
	"github.com/nrc-no/core/apps/api/pkg/api/defaultscheme"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/server/options"
	"github.com/nrc-no/core/apps/api/pkg/server/resourceconfig"
	"github.com/nrc-no/core/apps/api/pkg/server/storage"
	storagebackend "github.com/nrc-no/core/apps/api/pkg/storage/backend"
)

type StorageFactoryConfig struct {
	StorageConfig             storagebackend.Config
	Serializer                runtime.StorageSerializer
	APIResourceConfig         *storage.ResourceConfig
	DefaultResourceEncoding   *storage.DefaultResourceEncodingConfig
	ResourceEncodingOverrides []schema.GroupVersionResource
	DefaultStorageMediaType   string
}

func NewStorageFactoryConfig() *StorageFactoryConfig {
	return &StorageFactoryConfig{
		Serializer: defaultscheme.Codecs,
	}
}

func (c *StorageFactoryConfig) Complete(mongoOptions *options.MongoOptions) (*completedStorageFactoryConfig, error) {
	c.StorageConfig = mongoOptions.StorageConfig
	return &completedStorageFactoryConfig{c}, nil
}

type completedStorageFactoryConfig struct {
	*StorageFactoryConfig
}

func (c *completedStorageFactoryConfig) New() (*storage.DefaultStorageFactory, error) {
	resourceEncodingConfig := resourceconfig.MergeResourceEncodingConfigs(c.DefaultResourceEncoding, c.ResourceEncodingOverrides)
	storageFactory := storage.NewDefaultStorageFactory(
		c.StorageConfig,
		c.DefaultStorageMediaType,
		c.Serializer,
		resourceEncodingConfig,
		c.APIResourceConfig,
	)

	return storageFactory, nil
}
