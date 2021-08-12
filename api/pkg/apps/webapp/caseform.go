package webapp

import (
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/validation"
)

type ValidatedCaseTemplate struct {
	Template Template `json:"template"`
}

type Template struct {
	FormElements []FormElement `json:"formElements"`
}

type FormElement struct {
	form.FormElement
}

func NewValidatedTemplate(template *cms.CaseTemplate, errors validation.ErrorList) *ValidatedCaseTemplate {
	formElements := []FormElement{}
	for _, element := range template.FormElements {
		if errs := errors.FindFamily(element.Attributes.Name); len(*errs) > 0 {
			element.Errors = errs
		}
		formElements = append(formElements, FormElement{
			FormElement: element,
		})
	}
	return &ValidatedCaseTemplate{Template{formElements}}
}
