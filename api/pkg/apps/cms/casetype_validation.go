package cms

import "github.com/nrc-no/core/pkg/validation"

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
	if len(caseType.TeamID) == 0 {
		errList = append(errList, validation.Required(path.Child("teamId"), "team is required"))
	}
	if caseType.CaseTemplate == nil || caseType.CaseTemplate.FormElements == nil {
		errList = append(errList, validation.Required(path.Child("caseTemplate"), "template is required"))
	}

	return errList
}
