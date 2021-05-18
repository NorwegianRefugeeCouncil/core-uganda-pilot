package core

import (
	v12 "github.com/nrc-no/core/apps/api/pkg/client/informers/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/client/informers/internalinterfaces"
)

type Interface interface {
	V1() v12.Interface
}

type group struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &group{
		factory:          f,
		namespace:        namespace,
		tweakListOptions: tweakListOptions,
	}
}

func (g *group) V1() v12.Interface {
	return v12.New(g.factory, g.namespace, g.tweakListOptions)
}
