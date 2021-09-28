package form

import (
	"github.com/nrc-no/core/pkg/i18n"
	"github.com/nrc-no/core/pkg/validation"
)

// Control describes the attributes of a form input control. Some fields, like Name, Placeholder and Readonly are
//intended as 1:1 mappings to HTML form element attributes. Others like Validation, Label and CheckboxOptions correspond
// to cross-cutting concerns (Validation), sibling elements (Label) or richer data structures (CheckboxOptions).
type Control struct {
	Name            string            `json:"name" bson:"name"`
	Type            ControlType       `json:"type" bson:"type"`
	Placeholder     i18n.Strings      `json:"placeholder" bson:"placeholder"`
	Multiple        bool              `json:"multiple" bson:"multiple"`
	Readonly        bool              `json:"readonly" bson:"readonly"`
	Label           i18n.Strings      `json:"label" bson:"label"`
	Description     i18n.Strings      `json:"description" bson:"description"`
	DefaultValue    []string          `json:"defaultValue" bson:"defaultValue"`
	Options         []i18n.Strings    `json:"options" bson:"options"`
	CheckboxOptions []CheckboxOption  `json:"checkboxOptions" bson:"checkboxOptions"`
	Validation      ControlValidation `json:"validation" bson:"validation"`
}

// ControlType is the type of a form input element. Some ControlTypes correspond to built-in HTML input type such as
// <input type="email"> or <input type="date"> other refer to custom-build elements ("boolean", "taxonomy"). In all
// cases, it is up to the package user to associate a given ControlType with the appropriate representation on
// the front end.
type ControlType string

const (
	Text     ControlType = "text"
	Email    ControlType = "email"
	Phone    ControlType = "tel"
	URL      ControlType = "url"
	Date     ControlType = "date"
	Textarea ControlType = "textarea"
	Dropdown ControlType = "dropdown"
	Boolean  ControlType = "boolean"
	Checkbox ControlType = "checkbox"
	Radio    ControlType = "radio"
	Taxonomy ControlType = "taxonomy"
	File     ControlType = "file"
)

var ControlTypes = []ControlType{Text, Email, Phone, URL, Date, Textarea, Dropdown, Boolean, Checkbox, Radio, Taxonomy, File}

func NewControl(name string, typ ControlType, label i18n.Strings, required bool) *Control {
	return &Control{
		Name:       name,
		Type:       typ,
		Label:      label,
		Validation: ControlValidation{Required: required},
	}
}

// CheckboxOption describes the attributes of a checkbox input element intended as a Checkbox and Radio ControlType
//child element.
type CheckboxOption struct {
	Label    i18n.Strings `json:"label" bson:"label"`
	Value    string       `json:"value" bson:"value"`
	Required bool         `json:"required" bson:"required"`
}

// ControlValidation relates validation constraints applied to a Control
type ControlValidation struct {
	Required bool `json:"required" bson:"required"`
}

func (f *Form) FindControlByName(name string) *Control {
	for _, control := range f.Controls {
		if control.Name == name {
			return &control
		}
	}
	return nil
}

// ValuedControl is Control which has undergone validation and contains a Value as well as possible validation Errors.
type ValuedControl struct {
	Control
	Value  []string
	Errors *validation.ErrorList
}
