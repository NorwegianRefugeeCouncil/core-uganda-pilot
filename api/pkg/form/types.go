package form

import (
	"github.com/nrc-no/core/pkg/i18n"
	"github.com/nrc-no/core/pkg/validation"
	"net/url"
)

type ControlType string

const (
	Text     ControlType = "text"
	Email    ControlType = "email"
	Phone    ControlType = "tel"
	URL      ControlType = "url"
	Date     ControlType = "date"
	Textarea ControlType = "textarea"
	Dropdown ControlType = "dropdown"
	Checkbox ControlType = "checkbox"
	Radio    ControlType = "radio"
	Taxonomy ControlType = "taxonomy"
	File     ControlType = "file"
	Number   ControlType = "number"
)

var ControlTypes = []ControlType{Text, Email, Phone, URL, Date, Textarea, Dropdown, Checkbox, Radio, Taxonomy, File, Number}

type Control struct {
	Name            string            `json:"name" bson:"name"`
	Type            ControlType       `json:"type" bson:"type"`
	Label           i18n.Strings      `json:"label" bson:"label"`
	Validation      ControlValidation `json:"validation" bson:"validation"`
	DefaultValue    []string          `json:"defaultValue" bson:"defaultValue"`
	Description     i18n.Strings      `json:"description" bson:"description"`
	Placeholder     i18n.Strings      `json:"placeholder" bson:"placeholder"`
	Multiple        bool              `json:"multiple" bson:"multiple"`
	Options         []i18n.Strings    `json:"options" bson:"options"`
	CheckboxOptions []CheckboxOption  `json:"checkboxOptions" bson:"checkboxOptions"`
	Readonly        bool              `json:"readonly" bson:"readonly"`
}

func NewControl(name string, typ ControlType, label i18n.Strings, required bool) *Control {
	return &Control{
		Name:       name,
		Type:       typ,
		Label:      label,
		Validation: ControlValidation{Required: required},
	}
}

type CheckboxOption struct {
	Label    i18n.Strings `json:"label" bson:"label"`
	Value    string       `json:"value" bson:"value"`
	Required bool         `json:"required" bson:"required"`
}

type ControlValidation struct {
	Required bool `json:"required" bson:"required"`
}

type ValuedControl struct {
	Control
	Value  []string
	Errors *validation.ErrorList
}

type ValuedForm struct {
	Controls []ValuedControl
	Errors   *validation.ErrorList
}

type Form struct {
	Controls []Control `json:"controls" bson:"controls"`
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
		ctrl := ValuedControl{Control: control}
		if values != nil {
			value := values[control.Name]
			ctrl.Value = value
		}
		if errors != nil {
			errs := errors.FindFamily(control.Name)
			ctrl.Errors = errs
		}
		valuedControls = append(valuedControls, ctrl)
	}
	var result ValuedForm
	result.Controls = valuedControls
	if errors != nil {
		errs := errors.Find("")
		result.Errors = errs
	}
	return result
}
