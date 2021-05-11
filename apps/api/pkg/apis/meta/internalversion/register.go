package internalversion

import (
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
)

const GroupName = "meta"

var (
	SchemeBuilder      runtime.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func addToGroupVersion(scheme *runtime.Scheme) error {

	if err := scheme.AddIgnoredConversionType(&metav1.TypeMeta{}, &metav1.TypeMeta{}); err != nil {
		return err
	}

	scheme.AddKnownTypes(SchemeGroupVersion,
		&ListOptions{},
		&metav1.GetOptions{},
		&metav1.DeleteOptions{},
		&metav1.CreateOptions{},
		&metav1.UpdateOptions{},
	)

	scheme.AddUnversionedTypes(SchemeGroupVersion,
		&metav1.DeleteOptions{},
		&metav1.CreateOptions{},
		&metav1.UpdateOptions{},
	)

	metav1.AddToGroupVersion(scheme, metav1.SchemeGroupVersion)

	return nil
}

func init() {
	localSchemeBuilder.Register(addToGroupVersion)
}
