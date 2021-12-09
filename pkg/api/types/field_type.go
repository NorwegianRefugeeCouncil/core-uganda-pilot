package types

import (
	"errors"
	"fmt"
	"reflect"
)

// FieldType is a struct that contains the FieldType for a given FieldDefinition
// Only one of the fields might be specified. For example, a FieldType
// cannot have both FieldType.Text and FieldType.Reference defined. Only one is allowed.
type FieldType struct {
	// Text represents the configuration for a text field
	Text *FieldTypeText `json:"text,omitempty" yaml:"text,omitempty"`
	// Reference represents the configuration for a reference field
	Reference *FieldTypeReference `json:"reference,omitempty" yaml:"reference,omitempty"`
	// SubForm represents the configuration for a sub form field
	SubForm *FieldTypeSubForm `json:"subForm,omitempty" yaml:"subForm,omitempty"`
	// MultilineText represents the configuration for a multiline text field
	MultilineText *FieldTypeMultilineText `json:"multilineText,omitempty" yaml:"multilineText,omitempty"`
	// Date represents the configuration for a date field
	Date *FieldTypeDate `json:"date,omitempty" yaml:"date,omitempty"`
	// Quantity represents the configuration for a quantity field
	Quantity *FieldTypeQuantity `json:"quantity,omitempty" yaml:"quantity,omitempty"`
	// Week represents the configuration for a week field
	Week *FieldTypeWeek `json:"week,omitempty" yaml:"week,omitempty"`
	// Month represents the configuration for a month field
	Month *FieldTypeMonth `json:"month,omitempty" yaml:"month,omitempty"`
	// SingleSelect represents the configuration for a single select field
	SingleSelect *FieldTypeSingleSelect `json:"singleSelect,omitempty" yaml:"singleSelect,omitempty"`
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

func (f FieldType) GetFieldKind() (FieldKind, error) {
	for kind, accessor := range fieldAccessors {
		field := accessor(f)
		value := reflect.ValueOf(field)
		if value.Kind() == reflect.Ptr && !value.IsNil() {
			return kind, nil
		}
	}
	return FieldKindUnknown, errors.New("failed to get field kind")
}
func (f FieldType) IsKind(kind FieldKind) (bool, error) {
	fieldKind, err := f.GetFieldKind()
	if err != nil {
		return false, err
	}
	return fieldKind == kind, nil
}

// FieldTypeReference represents a field that is a reference to a record in another FormDefinition
//
// For example, given a form "Countries" and a form "Projects".
// The "Projects" form might have a field "Country" that references the "Countries" form.
// In this case, when adding a record in the "Projects", the user would be prompted to select a
// country.
type FieldTypeReference struct {
	// DatabaseID represents the DatabaseID of the referenced FormDefinition
	DatabaseID string `json:"databaseId" yaml:"databaseId"`
	// FormID represents the FormID of the referenced FormDefinition
	FormID string `json:"formId" yaml:"formId"`
}

// FieldTypeText represents a textual field
type FieldTypeText struct{}

// FieldTypeMultilineText represents a multiline text field
type FieldTypeMultilineText struct{}

// FieldTypeDate represents a Date field (calendar date, no time/timezone)
type FieldTypeDate struct{}

// FieldTypeWeek represents a Week field (YYYYKWww)
type FieldTypeWeek struct{}

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
	// Options represent the different options that the user can select from
	Options []*SelectOption `json:"options,omitempty" yaml:"options,omitempty"`
}

// SelectOption represent an option for a FieldTypeSingleSelect or FieldTypeMultiSelect
type SelectOption struct {
	// ID of the option
	ID string `json:"id" yaml:"id"`
	// Name of the option
	Name string `json:"name" yaml:"name"`
}

// FieldTypeSubForm represents a field that contains a nested form.
// A user could attach multiple records of that subform to the "parent" record.
//
// For example, given a form "Projects", this form could have a subform "Monthly Deliveries".
// The "Monthly Deliveries". There could be multiple "Monthly Deliveries" for a single "Project".
type FieldTypeSubForm struct {
	// Fields represent the fields for the SubForm
	Fields FieldDefinitions `json:"fields,omitempty" yaml:"fields,omitempty"`
}

// GetFields  returns the FieldDefinitions for the subform
func (f *FieldTypeSubForm) GetFields() FieldDefinitions {
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
	FieldKindWeek
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
	FieldKindWeek: func(fieldType FieldType) interface{} {
		return fieldType.Week
	},
	FieldKindSingleSelect: func(fieldType FieldType) interface{} {
		return fieldType.SingleSelect
	},
}
