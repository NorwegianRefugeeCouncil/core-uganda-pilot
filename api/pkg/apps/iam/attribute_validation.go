package iam

import (
	"fmt"
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/nrc-no/core/pkg/validation"
	"strconv"
)

func ValidateAttribute(attribute *Attribute, path *validation.Path) validation.ErrorList {
	errs := validation.ErrorList{}

	// Validate Name
	if len(attribute.Name) == 0 {
		err := validation.Required(path.Child("name"), "Name is required")
		errs = append(errs, err)
	} else if !validation.IsValidAlpha(attribute.Name) {
		err := validation.Invalid(path.Child("name"), attribute.Name, validation.InvalidAlphaDetail)
		errs = append(errs, err)
	}

	// Validate PartyTypeIDs
	if len(attribute.PartyTypeIDs) == 0 {
		err := validation.Required(path.Child("partyTypeIds"), "At least one party type is required")
		errs = append(errs, err)
	}

	// Validate Translations
	translationsPath := path.Child("translations")
	if len(attribute.Translations) == 0 {
		err := validation.Required(translationsPath, "At least one translation is required")
		errs = append(errs, err)
	} else {
		for i, translation := range attribute.Translations {
			translationPath := translationsPath.Index(i)
			errs = ValidateTranslation(translation, translationPath, errs)
		}
	}

	// Validate form elements
	errs = attribute.ValidateFormField(errs, path)

	return errs
}

func (a *Attribute) ValidateFormField(errList validation.ErrorList, path *validation.Path) validation.ErrorList {
	values := a.Attributes.Value
	fieldPath := path.Child(a.Attributes.Name)
	// no use validating empty values
	if !utils.AllEmpty(a.Attributes.Value) {
		switch a.Type {
		case form.Date:
			if !validation.IsValidDate(values[0]) {
				err := validation.Invalid(fieldPath, values[0], validation.InvalidDateDetail)
				errList = append(errList, err)
			}
		case form.Phone:
			if !validation.IsValidPhone(values[0]) {
				err := validation.Invalid(fieldPath, values[0], validation.InvalidPhoneDetail)
				errList = append(errList, err)
			}
		case form.Email:
			if !validation.IsValidEmail(values[0]) {
				err := validation.Invalid(fieldPath, values[0], validation.InvalidEmailDetail)
				errList = append(errList, err)
			}
		case form.Checkbox:
			for i, option := range a.Attributes.CheckboxOptions {
				if option.Required && !utils.Contains(a.Attributes.Value, strconv.Itoa(i)) {
					err := validation.Required(fieldPath.Index(i), fmt.Sprintf("%s is required", a.Attributes.Name))
					errList = append(errList, err)
				}
			}
		}
	}
	if a.Validation.Required && utils.AllEmpty(a.Attributes.Value) {
		err := validation.Required(fieldPath, fmt.Sprintf("%s is required", a.Attributes.Name))
		errList = append(errList, err)
	}
	return errList
}

func ValidateTranslation(translation AttributeTranslation, translationPath *validation.Path, errs validation.ErrorList) validation.ErrorList {
	// Validate Locale
	if len(translation.Locale) == 0 {
		err := validation.Required(translationPath.Child("locale"), "Locale is required")
		errs = append(errs, err)
	} else if !validation.IsValidAlpha(translation.Locale) {
		err := validation.Invalid(translationPath.Child("locale"), translation.Locale, validation.InvalidAlphaDetail)
		errs = append(errs, err)
	}
	// Validate LongFormulation
	if len(translation.LongFormulation) == 0 {
		err := validation.Required(translationPath.Child("long"), "Long formulation is required")
		errs = append(errs, err)
	}
	// Validate ShortFormulation
	if len(translation.ShortFormulation) == 0 {
		err := validation.Required(translationPath.Child("short"), "Short formulation is required")
		errs = append(errs, err)
	}
	return errs
}
