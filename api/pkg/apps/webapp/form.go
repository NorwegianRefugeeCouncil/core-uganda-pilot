package webapp

import (
	"bytes"
	"encoding/json"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/validation"
	"html/template"
	"net/url"
	"regexp"
	"strings"
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
	Type        FieldType
	Name        string
	Value       []string
	Label       string
	Placeholder string
	Description string
	Error       *validation.Error
}

type FieldType string

const (
	Textarea  FieldType = "textarea"
	TextInput FieldType = "textinput"
	Dropdown  FieldType = "dropdown"
	Checkbox  FieldType = "checkbox"
	Email     FieldType = "email"
	Date      FieldType = "date"
	File      FieldType = "file"
	Time      FieldType = "time"
)

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

func (f Form) RenderFormField() template.HTML {
	return template.HTML("")
}

func (f Form) IsValidated() bool {
	return f.WasValidated
}

func (f Form) GetValidationErrors() validation.ErrorList {
	return f.ValidationErrors
}

func (f Form) RenderValidationError(err validation.Error) template.HTML {
	field := err.Field
	// get only the "name" part of the field (in case we pass in a path like 'parent[1].child'
	rg, _ := regexp.Compile(`^*([\w-])$`)
	matches := rg.FindStringSubmatch(field)
	if len(matches) > 1 {
		field = matches[1]
	}
	data := map[string]interface{}{
		"field":  field,
		"errors": errs,
	}
	templateText := `
	<div id="{{index . "field"}}Feedback" class="invalid-feedback">
		{{- range (index . "errors") -}}
			<p>{{.Detail}}</p>
		{{- end -}}
	</div>
	`
	return parseAndExecuteTemplate(templateText, data)
}

func (f Form) RenderValidationErrorsFor(field string) template.HTML {
	errs := f.ValidationErrors.Find(field)
	if len(errs) == 0 {
		return template.HTML("")
	}
	data := map[string]interface{}{
		"field":  field,
		"errors": errs,
	}
	templateText := `
	<div id="{{index . "field"}}Feedback" class="invalid-feedback">
		{{- range (index . "errors") -}}
			<p>{{.Detail}}</p>
		{{- end -}}
	</div>
	`
	return parseAndExecuteTemplate(templateText, data)
}

func getFieldTemplateForType(t FieldType) template.Template {
	switch t {
	case Textarea:
	case TextInput:
	case Dropdown:
	case Checkbox:
	}
}

func parseAndExecuteTemplate(templateText string, data map[string]interface{}) template.HTML {
	tmpl, _ := template.New("").Parse(templateText)
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	s := strings.TrimSpace(buf.String())
	return template.HTML(s)
}

func textareaTemplate(name string) *template.Template {
	t := `
	<div class="form-floating mb-3">
            <textarea id="{{.Field}}"
                      name="{{.Field}}"
                      class="form-control {{***TODO***}}"
                      style="height: 150px"
                      data-testid="form"
                      placeholder="{{.Placeholder}}"
                      {{if .Validation.Required}}required{{end}}>{{if .Attributes.Value}}{{index .Attributes.Value 0}}{{end}}</textarea>
            <label for="{{.Attributes.ID}}">{{.Attributes.Label}}</label>
            <div class="form-text">{{.Attributes.Description}}</div>
	</div>
	`
	tmpl, _ := template.New(name).Parse(t)
	return tmpl
}
