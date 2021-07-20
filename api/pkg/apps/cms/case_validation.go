package cms

import (
	"fmt"
	"github.com/nrc-no/core/pkg/validation"
)

func ValidateCase(kase *Case, path *validation.Path) validation.ErrorList {
	errList := validation.ErrorList{}

	// Validate UUIDs
	uuids := map[string]string{
		"caseTypeId": kase.CaseTypeID,
		"partyId":    kase.PartyID,
		"parentId":   kase.ParentID,
		"teamId":     kase.TeamID,
		"creatorId":  kase.CreatorID}

	for name, uuid := range uuids {
		if !validation.IsValidUUID(uuid) {
			errList = append(errList, validation.Invalid(path.Child(name), uuid, fmt.Sprintf("%s is not a valid UUID", name)))
		}
	}

	// Validate form elements
	for _, elem := range kase.FormData.FormElements {
		if elem.Validation.Required && len(elem.Attributes.Value) == 0 {
			errList = append(errList, validation.Required(path.Child(elem.Attributes.ID), fmt.Sprintf("%s is required", elem.Attributes.ID)))
		}
		// TODO COR-156 implement validation for specific input controls (email, date, etc)
	}

	return errList
}
