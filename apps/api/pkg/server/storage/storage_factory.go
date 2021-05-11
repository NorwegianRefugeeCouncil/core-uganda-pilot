package storage

import (
	"crypto/tls"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	storagebackend "github.com/nrc-no/core/apps/api/pkg/storage/backend"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
	"strings"
)

// Backend describes the storage servers, the information here should be enough
// for health validations.
type Backend struct {
	Server    string
	TLSConfig *tls.Config
}

// StorageFactory is the interface to locate the storage for a given GroupResource
type StorageFactory interface {
	// New finds the storage destination for the given group and resource. It will
	// return an error if the group has no storage destination configured.
	NewConfig(groupResource schema.GroupResource) (*storagebackend.Config, error)

	// ResourcePrefix returns the overridden resource prefix for the GroupResource
	// This allows for cohabitation of resources with different native types and provides
	// centralized control over the shape of etcd directories
	ResourcePrefix(groupResource schema.GroupResource) string

	// Backends gets all backends for all registered storage destinations.
	// Used for getting all instances for health validations.
	Backends() []Backend
}

type DefaultStorageFactory struct {
	// StorageConfig describes how to create a storage backend in general.
	// Its authentication information will be used for every storage.Interface returned.
	StorageConfig storagebackend.Config

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
	config storagebackend.Config,
	defaultMediaType string,
	defaultSerializer runtime.StorageSerializer,
	resourceEncodingCOnfig ResourceEncodingConfig,
	resourceConfig APIResourceConfigSource,
) *DefaultStorageFactory {
	return &DefaultStorageFactory{
		StorageConfig:           config,
		DefaultMediaType:        defaultMediaType,
		DefaultSerializer:       defaultSerializer,
		ResourceEncodingConfig:  resourceEncodingCOnfig,
		APIResourceConfigSource: resourceConfig,
		DefaultResourcePrefixes: map[schema.GroupResource]string{},
		newStorageCodecFn:       NewStorageCodec,
	}
}

func (s *DefaultStorageFactory) NewConfig(groupResource schema.GroupResource) (*storagebackend.Config, error) {
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

	storageConfig.Prefix = storageConfig.Prefix + "/" + groupResource.Group + "/" + groupResource.Resource

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

// Backends returns all backends for all registered storage destinations.
// Used for getting all instances for health validations.
func (s *DefaultStorageFactory) Backends() []Backend {
	servers := sets.NewString(s.StorageConfig.Transport.ServerList...)
	//
	//for _, overrides := range s.Overrides {
	//  servers.Insert(overrides.etcdLocation...)
	//}
	//
	//tlsConfig := &tls.Config{
	//  InsecureSkipVerify: true,
	//}
	//if len(s.StorageConfig.Transport.CertFile) > 0 && len(s.StorageConfig.Transport.KeyFile) > 0 {
	//  cert, err := tls.LoadX509KeyPair(s.StorageConfig.Transport.CertFile, s.StorageConfig.Transport.KeyFile)
	//  if err != nil {
	//    klog.Errorf("failed to load key pair while getting backends: %s", err)
	//  } else {
	//    tlsConfig.Certificates = []tls.Certificate{cert}
	//  }
	//}
	//if len(s.StorageConfig.Transport.TrustedCAFile) > 0 {
	//  if caCert, err := ioutil.ReadFile(s.StorageConfig.Transport.TrustedCAFile); err != nil {
	//    klog.Errorf("failed to read ca file while getting backends: %s", err)
	//  } else {
	//    caPool := x509.NewCertPool()
	//    caPool.AppendCertsFromPEM(caCert)
	//    tlsConfig.RootCAs = caPool
	//    tlsConfig.InsecureSkipVerify = false
	//  }
	//}

	backends := []Backend{}
	for server := range servers {
		backends = append(backends, Backend{
			Server: server,
			// We can't share TLSConfig across different backends to avoid races.
			// For more details see: http://pr.k8s.io/59338
			// TLSConfig: tlsConfig.Clone(),
		})
	}
	return backends
}

func (s *DefaultStorageFactory) ResourcePrefix(groupResource schema.GroupResource) string {
	chosenStorageResource := s.getStorageGroupResource(groupResource)
	//groupOverride := s.Overrides[getAllResourcesAlias(chosenStorageResource)]
	//exactResourceOverride := s.Overrides[chosenStorageResource]

	mongoResourcePrefix := s.DefaultResourcePrefixes[chosenStorageResource]
	//if len(groupOverride.etcdResourcePrefix) > 0 {
	//  etcdResourcePrefix = groupOverride.etcdResourcePrefix
	//}
	//if len(exactResourceOverride.etcdResourcePrefix) > 0 {
	//  etcdResourcePrefix = exactResourceOverride.etcdResourcePrefix
	//}
	if len(mongoResourcePrefix) == 0 {
		mongoResourcePrefix = strings.ToLower(chosenStorageResource.Resource)
	}

	return mongoResourcePrefix
}
