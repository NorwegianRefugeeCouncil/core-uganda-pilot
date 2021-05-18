package validation

import (
	"github.com/nrc-no/core/api/pkg/apis/core"
	"k8s.io/apimachinery/pkg/util/sets"
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
		allErrs = append(allErrs, field.Required(fldPath.Child("group"), "group is required"))
	}
	allErrs = append(allErrs, ValidateFormDefinitionNames(&f.Names, fldPath.Child("names"))...)
	allErrs = append(allErrs, ValidateFormDefinitionVersions(f.Versions, fldPath.Child("versions"))...)
	return allErrs
}

func ValidateFormDefinitionVersions(f []core.FormDefinitionVersion, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(f) == 0 {
		allErrs = append(allErrs, field.Required(fldPath, "versions cannot be empty"))
	}

	seenVersionNames := sets.NewString()
	for i, version := range f {
		versionField := fldPath.Index(i)
		if seenVersionNames.Has(version.Name) {
			allErrs = append(allErrs, field.Duplicate(versionField.Child("name"), version.Name))
		}
		seenVersionNames.Insert(version.Name)
		allErrs = append(allErrs, ValidateFormDefinitionVersion(&version, versionField)...)
	}

	return allErrs
}

func ValidateFormDefinitionVersion(f *core.FormDefinitionVersion, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(f.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), "version name cannot be empty"))
	}
	allErrs = append(allErrs, ValidateFormDefinitionValidation(&f.Schema, fldPath.Child("schema"))...)
	return allErrs
}

func ValidateFormDefinitionValidation(f *core.FormDefinitionValidation, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, ValidateFormDefinitionSchema(&f.FormSchema, fldPath.Child("formSchema"))...)
	return allErrs
}

func ValidateFormDefinitionSchema(f *core.FormDefinitionSchema, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, ValidateFormDefinitionElement(&f.Root, fldPath.Child("root"))...)
	return allErrs
}

func ValidateFormDefinitionElement(f *core.FormElementDefinition, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if f.Type == core.SectionType {
		if len(f.Children) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("children"), "section elements must have at least 1 child"))
		}
	} else {
		if len(f.Children) > 0 {
			allErrs = append(allErrs, field.TooMany(fldPath.Child("children"), len(f.Children), 0))
		}
	}

	for i, child := range f.Children {
		allErrs = append(allErrs, ValidateFormDefinitionElement(&child, fldPath.Child("children").Index(i))...)
	}
	return allErrs
}

func ValidateFormDefinitionNames(f *core.FormDefinitionNames, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(f.Plural) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("plural"), "plural is required"))
	}
	if len(f.Singular) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("singular"), "singular is required"))
	}
	if len(f.Kind) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("kind"), "kind is required"))
	}
	return allErrs
}
