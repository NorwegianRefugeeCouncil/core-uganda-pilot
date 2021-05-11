package v1

import (
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/runtime/serializer"
)

// GroupName is the group name use in this package
const GroupName = "core"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

var (
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&FormDefinition{},
		&FormDefinitionList{},
	)
	scheme.AddKnownTypes(SchemeGroupVersion, &metav1.Status{})
	return nil
}

func init() {
	metav1.AddToGroupVersion(Scheme, metav1.SchemeGroupVersion)
	if err := AddToScheme(Scheme); err != nil {
		panic(err)
	}
}
