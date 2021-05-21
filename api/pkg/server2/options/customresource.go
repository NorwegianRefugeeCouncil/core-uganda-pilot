package options

import (
	"github.com/nrc-no/core/api/pkg/server2/endpoints/handlers"
	genericregistry "github.com/nrc-no/core/api/pkg/server2/registry/generic"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// NewCRDRESTOptionsGetter create a RESTOptionsGetter for CustomResources.
func NewCRDRESTOptionsGetter(mongoOptions MongoOptions) genericregistry.RESTOptionsGetter {
	ret := handlers.CRDRESTOptionsGetter{
		StorageConfig: mongoOptions.StorageConfig,
		//StoragePrefix:           mongoOptions.StorageConfig.Prefix,
		//EnableWatchCache:        mongoOptions.EnableWatchCache,
		//DefaultWatchCacheSize:   mongoOptions.DefaultWatchCacheSize,
		//EnableGarbageCollection: mongoOptions.EnableGarbageCollection,
		//DeleteCollectionWorkers: mongoOptions.DeleteCollectionWorkers,
		//CountMetricPollPeriod:   mongoOptions.StorageConfig.CountMetricPollPeriod,
	}
	ret.StorageConfig.Codec = unstructured.UnstructuredJSONScheme
	return ret
}
