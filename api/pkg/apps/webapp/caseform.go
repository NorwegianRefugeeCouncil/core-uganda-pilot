package webapp

import (
	"github.com/nrc-no/core/pkg/form"
)

type ValidatedCaseTemplate struct {
	Template Template `json:"template"`
}

type Template struct {
	FormElements []form.FormElement `json:"formElements"`
}
