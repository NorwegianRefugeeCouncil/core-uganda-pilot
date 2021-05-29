package v1

import (
	"context"
	discoveryv1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	coremetav1 "github.com/nrc-no/core/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/api/pkg/client/informers/internalinterfaces"
	listers "github.com/nrc-no/core/api/pkg/client/listers/discovery/v1"
	"github.com/nrc-no/core/api/pkg/client/typed"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"time"
)

// APIServiceInformer provides access to a shared informer and lister for
// APIServices.
type APIServiceInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() listers.APIServiceLister
}

type apiServiceInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewAPIServiceInformer constructs a new informer for APIService type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewAPIServiceInformer(client typed.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredAPIServiceInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredAPIServiceInformer constructs a new informer for APIService type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredAPIServiceInformer(client typed.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.DiscoveryV1().APIServices().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				opts := coremetav1.ListResourcesOptions{
					Watch:               options.Watch,
					AllowWatchBookmarks: options.AllowWatchBookmarks,
					ResourceVersion:     options.ResourceVersion,
					TimeoutSeconds:      options.TimeoutSeconds,
					Limit:               &options.Limit,
					Continue:            options.Continue,
				}
				return client.DiscoveryV1().APIServices().Watch(context.TODO(), opts)
			},
		},
		&discoveryv1.APIService{},
		resyncPeriod,
		indexers,
	)
}

func (f *apiServiceInformer) defaultInformer(client typed.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredAPIServiceInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *apiServiceInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&discoveryv1.APIService{}, f.defaultInformer)
}

func (f *apiServiceInformer) Lister() listers.APIServiceLister {
	return listers.NewAPIServiceLister(f.Informer().GetIndexer())
}
