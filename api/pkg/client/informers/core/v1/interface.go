package v1

import "github.com/nrc-no/core/api/pkg/client/informers/internalinterfaces"

// Interface provides access to all the informers in this group version.
type Interface interface {
	// CustomResourceDefinitions returns a CustomResourceDefinitionInformer.
	CustomResourceDefinitions() CustomResourceDefinitionInformer
	// FormDefinitions returns a FormDefinitionInformer.
	FormDefinitions() FormDefinitionInformer
	OperatingScopes() OperatingScopeInformer
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

// CustomResourceDefinitions returns a CustomResourceDefinitionInformer.
func (v *version) CustomResourceDefinitions() CustomResourceDefinitionInformer {
	return &customResourceDefinitionInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// FormDefinitions returns a FormDefinitionInformer.
func (v *version) FormDefinitions() FormDefinitionInformer {
	return &formDefinitionInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// OperatingScopes returns a FormDefinitionInformer.
func (v *version) OperatingScopes() OperatingScopeInformer {
	return &operatingScopeInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}
