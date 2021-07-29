package webapp

import (
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/validation"
)

type ValidatedCaseTemplate struct {
	Template struct {
		FormElements []FormElement
	}
}

type FormElement struct {
	cms.FormElement
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
	return &ValidatedCaseTemplate{struct{ FormElements []FormElement }{FormElements: formElements}}
}
