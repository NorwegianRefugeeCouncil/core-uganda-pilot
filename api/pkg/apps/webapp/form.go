package webapp

import (
	"encoding/json"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/validation"
	"html/template"
	"net/url"
)

type CaseTypeForm struct {
	Form
}

type CaseForm struct {
	Form
}

type Form struct {
	WasValidated     bool
	ValidationErrors validation.ErrorList
	Fields           []FormField
	Template         *template.Template
}

type FormField struct {
	Type            FieldType
	Name            string
	Value           []string
	Class           string
	Options         []string
	CheckboxOptions []CheckboxOption
	Label           string
	Placeholder     string
	Description     string
	Errors          *validation.ErrorList
	Required        bool
	Multiple        bool
}

type CheckboxOption struct {
	Label    string
	Required bool
}

type FieldType cms.FieldType

func UnmarshalCaseTypeFormData(c *cms.CaseType, values url.Values) error {
	c.Name = values.Get("name")
	c.PartyTypeID = values.Get("partyTypeId")
	c.TeamID = values.Get("teamId")
	templateString := values.Get("template")
	if templateString == "" {
		c.Template = &cms.CaseTemplate{}
	} else {
		if err := json.Unmarshal([]byte(templateString), &c.Template); err != nil {
			return err
		}
	}
	return nil
}

func UnmarshalCaseFormData(c *cms.Case, caseTemplate *cms.CaseTemplate, values url.Values) error {
	c.CaseTypeID = values.Get("caseTypeId")
	c.PartyID = values.Get("partyId")
	c.Done = values.Get("done") == "on"
	c.ParentID = values.Get("parentId")
	c.TeamID = values.Get("teamId")
	var formElements []cms.FormElement
	for _, formElement := range caseTemplate.FormElements {
		formElement.Attributes.Value = values[formElement.Attributes.ID]
		formElements = append(formElements, formElement)
	}
	c.FormData = &cms.CaseTemplate{FormElements: formElements}
	return nil
}

func (f Form) FromCaseTemplate(tmpl *cms.CaseTemplate) {
	for _, element := range tmpl.FormElements {
		checkboxOptions := []CheckboxOption{}
		for _, option := range element.Attributes.CheckboxOptions {
			checkboxOptions = append(checkboxOptions, CheckboxOption{
				Label:    option.Label,
				Required: option.Required,
			})
		}
		formField := FormField{
			Type:            FieldType(element.Type),
			Name:            element.Attributes.ID,
			Value:           element.Attributes.Value,
			Options:         element.Attributes.Options,
			CheckboxOptions: checkboxOptions,
			Label:           element.Attributes.Label,
			Placeholder:     element.Attributes.Placeholder,
			Description:     element.Attributes.Description,
		}
		f.Fields = append(f.Fields, formField)
	}
}

func (f Form) IsValidated() bool {
	return f.WasValidated
}

func (f Form) GetValidationErrors() validation.ErrorList {
	return f.ValidationErrors
}
