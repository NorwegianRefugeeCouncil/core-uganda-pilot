package webapp

import (
	"github.com/nrc-no/core/internal/form"
	"github.com/nrc-no/core/internal/validation"
	"strings"
)

// controlValidationResponse describes the structure of validation responses to send to the client
type controlValidationResponse struct {
	Type   form.ControlType `json:"type" bson:"type"`
	Name   string           `json:"name" bson:"name"`
	Errors []string         `json:"errors" bson:"errors"`
}

type formValidation []controlValidationResponse

// makeFormValidation returns a slice of form.Control populated with validated template form elements.
func makeFormValidation(errors validation.ErrorList, f form.Form) formValidation {
	var result []controlValidationResponse
	errorMessagesFromName := make(map[string][]string)
	for _, err := range errors {
		names := strings.Split(err.Field, ".")
		name := names[len(names)-1]

		errorMessagesFromName[name] = append(errorMessagesFromName[name], err.Detail)
	}
	for name := range errorMessagesFromName {
		var cvr controlValidationResponse
		cvr.Name = name
		cvr.Errors = errorMessagesFromName[name]
		if control := f.Controls.FindByName(name); control != nil {
			cvr.Type = control.Type
		}
		result = append(result, cvr)
	}
	return result
}
