package server

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/server/options/encryptionconfig"
	"k8s.io/apiserver/pkg/server/storage"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"strings"
)

type StorageFactoryConfig struct {
	Serializer                   runtime.StorageSerializer
	APIResourceConfig            *storage.ResourceConfig
	DefaultResourceEncoding      *storage.DefaultResourceEncodingConfig
	StorageConfig                storagebackend.Config
	DefaultStorageMediaType      string
	EtcdServerOverrides          []string
	EncryptionProviderConfigPath string
}

func NewStorageFactoryConfig() *StorageFactoryConfig {
	return &StorageFactoryConfig{
		Serializer:              Codecs,
		DefaultResourceEncoding: storage.NewDefaultResourceEncodingConfig(Scheme),
	}
}

func (c *StorageFactoryConfig) Complete(etcdOptions *EtcdOptions) (*completedStorageFactoryConfig, error) {
	c.StorageConfig = etcdOptions.StorageConfig
	c.DefaultStorageMediaType = etcdOptions.DefaultStorageMediaType
	c.EtcdServerOverrides = etcdOptions.EtcdServersOverrides
	c.EncryptionProviderConfigPath = etcdOptions.EncryptionProviderConfigFilepath
	return &completedStorageFactoryConfig{c}, nil
}

type completedStorageFactoryConfig struct {
	*StorageFactoryConfig
}

func (c *completedStorageFactoryConfig) New() (*storage.DefaultStorageFactory, error) {
	storageFactory := storage.NewDefaultStorageFactory(
		c.StorageConfig,
		c.DefaultStorageMediaType,
		c.Serializer,
		c.DefaultResourceEncoding,
		c.APIResourceConfig,
		map[schema.GroupResource]string{},
	)

	for _, override := range c.EtcdServerOverrides {
		tokens := strings.Split(override, "#")
		apiresource := strings.Split(tokens[0], "/")

		group := apiresource[0]
		resource := apiresource[1]
		groupResource := schema.GroupResource{Group: group, Resource: resource}

		servers := strings.Split(tokens[1], ";")
		storageFactory.SetEtcdLocation(groupResource, servers)
	}
	if len(c.EncryptionProviderConfigPath) != 0 {
		transformerOverrides, err := encryptionconfig.GetTransformerOverrides(c.EncryptionProviderConfigPath)
		if err != nil {
			return nil, err
		}
		for groupResource, transformer := range transformerOverrides {
			storageFactory.SetTransformer(groupResource, transformer)
		}
	}
	return storageFactory, nil
}
