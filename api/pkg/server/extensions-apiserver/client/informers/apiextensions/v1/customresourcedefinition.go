package v1

import (
  "context"
  apiextensionsv1 "github.com/nrc-no/core/apps/api/pkg/apis/apiextensions/v1"
  metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
  "github.com/nrc-no/core/apps/api/pkg/runtime"
  "github.com/nrc-no/core/apps/api/pkg/server/extensions-apiserver/client/clientset"
  "github.com/nrc-no/core/apps/api/pkg/server/extensions-apiserver/client/informers/internalinterfaces"
  v1 "github.com/nrc-no/core/apps/api/pkg/server/extensions-apiserver/client/listers/apiextensions/v1"
  "github.com/nrc-no/core/apps/api/pkg/tools/cache"
  "github.com/nrc-no/core/apps/api/pkg/watch"
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
func NewCustomResourceDefinitionInformer(client clientset.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
  return NewFilteredCustomResourceDefinitionInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredCustomResourceDefinitionInformer constructs a new informer for CustomResourceDefinition type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCustomResourceDefinitionInformer(client clientset.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
  return cache.NewSharedIndexInformer(
    &cache.ListWatch{
      ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
        if tweakListOptions != nil {
          tweakListOptions(&options)
        }
        return client.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), options)
      },
      WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
        if tweakListOptions != nil {
          tweakListOptions(&options)
        }
        return client.ApiextensionsV1().CustomResourceDefinitions().Watch(context.TODO(), options)
      },
    },
    &apiextensionsv1.CustomResourceDefinition{},
    resyncPeriod,
    indexers,
  )
}

func (f *customResourceDefinitionInformer) defaultInformer(client clientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
  return NewFilteredCustomResourceDefinitionInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *customResourceDefinitionInformer) Informer() cache.SharedIndexInformer {
  return f.factory.InformerFor(&apiextensionsv1.CustomResourceDefinition{}, f.defaultInformer)
}

func (f *customResourceDefinitionInformer) Lister() v1.CustomResourceDefinitionLister {
  return v1.NewCustomResourceDefinitionLister(f.Informer().GetIndexer())
}
