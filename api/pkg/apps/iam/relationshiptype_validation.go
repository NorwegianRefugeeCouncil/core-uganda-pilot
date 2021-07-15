package iam

import "github.com/nrc-no/core/pkg/validation"

func ValidateRelationshipType(relationshipType *RelationshipType, path *validation.Path) validation.ErrorList {
	errs := validation.ErrorList{}

	// Validate Name
	if len(relationshipType.Name) == 0 {
		err := validation.Required(path.Child("name"), "name is required")
		errs = append(errs, err)
	} else if !validation.IsValidAlphaNumeric(relationshipType.Name) {
		err := validation.Invalid(path.Child("name"), relationshipType.Name, validation.InvalidAlphaNumericDetail)
		errs = append(errs, err)
	}

	// Validate FirstPartyRole
	if len(relationshipType.FirstPartyRole) == 0 {
		err := validation.Required(path.Child("firstPartyRole"), "First party role is required")
		errs = append(errs, err)
	} else if !validation.IsValidAlphaNumeric(relationshipType.FirstPartyRole) {
		err := validation.Invalid(path.Child("firstPartyRole"), relationshipType.FirstPartyRole, validation.InvalidAlphaNumericDetail)
		errs = append(errs, err)
	}

	// Validate SecondPartyRole
	if relationshipType.IsDirectional && len(relationshipType.SecondPartyRole) == 0 {
		err := validation.Required(path.Child("secondPartyRole"), "Second party role is required")
		errs = append(errs, err)
	} else if !validation.IsValidAlphaNumeric(relationshipType.SecondPartyRole) {
		err := validation.Invalid(path.Child("secondPartyRole"), relationshipType.SecondPartyRole, validation.InvalidAlphaNumericDetail)
		errs = append(errs, err)
	}

	// Validate Rules
	rulesPath := path.Child("rules")
	if len(relationshipType.Rules) == 0 {
		err := validation.Required(rulesPath, "At least one rule is required")
		errs = append(errs, err)
	} else {
		for i, rule := range relationshipType.Rules {
			p := rulesPath.Index(i)
			if len(rule.PartyTypeRule.FirstPartyTypeID) == 0 {
				err := validation.Required(p.Child("firstPartyTypeId"), "First party type is required")
				errs = append(errs, err)
			}
			if relationshipType.IsDirectional {
				if len(rule.PartyTypeRule.SecondPartyTypeID) == 0 {
					err := validation.Required(p.Child("secondPartyTypeId"), "Second party type is required")
					errs = append(errs, err)
				}
			}
		}
	}

	return errs
}
