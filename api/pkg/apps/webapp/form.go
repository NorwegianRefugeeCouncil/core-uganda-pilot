package webapp

import (
	"bytes"
	"github.com/nrc-no/core/pkg/validation"
	"html/template"
	"strings"
)

type CaseForm struct {
	Form
}

type Form struct {
	WasValidated     bool
	ValidationErrors validation.ErrorList
	Fields           []FormField
}

type FormField struct {
	Type  FieldType
	Name  string
	Value []string
}

type FieldType string

const (
	Textarea  FieldType = "textarea"
	Textinput FieldType = "textinput"
	Dropdown  FieldType = "dropdown"
	Checkbox  FieldType = "checkbox"
	Email     FieldType = "email"
	Date      FieldType = "date"
	File      FieldType = "file"
	Time      FieldType = "time"
)

func (f Form) IsValidated() bool {
	return f.WasValidated
}

func (f Form) GetValidationErrors() validation.ErrorList {
	return f.ValidationErrors
}

func (f Form) RenderValidationError(field string) template.HTML {
	errs := f.ValidationErrors.Find(field)
	if len(errs) == 0 {
		return template.HTML("")
	}
	templateText := `
	<div id="{{index . "name"}}Feedback" class="invalid-feedback">
		{{- range (index . "errors") -}}
			<p>{{.Detail}}</p>
		{{- end -}}
	</div>
	`
	tmpl, _ := template.New("feedback").Parse(templateText)
	data := map[string]interface{}{
		"name":   field,
		"errors": errs,
	}
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	s := strings.TrimSpace(buf.String())
	return template.HTML(s)
}
