package v1

import (
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
)

const GroupName = "meta"

var (
	schemeBuilder      runtime.SchemeBuilder
	localSchemeBuilder = &schemeBuilder
)

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
var scheme = runtime.NewScheme()

func AddToGroupVersion(scheme *runtime.Scheme, groupVersion schema.GroupVersion) {
	scheme.AddKnownTypes(groupVersion, &Status{})
}
