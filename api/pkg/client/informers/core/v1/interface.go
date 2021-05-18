package v1

import "github.com/nrc-no/core/apps/api/pkg/client/informers/internalinterfaces"

type Interface interface {
	FormDefinitions() FormDefinitionInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{
		factory:          f,
		namespace:        namespace,
		tweakListOptions: tweakListOptions,
	}
}

func (v *version) FormDefinitions() FormDefinitionInformer {
	return &formDefinitionInformer{
		factory:          v.factory,
		tweakListOptions: v.tweakListOptions,
		namespace:        v.namespace,
	}
}
