package form

import (
	"github.com/nrc-no/core/pkg/validation"
	"net/url"
)

import "github.com/nrc-no/core/pkg/validation"

type FieldType string

const (
	Text          FieldType = "text"
	Email         FieldType = "email"
	Phone         FieldType = "tel"
	URL           FieldType = "url"
	Date          FieldType = "date"
	Textarea      FieldType = "textarea"
	Dropdown      FieldType = "dropdown"
	Checkbox      FieldType = "checkbox"
	Radio         FieldType = "radio"
	TaxonomyInput FieldType = "taxonomyinput"
	File          FieldType = "file"
	CustomDiv     FieldType = "div"
)

type FormElement struct {
	Type       FieldType             `json:"type" bson:"type"`
	Attributes FormElementAttributes `json:"attributes" bson:"attributes"`
	Validation FormElementValidation `json:"validation" bson:"validation"`
	Errors     *validation.ErrorList `json:"errors"`
	Readonly   bool
}

type I18nString struct {
	Locale string `json:"locale" bson:"locale"`
	Value  string `json:"value" bson:"value"`
}

type I18nStringList []I18nString

// TODO COR-209 change static string fields (labels, descriptors) to []I18nString?
type FormElementAttributes struct {
	Label           string           `json:"label" bson:"label"`
	Name            string           `json:"name" bson:"name"`
	Value           []string         `json:"value" bson:"value"`
	Description     string           `json:"description" bson:"description"`
	Placeholder     string           `json:"placeholder" bson:"placeholder"`
	Multiple        bool             `json:"multiple" bson:"multiple"`
	Options         []string         `json:"options" bson:"options"`
	CheckboxOptions []CheckboxOption `json:"checkboxOptions" bson:"checkboxOptions"`
}

type CheckboxOption struct {
	Label    string `json:"label" bson:"label"`
	Required bool   `json:"required" bson:"required"`
}

type ControlValidation struct {
	Required bool `json:"required" bson:"required"`
}

type ValuedControl struct {
	*Control
	Value  []string
	Errors *validation.ErrorList
}

type ValuedForm struct {
	Controls []ValuedControl
	Errors   *validation.ErrorList
}

// Case templates
// https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema

// CaseTemplate contains a list of form elements used to construct a case form
type Form struct {
	// FormControls is an ordered list of the elements found in the form
	Controls []Control `json:"formcontrols" bson:"formcontrols"`
}

func (f *Form) FindControlByName(name string) *Control {
	for _, control := range f.Controls {
		if control.Name == name {
			return &control
		}
	}
	return nil
}

func NewValuedForm(form Form, values url.Values, errors validation.ErrorList) ValuedForm {
	var valuedControls []ValuedControl
	for _, control := range form.Controls {
		value := values[control.Name]
		errs := errors.FindFamily(control.Name)
		valuedControls = append(valuedControls, ValuedControl{
			Control: &control,
			Value:   value,
			Errors:  errs,
		})
	}
	errs := errors.Find("")
	return ValuedForm{
		Controls: valuedControls,
		Errors:   errs,
	}
}
