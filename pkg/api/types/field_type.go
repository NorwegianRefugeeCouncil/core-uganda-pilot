package types

type FieldType struct {
	Text          *FieldTypeText          `json:"text,omitempty"`
	Reference     *FieldTypeReference     `json:"reference,omitempty"`
	SubForm       *FieldTypeSubForm       `json:"subForm,omitempty"`
	MultilineText *FieldTypeMultilineText `json:"multilineText,omitempty"`
	Date          *FieldTypeDate          `json:"date,omitempty"`
	Quantity      *FieldTypeQuantity      `json:"quantity,omitempty"`
	SingleSelect  *FieldTypeSingleSelect  `json:"singleSelect,omitempty"`
}

type FieldTypeReference struct {
	DatabaseID string `json:"databaseId,omitempty"`
	FormID     string `json:"formId,omitempty"`
}

type FieldTypeText struct{}

type FieldTypeMultilineText struct{}

type FieldTypeDate struct{}

type FieldTypeQuantity struct{}

type FieldTypeSingleSelect struct{}

type FieldTypeSubForm struct {
	ID     string             `json:"id"`
	Name   string             `json:"name"`
	Code   string             `json:"code"`
	Fields []*FieldDefinition `json:"fields,omitempty"`
}

func (f *FieldTypeSubForm) GetID() string {
	return f.ID
}

func (f *FieldTypeSubForm) GetFields() []*FieldDefinition {
	return f.Fields
}

type FieldKind string

const (
	FieldKindUnknown       FieldKind = "unknown"
	FieldKindText          FieldKind = "text"
	FieldKindSubForm       FieldKind = "subform"
	FieldKindReference     FieldKind = "reference"
	FieldKindMultilineText FieldKind = "multilineText"
	FieldKindDate          FieldKind = "date"
	FieldKindQuantity      FieldKind = "quantity"
	FieldKindSingleSelect  FieldKind = "singleSelect"
)
