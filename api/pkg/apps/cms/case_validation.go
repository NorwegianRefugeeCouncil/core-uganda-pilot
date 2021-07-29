package cms

import (
	"fmt"
	"github.com/nrc-no/core/pkg/validation"
	"strconv"
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
		if uuid != "" && !validation.IsValidUUID(uuid) {
			errList = append(errList, validation.Invalid(path.Child(name), uuid, fmt.Sprintf("%s is not a valid UUID", name)))
		}
	}

	// Validate form elements
	for _, elem := range kase.Template.FormElements {
		switch elem.Type {
		case Checkbox:
			for i, option := range elem.Attributes.CheckboxOptions {
				if option.Required && !contains(elem.Attributes.Value, strconv.Itoa(i)) {
					err := validation.Required(path.Child(elem.Attributes.Name).Index(i), fmt.Sprintf("%s is required", elem.Attributes.Name))
					errList = append(errList, err)
				}
			}
			fallthrough
		default:
			if elem.Validation.Required && len(elem.Attributes.Value) == 0 {
				err := validation.Required(path.Child(elem.Attributes.Name), fmt.Sprintf("%s is required", elem.Attributes.Name))
				errList = append(errList, err)
			}
			break
		}
		// TODO COR-156 implement validation for specific input controls (email, date, etc)
	}

	return errList
}

func contains(slice []string, elem string) bool {
	for _, s := range slice {
		if s == elem {
			return true
		}
	}
	return false
}
