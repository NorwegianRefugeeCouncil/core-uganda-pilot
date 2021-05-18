package validation

import (
	"fmt"
	"github.com/nrc-no/core/api/pkg/apis/core"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"regexp"
	"strconv"
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
	allErrs = append(allErrs, ValidateFormDefinitionElement(&f.Root, fldPath.Child("root"), 0)...)
	return allErrs
}

func ValidateFormDefinitionElement(f *core.FormElementDefinition, fldPath *field.Path, level int) field.ErrorList {
	allErrs := field.ErrorList{}

	isRoot := level == 0

	if isRoot && len(f.Key) > 0 {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("key"), f.Key, []string{""}))
	}

	if !isRoot && f.Type == core.SectionType && len(f.Key) != 0 {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("key"), f.Key, []string{""}))
	}

	if !isRoot && f.Type != core.SectionType && len(f.Key) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("key"), "key is required"))
	}

	if len(f.Children) > 0 && !core.CanHaveChildren(f.Type) {
		allErrs = append(allErrs, field.TooMany(fldPath.Child("children"), len(f.Children), 0))
	}

	if len(f.Children) == 0 && core.MustHaveChildren(f.Type) {
		allErrs = append(allErrs, field.Required(fldPath.Child("children"), "section elements must have at least 1 child"))
	}

	for i, child := range f.Children {
		allErrs = append(allErrs, ValidateFormDefinitionElement(&child, fldPath.Child("children").Index(i), level+1)...)
	}

	if f.Type != core.IntegerType {

		if len(f.Min) != 0 {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("min"), f.Min, []string{""}))
		}
		if len(f.Max) != 0 {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("max"), f.Max, []string{""}))
		}

	} else {

		var min int64
		var hasMin bool
		var err error
		if len(f.Min) != 0 {
			min, err = strconv.ParseInt(f.Min, 10, 64)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("min"), f.Min, "invalid number"))
			}
			hasMin = true
		}

		var max int64
		var hasMax bool
		if len(f.Max) != 0 {
			max, err = strconv.ParseInt(f.Max, 10, 64)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("max"), f.Max, "invalid number"))
			}
			hasMax = true
		}

		if hasMin && hasMax {
			if min > max {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("max"), f.Max, "minimum cannot be greater than maximum"))
				allErrs = append(allErrs, field.Invalid(fldPath.Child("min"), f.Min, "minimum cannot be greater than maximum"))
			}
		}
	}

	if f.Type == core.LongTextType || f.Type == core.ShortTextType {

		if f.MinLength < 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("minLength"), f.MinLength, "cannot have negative minimum length"))
		}

		if f.MaxLength != nil && *f.MaxLength < 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("minLength"), f.MinLength, "cannot have negative minimum length"))
		}

		if f.MinLength != 0 && f.MaxLength != nil && *f.MaxLength < f.MinLength {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("maxLength"), *f.MaxLength, "maximum length cannot be smaller than minimum length"))
			allErrs = append(allErrs, field.Invalid(fldPath.Child("minLength"), f.MinLength, "maximum length cannot be smaller than minimum length"))
		}

		if len(f.Pattern) != 0 {
			_, err := regexp.Compile(f.Pattern)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("pattern"), f.Pattern, fmt.Sprintf("invalid pattern: %s", err.Error())))
			}
		}

	} else {

		if f.MinLength != 0 {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("minLength"), f.MinLength, []string{""}))
		}
		if f.MaxLength != nil {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("maxLength"), f.MaxLength, []string{""}))
		}
		if len(f.Pattern) != 0 {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("pattern"), f.Pattern, []string{""}))
		}

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
