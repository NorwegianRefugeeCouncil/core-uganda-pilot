package internalinterfaces

import (
	"github.com/nrc-no/core/api/pkg/client/typed"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	"time"
)

// NewInformerFunc takes versioned.Interface and time.Duration to return a SharedIndexInformer.
type NewInformerFunc func(typed.Interface, time.Duration) cache.SharedIndexInformer

// SharedInformerFactory a small interface to allow for adding an informer without an import cycle
type SharedInformerFactory interface {
	Start(stopCh <-chan struct{})
	InformerFor(obj runtime.Object, newFunc NewInformerFunc) cache.SharedIndexInformer
}

// TweakListOptionsFunc is a function that transforms a v1.ListResourcesOptions.
type TweakListOptionsFunc func(*v1.ListOptions)
