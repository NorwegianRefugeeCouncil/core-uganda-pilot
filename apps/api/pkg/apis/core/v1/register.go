package v1

import (
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/runtime/serializer/json"
)

const GroupName = "core"

var (
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	Scheme             = runtime.NewScheme()
	Codecs             = json.NewSerializer(json.DefaultMetaFactory, Scheme, Scheme)
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
