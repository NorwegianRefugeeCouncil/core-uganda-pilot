package webapp

import (
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/validation"
)

type ValidatedTemplate struct {
	FormElements []FormElement
}

type FormElement struct {
	cms.FormElement
	Errors validation.ErrorList
}

func NewValidatedTemplate(template *cms.CaseTemplate, errors *validation.ErrorList) *ValidatedTemplate {
	formFields := []FormElement{}
	for _, element := range template.FormElements {
		errs := errors.Find(element.Attributes.ID)
		formFields = append(formFields, FormElement{
			FormElement: element,
			Errors:      errs,
		})
	}
	return &ValidatedTemplate{formFields}
}
