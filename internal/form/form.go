package form

import (
	"github.com/nrc-no/core/internal/validation"
	"net/url"
)

// Form describes the content of an HTML form element as an ordered list of Control elements optionally grouped into
// sections
type Form struct {
	Controls Controls              `json:"controls" bson:"controls"`
	Sections []Section             `json:"sections" bson:"sections"`
	Errors   *validation.ErrorList `json:"errors" bson:"errors"`
}

// NewValidatedForm takes a Form, url.Values corresponding to that form (originating from an HTTP form submission) and a
// validation.ErrorList containing 0 or more validation errors and combines these three structures into a new Form.
func NewValidatedForm(phorm Form, values url.Values, errors validation.ErrorList) Form {
	var result Form
	var resultControls []Control
	for _, control := range phorm.Controls {
		ctrl := control
		if values != nil {
			value := values[control.Name]
			ctrl.Value = value
		}
		if errors != nil {
			errs := errors.FindFamily(control.Name)
			ctrl.Errors = errs
		}
		resultControls = append(resultControls, ctrl)
	}
	result.Controls = resultControls
	result.Sections = phorm.Sections
	if errors != nil {
		errs := errors.Find("")
		result.Errors = errs
	}
	return result
}
