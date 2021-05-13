package internalinterfaces

import (
  v1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
  "github.com/nrc-no/core/apps/api/pkg/runtime"
  "github.com/nrc-no/core/apps/api/pkg/server/extensions-apiserver/client/clientset"
  "github.com/nrc-no/core/apps/api/pkg/tools/cache"
  "time"
)

// NewInformerFunc takes clientset.Interface and time.Duration to return a SharedIndexInformer.
type NewInformerFunc func(clientset.Interface, time.Duration) cache.SharedIndexInformer

// SharedInformerFactory a small interface to allow for adding an informer without an import cycle
type SharedInformerFactory interface {
  Start(stopCh <-chan struct{})
  InformerFor(obj runtime.Object, newFunc NewInformerFunc) cache.SharedIndexInformer
}

// TweakListOptionsFunc is a function that transforms a v1.ListOptions.
type TweakListOptionsFunc func(*v1.ListOptions)
