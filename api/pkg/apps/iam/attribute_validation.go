package iam

import "github.com/nrc-no/core/pkg/validation"

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
	attribute.Validate(errs, path)

	return errs
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
