package form

import (
	"github.com/nrc-no/core/pkg/validation"
	"net/url"
)

// Form describes the content of an HTML form element as an ordered list of Control elements
type Form struct {
	Controls []Control `json:"controls" bson:"controls"`
	Groups   []Group   `json:"groups" bson:"groups"`
}

// ValuedForm is a Form which has undergone validation. It is composed of an ordered list of ValuedControl as well as
// possible validation Errors
type ValuedForm struct {
	Controls []ValuedControl
	Errors   *validation.ErrorList
}

// NewValuedForm takes a Form, url.Values corresponding to that form (originating from an HTTP form submission) and a
// validation.ErrorList containing 0 or more validation errors. NewValuedForm combines these three structures into
// a ValuedForm.
func NewValuedForm(form Form, values url.Values, errors validation.ErrorList) ValuedForm {
	var valuedControls []ValuedControl
	for _, control := range form.Controls {
		ctrl := ValuedControl{Control: control}
		if values != nil {
			value := values[control.Name]
			ctrl.Value = value
		}
		if errors != nil {
			errs := errors.FindFamily(control.Name)
			ctrl.Errors = errs
		}
		valuedControls = append(valuedControls, ctrl)
	}
	var result ValuedForm
	result.Controls = valuedControls
	if errors != nil {
		errs := errors.Find("")
		result.Errors = errs
	}
	return result
}
