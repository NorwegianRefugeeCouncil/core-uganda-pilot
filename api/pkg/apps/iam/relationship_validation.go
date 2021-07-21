package iam

import "github.com/nrc-no/core/pkg/validation"

func ValidateRelationship(relationship *Relationship, path *validation.Path) validation.ErrorList {
	errs := validation.ErrorList{}

	// Validate RelationshipType
	if len(relationship.RelationshipTypeID) == 0 {
		err := validation.Required(path.Child("relationshipTypeId"), "relationship type is required")
		errs = append(errs, err)
	}
	// Validate FirstParty
	if len(relationship.FirstPartyID) == 0 {
		err := validation.Required(path.Child("firstParty"), "first party is required")
		errs = append(errs, err)
	}
	// Validate SecondParty
	if len(relationship.SecondPartyID) == 0 {
		err := validation.Required(path.Child("secondParty"), "second party is required")
		errs = append(errs, err)
	}

	return errs
}
