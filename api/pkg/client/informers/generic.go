package informers

import (
	"fmt"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	discoveryv1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=core.nrc.no, Version=v1
	case corev1.SchemeGroupVersion.WithResource("customresourcedefinitions"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().CustomResourceDefinitions().Informer()}, nil
	case corev1.SchemeGroupVersion.WithResource("formdefinitions"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().FormDefinitions().Informer()}, nil
	case discoveryv1.SchemeGroupVersion.WithResource("apiservices"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Discovery().V1().APIServices().Informer()}, nil
	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
