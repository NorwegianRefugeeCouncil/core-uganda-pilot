package v1

import (
	coremetav1 "github.com/nrc-no/core/api/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "core.nrc.no"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

var (
	SchemeBuilder      runtime.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)

func init() {
	localSchemeBuilder.Register(addKnownTypes)
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&FormDefinition{},
		&FormDefinitionList{},
		&CustomResourceDefinition{},
		&CustomResourceDefinitionList{},
		&OrganizationScope{},
		&OperatingScopeList{},
		&User{},
		&UserList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	coremetav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}
