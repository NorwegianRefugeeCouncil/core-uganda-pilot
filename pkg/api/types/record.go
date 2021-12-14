package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RecordRef represents a key that allows referencing a single Record.
type RecordRef struct {
	// ID is the Record.ID
	ID string `json:"id"`
	// DatabaseID is the Record.DatabaseID
	DatabaseID string `json:"databaseId"`
	// FormID is the Record.FormID
	FormID string `json:"formId"`
}

// FormRef represents a key that allows referencing a single FormDefinition.
type FormRef struct {
	// DatabaseID represents the FormDefinition.DatabaseID
	DatabaseID string `json:"databaseId"`
	// FormID represents the FormDefinition.ID
	FormID string `json:"formId"`
}

// GetDatabaseID implements FormReference.GetDatabaseID
func (f FormRef) GetDatabaseID() string {
	return f.DatabaseID
}

// GetFormID implements FormReference.GetFormID
func (f FormRef) GetFormID() string {
	return f.FormID
}

// Record represents an entry in a Form.
type Record struct {
	// ID of the record
	ID string `json:"id"`
	// Seq of the Record. This value is automatically increased by the database.
	// The presence of this field allows us to sort the table by insertion order.
	Seq int64 `json:"seq"`
	// DatabaseID of the Record
	DatabaseID string `json:"databaseId"`
	// FormID of the Record. Represents in which Form this record belongs.
	FormID string `json:"formId"`
	// OwnerID represents the owner of the Record. In cases where
	// a Record is part of a SubForm, this field records the "Owner" form ID.
	OwnerID *string `json:"ownerFormID"`
	// Values is a list of values that correspond to the FormDefinition.Fields.
	Values FieldValues `json:"values"`
}

type Records []Record

type StringOrArrayValue uint8

const (
	StringValue StringOrArrayValue = iota + 1
	ArrayValue
	NullValue
)

type StringOrArray struct {
	Kind        StringOrArrayValue
	StringValue string
	ArrayValue  []string
}

func (s StringOrArray) GetValue() interface{} {
	if s.Kind == StringValue {
		return s.StringValue
	} else if s.Kind == ArrayValue {
		return s.ArrayValue
	}
	return nil
}

func NewStringValue(value string) StringOrArray {
	return StringOrArray{
		Kind:        StringValue,
		StringValue: value,
	}
}

func NewArrayValue(value []string) StringOrArray {
	return StringOrArray{
		Kind:       ArrayValue,
		ArrayValue: value,
	}
}

func NewNullValue() StringOrArray {
	return StringOrArray{
		Kind: NullValue,
	}
}

func (s StringOrArray) MarshalJSON() ([]byte, error) {
	switch s.Kind {
	case NullValue:
		return json.Marshal(nil)
	case ArrayValue:
		return json.Marshal(append(s.ArrayValue))
	case StringValue:
		return json.Marshal(s.StringValue)
	default:
		return nil, fmt.Errorf("unknown value kind")
	}
}

func (s *StringOrArray) UnmarshalJSON(b []byte) error {
	jsonStr := string(b)
	if jsonStr == "null" {
		s.Kind = NullValue
		return nil
	}

	if strings.HasPrefix(jsonStr, "[") {
		s.Kind = ArrayValue
		return json.Unmarshal(b, &s.ArrayValue)
	}

	s.Kind = StringValue
	return json.Unmarshal(b, &s.StringValue)
}

type FieldValue struct {
	FieldID string        `json:"fieldId"`
	Value   StringOrArray `json:"value"`
}

func NewFieldStringValue(fieldID string, value string) FieldValue {
	return FieldValue{
		FieldID: fieldID,
		Value:   NewStringValue(value),
	}
}

func NewFieldArrayValue(fieldID string, value []string) FieldValue {
	return FieldValue{
		FieldID: fieldID,
		Value:   NewArrayValue(value),
	}
}
func NewFieldNullValue(fieldID string) FieldValue {
	return FieldValue{
		FieldID: fieldID,
		Value:   NewNullValue(),
	}
}

type FieldValues []FieldValue

func (f FieldValues) Find(fieldID string) (FieldValue, bool) {
	for _, value := range f {
		if value.FieldID == fieldID {
			return value, true
		}
	}
	return FieldValue{}, false
}

func (f FieldValues) FindIndex(fieldID string) int {
	for i, value := range f {
		if value.FieldID == fieldID {
			return i
		}
	}
	return -1
}

func (f FieldValues) SetValue(fieldID string, value StringOrArray) FieldValues {
	values := f
	for i, v := range values {
		if v.FieldID == fieldID {
			values[i].Value = value
			return values
		}
	}
	values = append(values, FieldValue{FieldID: fieldID, Value: value})
	return values
}

func (f FieldValues) SetValueForFieldName(formDefinition *FormDefinition, fieldName string, value StringOrArray) (FieldValues, error) {
	values := f
	field, err := formDefinition.Fields.GetFieldByName(fieldName)
	if err != nil {
		return FieldValues{}, err
	}
	values = values.SetValue(field.ID, value)
	return values, nil
}

// RecordList represents a list of Record
type RecordList struct {
	Items []*Record `json:"items"`
}

// RecordListOptions represents the options for listing Record.
type RecordListOptions struct {
	DatabaseID string
	FormID     string
}
