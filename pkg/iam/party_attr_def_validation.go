package iam

import (
	"github.com/nrc-no/core/pkg/validation"
)

func ValidatePartyAttributeDefinition(attribute *PartyAttributeDefinition, path *validation.Path) validation.ErrorList {
	result := validation.ErrorList{}

	result = validateUUID(path.Child("id"), attribute.ID, result)
	result = validateUUID(path.Child("countryId"), attribute.CountryID, result)
	if len(attribute.PartyTypeIDs) == 0 {
		err := validation.Required(path.Child("partyTypeIds"), "at least one party type is required")
		result = append(result, err)
	} else {
		for i, id := range attribute.PartyTypeIDs {
			result = validateUUID(path.Child("partyTypeIds").Index(i), id, result)
		}
	}
	return result
}

func validateUUID(path *validation.Path, uuid string, result validation.ErrorList) validation.ErrorList {
	var err *validation.Error
	if len(uuid) == 0 {
		err = validation.Required(path, "field is required")
	} else if !validation.IsValidUUID(uuid) {
		err = validation.Invalid(path, uuid, validation.InvalidUUIDDetail)
	}
	if err != nil {
		result = append(result, err)
	}
	return result
}
