package iam

import "github.com/nrc-no/core/pkg/validation"

func ValidateParty(party *Party, path *validation.Path) validation.ErrorList {
	errs := validation.ErrorList{}

	// Validate PartyTypeIDs
	if len(party.PartyTypeIDs) == 0 {
		err := validation.Required(path.Child("partyTypeIds"), "At least one party type is required")
		errs = append(errs, err)
	}

	return errs
}
