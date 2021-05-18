package install

import (
	"github.com/nrc-no/core/apps/api/pkg/apis/apiextensions"
	v1 "github.com/nrc-no/core/apps/api/pkg/apis/apiextensions/v1"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	utilruntime "github.com/nrc-no/core/apps/api/pkg/util/runtime"
)

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(apiextensions.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1.SchemeGroupVersion))
}
