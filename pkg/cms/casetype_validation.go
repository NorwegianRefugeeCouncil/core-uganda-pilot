package cms

import "github.com/nrc-no/core/internal/validation"

func ValidateCaseType(caseType *CaseType, path *validation.Path) validation.ErrorList {
	errList := validation.ErrorList{}

	if len(caseType.Name) == 0 {
		errList = append(errList, validation.Required(path.Child("name"), "name is required"))
	} else if !validation.IsValidAlpha(caseType.Name) {
		errList = append(errList, validation.Invalid(path.Child("name"), caseType.Name, "name should only contain letters and spaces"))
	}
	if len(caseType.PartyTypeID) == 0 {
		errList = append(errList, validation.Required(path.Child("partyTypeId"), "party type is required"))
	}
	// todo validate that PartyTypeIds are valid uuids
	if len(caseType.TeamID) == 0 {
		errList = append(errList, validation.Required(path.Child("teamId"), "team is required"))
	}
	// todo validate that TeamID is a valid uuid
	// todo form.ValidateForm()
	// the form package should have a ValidateForm(form form.Form, path *validation.Path) validation.ErrorList method

	return errList
}
