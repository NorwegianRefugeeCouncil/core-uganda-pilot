package v1

import "github.com/nrc-no/core/api/pkg/client/informers/internalinterfaces"

// Interface provides access to all the informers in this group version.
type Interface interface {
	// APIServices returns a CustomResourceDefinitionInformer.
	APIServices() APIServiceInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// APIServices returns a CustomResourceDefinitionInformer.
func (v *version) APIServices() APIServiceInformer {
	return &apiServiceInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}
