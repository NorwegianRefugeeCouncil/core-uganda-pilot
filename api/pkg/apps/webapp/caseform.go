package webapp

import (
	"github.com/nrc-no/core/pkg/form"
)

type ValidatedCaseTemplate struct {
	Template Template `json:"template"`
}

type Template struct {
	Formcontrols []form.Control `json:"formcontrols"`
}

func (t *Template) Fill() *form.ValuedForm {
	// TODO implement me
	return nil
}
