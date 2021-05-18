package v1

import (
	"context"
	core "github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/client/informers/internalinterfaces"
	v1 "github.com/nrc-no/core/apps/api/pkg/client/listers/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/client/nrc"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/tools/cache"
	"github.com/nrc-no/core/apps/api/pkg/watch"
	"time"
)

type FormDefinitionInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.FormDefinitionLister
}

type formDefinitionInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

func NewFormDefinitionsInformer(client nrc.Interface, namespace string, resyncPeriod time.Duration, indexer cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredFormDefinitionsInformer(client, namespace, resyncPeriod, indexer, nil)
}

func NewFilteredFormDefinitionsInformer(client nrc.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1().FormDefinitions().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1().FormDefinitions().Watch(context.TODO(), options)
			},
		},
		&core.FormDefinition{},
		resyncPeriod,
		indexers,
	)
}

func (f *formDefinitionInformer) defaultInformer(client nrc.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredFormDefinitionsInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *formDefinitionInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&core.FormDefinition{}, f.defaultInformer)
}

func (f *formDefinitionInformer) Lister() v1.FormDefinitionLister {
	return v1.NewFormDefinitionLister(f.Informer().GetIndexer())
}
