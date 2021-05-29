package v1

import (
	"context"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	v12 "github.com/nrc-no/core/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/api/pkg/client/informers/internalinterfaces"
	v1 "github.com/nrc-no/core/api/pkg/client/listers/core/v1"
	"github.com/nrc-no/core/api/pkg/client/typed"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"time"
)

// CustomResourceDefinitionInformer provides access to a shared informer and lister for
// CustomResourceDefinitions.
type CustomResourceDefinitionInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.CustomResourceDefinitionLister
}

type customResourceDefinitionInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewCustomResourceDefinitionInformer constructs a new informer for CustomResourceDefinition type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCustomResourceDefinitionInformer(client typed.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredCustomResourceDefinitionInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredCustomResourceDefinitionInformer constructs a new informer for CustomResourceDefinition type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCustomResourceDefinitionInformer(client typed.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1().CustomResourceDefinitions().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				opts := v12.ListResourcesOptions{
					Watch:               options.Watch,
					AllowWatchBookmarks: options.AllowWatchBookmarks,
					ResourceVersion:     options.ResourceVersion,
					TimeoutSeconds:      options.TimeoutSeconds,
					Limit:               &options.Limit,
					Continue:            options.Continue,
				}
				return client.CoreV1().CustomResourceDefinitions().Watch(context.TODO(), opts)
			},
		},
		&corev1.CustomResourceDefinition{},
		resyncPeriod,
		indexers,
	)
}

func (f *customResourceDefinitionInformer) defaultInformer(client typed.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredCustomResourceDefinitionInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *customResourceDefinitionInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&corev1.CustomResourceDefinition{}, f.defaultInformer)
}

func (f *customResourceDefinitionInformer) Lister() v1.CustomResourceDefinitionLister {
	return v1.NewCustomResourceDefinitionLister(f.Informer().GetIndexer())
}
