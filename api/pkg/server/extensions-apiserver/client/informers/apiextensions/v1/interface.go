package v1

import "github.com/nrc-no/core/apps/api/pkg/server/extensions-apiserver/client/informers/internalinterfaces"

// Interface provides access to all the informers in this group version.
type Interface interface {
  // CustomResourceDefinitions returns a CustomResourceDefinitionInformer.
  CustomResourceDefinitions() CustomResourceDefinitionInformer
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

