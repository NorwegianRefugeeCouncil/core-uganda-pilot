package install

import (
	"github.com/nrc-no/core/api/pkg/apis/discovery"
	discoveryv1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func Install(scheme *runtime.Scheme) {
	utilruntime.Must(discovery.AddToScheme(scheme))
	utilruntime.Must(discoveryv1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(discoveryv1.SchemeGroupVersion))
}
