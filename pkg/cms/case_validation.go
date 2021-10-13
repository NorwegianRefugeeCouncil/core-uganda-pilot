package cms

import (
	"fmt"
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/validation"
)

func ValidateCase(kase *Case, path *validation.Path) validation.ErrorList {
	result := validation.ErrorList{}

	result = validateUUIDs(kase, path, result)
	result = validateFormData(kase, path, result)

	return result
}

func validateUUIDs(kase *Case, path *validation.Path, errList validation.ErrorList) validation.ErrorList {
	uuids := map[string]string{
		"id":         kase.ID,
		"caseTypeId": kase.CaseTypeID,
		"partyId":    kase.PartyID,
		"parentId":   kase.ParentID,
		"teamId":     kase.TeamID,
		"creatorId":  kase.CreatorID}

	requiredUUIDs := []string{"id", "caseTypeId", "partyId", "teamId", "creatorId"}

	for name, uuid := range uuids {
		// check for required uuids
		for _, required := range requiredUUIDs {
			if name == required && len(uuid) == 0 {
				msg := fmt.Sprintf("%s was empty but is required", name)
				errList = append(errList, validation.Required(path.Child(name), msg))
				break
			}
		}
		// check that uuids are valid
		if len(uuid) > 0 && !validation.IsValidUUID(uuid) {
			errList = append(errList, validation.Invalid(path.Child(name), uuid, fmt.Sprintf("%s is not a valid UUID", name)))
		}
	}
	return errList
}

func validateFormData(kase *Case, path *validation.Path, errList validation.ErrorList) validation.ErrorList {
	// skip if no data
	if kase.FormData == nil {
		return errList
	}

	// range over form data and validate according to control
	for name, value := range kase.FormData {
		control := kase.Form.Controls.FindByName(name)
		if control != nil {
			errs := form.ValidateControlValue(*control, value, path)
			if len(errs) > 0 {
				errList = append(errList, errs...)
			}
		}
	}
	return errList
}
