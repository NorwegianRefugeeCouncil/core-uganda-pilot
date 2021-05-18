package validation

import (
	"github.com/nrc-no/core/api/pkg/apis/core"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateFormDefinition(f *core.FormDefinition) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateFormDefinitionSpec(&f.Spec, field.NewPath("spec"))...)

	return allErrs
}

func ValidateFormDefinitionSpec(f *core.FormDefinitionSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(f.Group) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("group"), f.Group, "cannot be set without group"))
	}

	return allErrs
}
