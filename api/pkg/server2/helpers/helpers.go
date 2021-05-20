package helpers

import (
	"fmt"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
)

// GetSchemaForVersion returns the validation schema for the given version or nil.
func GetSchemaForVersion(crd *corev1.CustomResourceDefinition, version string) (corev1.CustomResourceDefinitionValidation, error) {
	for _, v := range crd.Spec.Versions {
		if version == v.Name {
			return v.Schema, nil
		}
	}
	return corev1.CustomResourceDefinitionValidation{}, fmt.Errorf("version %s not found in apiextensionsv1.CustomResourceDefinition: %v", version, crd.Name)
}
