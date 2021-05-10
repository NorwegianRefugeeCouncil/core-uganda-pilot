package storage

import (
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/storage/backend"
	"github.com/sirupsen/logrus"
)

type DefaultStorageFactory struct {
	// StorageConfig describes how to create a storage backend in general.
	// Its authentication information will be used for every storage.Interface returned.
	StorageConfig backend.Config

	// Overrides map[schema.GroupResource]groupResourceOverrides

	DefaultResourcePrefixes map[schema.GroupResource]string

	// DefaultMediaType is the media type used to store resources. If it is not set, "application/json" is used.
	DefaultMediaType string

	// DefaultSerializer is used to create encoders and decoders for the storage.Interface.
	DefaultSerializer runtime.StorageSerializer

	// ResourceEncodingConfig describes how to encode a particular GroupVersionResource
	ResourceEncodingConfig ResourceEncodingConfig

	// APIResourceConfigSource indicates whether the *storage* is enabled, NOT the API
	// This is discrete from resource enablement because those are separate concerns.  How this source is configured
	// is left to the caller.
	APIResourceConfigSource APIResourceConfigSource

	// newStorageCodecFn exists to be overwritten for unit testing.
	newStorageCodecFn func(opts StorageCodecConfig) (codec runtime.Codec, encodeVersioner runtime.GroupVersioner, err error)
}

func NewDefaultStorageFactory(
	config backend.Config,
) *DefaultStorageFactory {
	return &DefaultStorageFactory{
		StorageConfig: config,
	}
}

func (s *DefaultStorageFactory) NewConfig(groupResource schema.GroupResource) (*backend.Config, error) {
	chosenStorageResource := s.getStorageGroupResource(groupResource)

	storageConfig := s.StorageConfig
	codecConfig := StorageCodecConfig{
		StorageMediaType:  s.DefaultMediaType,
		StorageSerializer: s.DefaultSerializer,
	}

	var err error
	codecConfig.StorageVersion, err = s.ResourceEncodingConfig.StorageEncodingFor(chosenStorageResource)
	if err != nil {
		return nil, err
	}
	codecConfig.MemoryVersion, err = s.ResourceEncodingConfig.InMemoryEncodingFor(groupResource)
	if err != nil {
		return nil, err
	}
	codecConfig.Config = storageConfig

	storageConfig.Codec, storageConfig.EncodeVersioner, err = s.newStorageCodecFn(codecConfig)
	if err != nil {
		return nil, err
	}
	logrus.Infof("storing %v in %v, reading as %v from %#v", groupResource, codecConfig.StorageVersion, codecConfig.MemoryVersion, codecConfig.Config)

	return &storageConfig, nil
}

func (s *DefaultStorageFactory) getStorageGroupResource(groupResource schema.GroupResource) schema.GroupResource {
	//for _, potentialStorageResource := range s.Overrides[groupResource].cohabitatingResources {
	//  if s.APIResourceConfigSource.AnyVersionForGroupEnabled(potentialStorageResource.Group) {
	//    return potentialStorageResource
	//  }
	//}
	return groupResource
}
