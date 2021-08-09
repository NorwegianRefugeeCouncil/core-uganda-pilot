package form

import "github.com/nrc-no/core/pkg/validation"

type FieldType string

const (
	Textarea      FieldType = "textarea"
	Text          FieldType = "textinput"
	Dropdown      FieldType = "dropdown"
	Radio         FieldType = "radio"
	Checkbox      FieldType = "checkbox"
	Email         FieldType = "email"
	Date          FieldType = "date"
	File          FieldType = "file"
	Time          FieldType = "time"
	TaxonomyInput FieldType = "taxonomyinput"
)

type FormElement struct {
	Type       FieldType             `json:"type" bson:"type"`
	Attributes FormElementAttribute  `json:"attributes" bson:"attributes"`
	Validation FormElementValidation `json:"validation" bson:"validation"`
	Errors     *validation.ErrorList `json:"errors"`
	Readonly   bool
}

type FormElementAttribute struct {
	Label           string           `json:"label" bson:"label"`
	Name            string           `json:"name" bson:"name"`
	Value           []string         `json:"value" bson:"value"`
	Description     string           `json:"description" bson:"description"`
	Placeholder     string           `json:"placeholder" bson:"placeholder"`
	Multiple        bool             `json:"multiple" bson:"multiple"`
	Options         []string         `json:"options" bson:"options"`
	CheckboxOptions []CheckboxOption `json:"checkboxOptions" bson:"checkboxOptions"`
}

type FormElementValidation struct {
	Required bool `json:"required" bson:"required"`
}

type CheckboxOption struct {
	Label    string `json:"label" bson:"label"`
	Required bool   `json:"required" bson:"required"`
}
