package webapp

import (
	"bytes"
	"fmt"
	"github.com/nrc-no/core/pkg/validation"
	"html"
	"html/template"
	"path"
	"regexp"
	"strings"
)

type FormRenderer struct {
	template *template.Template
}

func NewFormRenderer(templateDirectory string) *FormRenderer {
	f := &FormRenderer{}
	funcMap := template.FuncMap{
		"validationfeedback": validationFeedback,
	}
	fieldTemplates := template.Must(template.ParseGlob(path.Join(templateDirectory, "form/*.gohtml"))).Funcs(funcMap)
	f.template = fieldTemplates
	return f
}
func (f FormRenderer) Render(form Form) template.HTML {
	buf := bytes.Buffer{}
	for _, field := range form.Fields {
		err := f.template.ExecuteTemplate(&buf, "formfield", field)
		if err != nil {
			panic(err)
		}
	}
	return template.HTML(buf.String())
}

func (f FormRenderer) RenderFormField() template.HTML {
	return template.HTML("")
}

func (f FormRenderer) RenderValidationErrorsFor(field string) template.HTML {
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

func validationClass(errorList validation.ErrorList) string {
	if len(errorList) > 0 {
		return "is-invalid"
	}
	return "is-valid"
}

func validationFeedback(errorList validation.ErrorList, id string) template.HTML {
	if len(errorList) == 0 {
		return `<div class="valid-feedback">Looks good!</div>`
	}
	s := fmt.Sprintf(`<div id="%sFeedback" class="invalid-feedback">`, id)
	for i, e := range errorList {
		if i > 0 {
			s += `<br>`
		}
		s += html.EscapeString(e.Detail)
	}
	s += `</div>`
	return template.HTML(s)
}
