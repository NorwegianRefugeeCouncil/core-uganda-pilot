package types

import "fmt"

// FieldType is a struct that contains the FieldType for a given FieldDefinition
// Only one of the fields might be specified. For example, a FieldType
// cannot have both FieldType.Text and FieldType.Reference defined. Only one is allowed.
type FieldType struct {
	// Text represents the configuration for a text field
	Text *FieldTypeText `json:"text,omitempty"`
	// Reference represents the configuration for a reference field
	Reference *FieldTypeReference `json:"reference,omitempty"`
	// SubForm represents the configuration for a sub form field
	SubForm *FieldTypeSubForm `json:"subForm,omitempty"`
	// MultilineText represents the configuration for a multiline text field
	MultilineText *FieldTypeMultilineText `json:"multilineText,omitempty"`
	// Date represents the configuration for a date field
	Date *FieldTypeDate `json:"date,omitempty"`
	// Quantity represents the configuration for a quantity field
	Quantity *FieldTypeQuantity `json:"quantity,omitempty"`
	// Month represents the configuration for a month field
	Month *FieldTypeMonth `json:"month,omitempty"`
	// SingleSelect represents the configuration for a single select field
	SingleSelect *FieldTypeSingleSelect `json:"singleSelect,omitempty"`
}

const accessorMessage = `
No accessor for field %s is defined in types.FieldAccessors.
This means that you added a field type, but did not add the accessor for it.
Add the accessor in pkg/api/field_type`

func (f FieldType) GetFieldType(kind FieldKind) (interface{}, error) {
	accessor, ok := fieldAccessors[kind]
	if !ok {
		return nil, fmt.Errorf(accessorMessage, kind)
	}
	if accessor == nil {
		return nil, fmt.Errorf("the accessor for field kind %v is nil", kind)
	}
	return accessor(f), nil
}

// FieldTypeReference represents a field that is a reference to a record in another FormDefinition
//
// For example, given a form "Countries" and a form "Projects".
// The "Projects" form might have a field "Country" that references the "Countries" form.
// In this case, when adding a record in the "Projects", the user would be prompted to select a
// country.
type FieldTypeReference struct {
	// DatabaseID represents the DatabaseID of the referenced FormDefinition
	DatabaseID string `json:"databaseId,omitempty"`
	// FormID represents the FormID of the referenced FormDefinition
	FormID string `json:"formId,omitempty"`
}

// FieldTypeText represents a textual field
type FieldTypeText struct{}

// FieldTypeMultilineText represents a multiline text field
type FieldTypeMultilineText struct{}

// FieldTypeDate represents a Date field (calendar date, no time/timezone)
type FieldTypeDate struct{}

// FieldTypeMonth represents a Month field (YYYY-mm)
type FieldTypeMonth struct{}

// FieldTypeQuantity represents a quantity field. A quantity field is
// a "number" of something.
type FieldTypeQuantity struct {
	// TODO: add "units"
	// TODO: add "decimals"
}

// FieldTypeSingleSelect represents a field from which the user can select a single option
type FieldTypeSingleSelect struct {
	// todo: move FieldDefinition.Options here
}

// FieldTypeSubForm represents a field that contains a nested form.
// A user could attach multiple records of that subform to the "parent" record.
//
// For example, given a form "Projects", this form could have a subform "Monthly Deliveries".
// The "Monthly Deliveries". There could be multiple "Monthly Deliveries" for a single "Project".
type FieldTypeSubForm struct {
	// ID represents the ID of the sub form
	ID string `json:"id"`
	// Name represents the Name of the sub form
	Name string `json:"name"`
	// Code represents the unique Code for the subform Field
	Code string `json:"code"`
	// Fields represent the fields for the SubForm
	Fields []*FieldDefinition `json:"fields,omitempty"`
}

// GetID returns the ID of the sub form
func (f *FieldTypeSubForm) GetID() string {
	return f.ID
}

// GetFields  returns the FieldDefinitions for the subform
func (f *FieldTypeSubForm) GetFields() []*FieldDefinition {
	return f.Fields
}

var allFieldKinds []FieldKind

func GetAllFieldKinds() []FieldKind {
	result := make([]FieldKind, len(allFieldKinds))
	_ = copy(result, allFieldKinds)
	return result
}

func init() {
	for i := 0; i < len(_FieldKind_index)-1; i++ {
		allFieldKinds = append(allFieldKinds, FieldKind(i))
	}
}

// FieldKind is a struct that contains the different types of fields
type FieldKind int

//go:generate go run golang.org/x/tools/cmd/stringer -type=FieldKind

const (
	FieldKindUnknown FieldKind = iota
	FieldKindText
	FieldKindSubForm
	FieldKindReference
	FieldKindMultilineText
	FieldKindDate
	FieldKindQuantity
	FieldKindMonth
	FieldKindSingleSelect
)

var fieldAccessors = map[FieldKind]func(fieldType FieldType) interface{}{
	FieldKindUnknown: func(fieldType FieldType) interface{} {
		return nil
	},
	FieldKindText: func(fieldType FieldType) interface{} {
		return fieldType.Text
	},
	FieldKindSubForm: func(fieldType FieldType) interface{} {
		return fieldType.SubForm
	},
	FieldKindReference: func(fieldType FieldType) interface{} {
		return fieldType.Reference
	},
	FieldKindMultilineText: func(fieldType FieldType) interface{} {
		return fieldType.MultilineText
	},
	FieldKindDate: func(fieldType FieldType) interface{} {
		return fieldType.Date
	},
	FieldKindQuantity: func(fieldType FieldType) interface{} {
		return fieldType.Quantity
	},
	FieldKindMonth: func(fieldType FieldType) interface{} {
		return fieldType.Month
	},
	FieldKindSingleSelect: func(fieldType FieldType) interface{} {
		return fieldType.SingleSelect
	},
}
