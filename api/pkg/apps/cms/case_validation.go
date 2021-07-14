package cms

import (
	"fmt"
	"github.com/nrc-no/core/pkg/validation"
)

func ValidateCase(kase *Case, path *validation.Path) validation.ErrorList {
	errList := validation.ErrorList{}

	// Validate form elements
	for _, elem := range kase.FormData.FormElements {
		if elem.Validation.Required && len(elem.Attributes.Value) == 0 {
			errList = append(errList, validation.Required(path.Child(elem.Attributes.ID), fmt.Sprintf("%s is required", elem.Attributes.ID)))
		}
		// TODO implement validation for specific input controls (email, date, etc)
	}

	return errList
}
