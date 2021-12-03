package types

import (
	"fmt"
	"github.com/nrc-no/core/pkg/utils/sets"
)

// FieldDefinition usually represents a question in a FormDefinition.
// A FieldDefinition defines the name, description, boundaries of data collection.
type FieldDefinition struct {
	// ID is the ID of the FieldDefinition
	ID string `json:"id,omitempty" yaml:"id,omitempty"`
	// Code is the unique Code of the FieldDefinition within the FormDefinition
	Code string `json:"code,omitempty" yaml:"code,omitempty"`
	// Name is the Name of the FieldDefinition
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Description is a helpful text helping the users to understand the question
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Options TODO: Remove this, put inside of FieldTypeMultiSelect / FieldTypeSelect
	Options []string `json:"options,omitempty" yaml:"options,omitempty"`
	// Key indicates that the FieldDefinition is part of the Unique Keys for the FormDefinition.
	// When a FormDefinition is created with Key fields, this means that there will be no
	// two records with the same combination of Key field values.
	//
	// For example, if a FormDefinition has 2 key fields, "Year" and "Month", then there could
	// be no two records with "2021" and "January".
	Key bool `json:"key" yaml:"key,omitempty"`
	// Required indicates that the user must enter data for that FieldDefinition
	Required bool `json:"required" yaml:"required"`
	// FieldType contains the type of FieldDefinition
	FieldType FieldType `json:"fieldType,omitempty" yaml:"fieldType,omitempty"`
}

// IsKind is a helper method that checks if the field is of the given FieldKind
func (f FieldDefinition) IsKind(kind FieldKind) bool {
	return f.GetKind() == kind
}

// IsReferenceField returns whether the field is a FieldTypeReference field or not
func (f FieldDefinition) IsReferenceField() bool {
	return f.IsKind(FieldKindReference)
}

// IsSubFormField returns whether a the field is a FieldTypeSubForm or not
func (f FieldDefinition) IsSubFormField() bool {
	return f.IsKind(FieldKindSubForm)
}

// GetKind returns the FieldKind of the field
func (f FieldDefinition) GetKind() FieldKind {
	if f.FieldType.Reference != nil {
		return FieldKindReference
	}
	if f.FieldType.Text != nil {
		return FieldKindText
	}
	if f.FieldType.Date != nil {
		return FieldKindDate
	}
	if f.FieldType.Quantity != nil {
		return FieldKindQuantity
	}
	if f.FieldType.MultilineText != nil {
		return FieldKindMultilineText
	}
	if f.FieldType.Month != nil {
		return FieldKindMonth
	}
	if f.FieldType.SubForm != nil {
		return FieldKindSubForm
	}
	return FieldKindUnknown
}

// FieldDefinitions represent a list of FieldDefinition
type FieldDefinitions []*FieldDefinition

// GetByID returns a FieldDefinition by its FieldDefinition.ID
func (f FieldDefinitions) GetByID(fieldId string) (*FieldDefinition, error) {
	for _, definition := range f {
		if definition.ID == fieldId {
			return definition, nil
		}
	}
	return nil, fmt.Errorf("field definition with id %s not found", fieldId)
}

// OfKind filters and returns the FieldDefinitions that are of the given FieldKind
func (f FieldDefinitions) OfKind(fieldKind FieldKind) FieldDefinitions {
	var result FieldDefinitions
	for _, definition := range f {
		if definition.GetKind() == fieldKind {
			result = append(result, definition)
		}
	}
	return result
}

// ThatAreKeys filters and returns the FieldDefinitions that have FieldDefinition.Key = true
func (f FieldDefinitions) ThatAreKeys(isKey bool) FieldDefinitions {
	var result FieldDefinitions
	for _, definition := range f {
		if definition.Key == isKey {
			result = append(result, definition)
		}
	}
	return result
}

// FieldIDs returns the of FieldDefinition.ID for all the FieldDefinitions
func (f FieldDefinitions) FieldIDs() sets.String {
	result := sets.NewString()
	for _, definition := range f {
		result.Insert(definition.ID)
	}
	return result
}
