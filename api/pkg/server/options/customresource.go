package options

import (
	handlers2 "github.com/nrc-no/core/api/pkg/endpoints/handlers"
	"github.com/nrc-no/core/api/pkg/registry/generic"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// NewCRDRESTOptionsGetter create a RESTOptionsGetter for CustomResources.
func NewCRDRESTOptionsGetter(mongoOptions MongoOptions) generic.RESTOptionsGetter {
	ret := handlers2.CRDRESTOptionsGetter{
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
