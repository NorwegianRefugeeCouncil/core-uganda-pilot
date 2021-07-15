package iam

import "github.com/nrc-no/core/pkg/validation"

func ValidatePartyType(partyType *PartyType, path *validation.Path) validation.ErrorList {
	errs := validation.ErrorList{}

	// Check name
	if len(partyType.Name) == 0 {
		err := validation.Required(path.Child("name"), "name is required")
		errs = append(errs, err)
	} else if !validation.IsValidAlphaNumeric(partyType.Name) {
		err := validation.Invalid(path.Child("name"), partyType.Name, validation.InvalidAlphaNumericDetail)
		errs = append(errs, err)
	}

	return errs
}
