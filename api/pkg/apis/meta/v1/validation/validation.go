package validation

import (
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
	"github.com/nrc-no/core/apps/api/pkg/util/validation"
	"github.com/nrc-no/core/apps/api/pkg/util/validation/field"
	"strings"
)

const (
	totalAnnotationSizeLimitB int = 256 * (1 << 10) // 256 kB
)

func ValidateLabels(labels map[string]string, fldPath *field.Path) exceptions.ErrorList {
	allErrs := exceptions.ErrorList{}
	for k, v := range labels {
		allErrs = append(allErrs, ValidateLabelName(k, fldPath)...)
		for _, msg := range validation.IsValidLabelValue(v) {
			allErrs = append(allErrs, exceptions.Invalid(fldPath, v, msg))
		}
	}
	return allErrs
}

// ValidateLabelName validates that the label name is correctly defined.
func ValidateLabelName(labelName string, fldPath *field.Path) exceptions.ErrorList {
	allErrs := exceptions.ErrorList{}
	for _, msg := range validation.IsQualifiedName(labelName) {
		allErrs = append(allErrs, exceptions.Invalid(fldPath, labelName, msg))
	}
	return allErrs
}

// ValidateAnnotations validates that a set of annotations are correctly defined.
func ValidateAnnotations(annotations map[string]string, fldPath *field.Path) exceptions.ErrorList {
	allErrs := exceptions.ErrorList{}
	var totalSize int64
	for k, v := range annotations {
		for _, msg := range validation.IsQualifiedName(strings.ToLower(k)) {
			allErrs = append(allErrs, exceptions.Invalid(fldPath, k, msg))
		}
		totalSize += (int64)(len(k)) + (int64)(len(v))
	}
	if totalSize > (int64)(totalAnnotationSizeLimitB) {
		allErrs = append(allErrs, exceptions.TooLong(fldPath, "", totalAnnotationSizeLimitB))
	}
	return allErrs
}
