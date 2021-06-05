package v1

import (
	"context"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	coremetav1 "github.com/nrc-no/core/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/api/pkg/client/informers/internalinterfaces"
	listers "github.com/nrc-no/core/api/pkg/client/listers/core/v1"
	"github.com/nrc-no/core/api/pkg/client/typed"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"time"
)

// OperatingScopeInformer provides access to a shared informer and lister for
// OperatingScopes.
type OperatingScopeInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() listers.OperatingScopeLister
}

type operatingScopeInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewOperatingScopeInformer constructs a new informer for OrganizationScope type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewOperatingScopeInformer(client typed.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredOperatingScopeInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredOperatingScopeInformer constructs a new informer for OrganizationScope type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredOperatingScopeInformer(client typed.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1().OperatingScopes().List(context.TODO(), options)
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
				return client.CoreV1().OperatingScopes().Watch(context.TODO(), opts)
			},
		},
		&corev1.OrganizationScope{},
		resyncPeriod,
		indexers,
	)
}

func (f *operatingScopeInformer) defaultInformer(client typed.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredOperatingScopeInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *operatingScopeInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&corev1.OrganizationScope{}, f.defaultInformer)
}

func (f *operatingScopeInformer) Lister() listers.OperatingScopeLister {
	return listers.NewOperatingScopeLister(f.Informer().GetIndexer())
}
