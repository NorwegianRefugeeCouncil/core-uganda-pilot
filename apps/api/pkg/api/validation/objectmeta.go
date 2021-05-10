package validation

import (
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	metav1validation "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1/validation"
	"github.com/nrc-no/core/apps/api/pkg/util/validation/field"
)

func ValidateObjectMetaAccessor(meta metav1.Object, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, metav1validation.ValidateLabels(meta.GetLabels(), fieldPath.Child("labels"))...)
	allErrs = append(allErrs, metav1validation.ValidateAnnotations(meta.GetAnnotations(), fieldPath.Child("annotations"))...)
	return allErrs
}
