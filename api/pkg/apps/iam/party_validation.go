package iam

import "github.com/nrc-no/core/pkg/validation"

func ValidateParty(party *Party, path *validation.Path) validation.ErrorList {
	errs := validation.ErrorList{}

	// Validate PartyTypeIDs
	if len(party.PartyTypeIDs) == 0 {
		err := validation.Required(path.Child("partyTypeIds"), "At least one party type is required")
		errs = append(errs, err)
	}

	// Validate Attributes
	//	// For now Party Attributes are just a map[string][]string but eventually
	//	// we will want TODO some changes to have more information about value types
	//	// and validation
	//attributesPath := path.Child("attributes")
	//for name, values := range party.Attributes {
	//	attributePath := attributesPath.Key(name)
	//}

	return errs
}
