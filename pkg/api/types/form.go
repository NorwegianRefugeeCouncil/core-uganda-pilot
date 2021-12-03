package types

import (
	"fmt"
	"github.com/nrc-no/core/pkg/utils/sets"
)

// FormDefinition represents the definition of a Form for data collection.
type FormDefinition struct {
	// ID is the unique ID of the FormDefinition
	ID string `json:"id" yaml:"id"`
	// Code of the FormDefinition
	// TODO remove this. It's not used yet.
	Code string `json:"code,omitempty" yaml:"code,omitempty"`
	// DatabaseID of the FormDefinition
	DatabaseID string `json:"databaseId" yaml:"databaseId"`
	// FolderID of the FormDefinition. If the FolderID is empty,
	// this means that the FormDefinition exists at the root
	// of the DatabaseID
	FolderID string `json:"folderId,omitempty" yaml:"folderId,omitempty"`
	// Name of the FormDefinition
	Name string `json:"name" yaml:"name"`
	// Fields that constitute the FormDefinition
	Fields FieldDefinitions `json:"fields" yaml:"fields"`
}

// GetID implements FormInterface.GetID
func (f *FormDefinition) GetID() string {
	return f.ID
}

// GetFields implements FormInterface.GetFields
func (f *FormDefinition) GetFields() FieldDefinitions {
	return f.Fields
}

// GetParentID implements FormInterface.GetParentID
func (f *FormDefinition) GetParentID() *string {
	return nil
}

// GetFieldByID finds and returns a field by ID
func (f FormDefinition) GetFieldByID(fieldID string) (*FieldDefinition, error) {
	for _, field := range f.Fields {
		if field.ID == fieldID {
			return field, nil
		}
	}
	return nil, fmt.Errorf("could not find field with id %s", fieldID)
}

// RemoveFieldByID removes a field from the form
func (f *FormDefinition) RemoveFieldByID(fieldID string) (*FieldDefinition, error) {
	var result []*FieldDefinition
	var fld *FieldDefinition
	for _, field := range f.Fields {
		if field.ID == fieldID {
			fld = field
			continue
		}
		result = append(result, field)
	}
	if fld == nil {
		return nil, fmt.Errorf("could not find field with id %s", fieldID)
	}
	f.Fields = result
	return fld, nil
}

// GetFormInterface will return a FormInterface for the given form or sub form ID.
func (f *FormDefinition) GetFormInterface(formOrSubFormID string) (FormInterface, error) {
	// if the given form ID == the root form ID, just return the root form ID
	if f.ID == formOrSubFormID {
		return f, nil
	}
	// recursively find the sub form with the given ID
	subFormInterface := findSubFormInterface(f.ID, formOrSubFormID, f.Fields)
	if subFormInterface == nil {
		return nil, fmt.Errorf("could not find form or subform with id %s", formOrSubFormID)
	}
	return subFormInterface, nil
}

// FormDefinitionList represents a list of FormDefinition
type FormDefinitionList struct {
	Items []*FormDefinition `json:"items" yaml:"items"`
}

// NewFormDefinitionList creates a new FormDefinitionList
func NewFormDefinitionList(items ...*FormDefinition) *FormDefinitionList {
	return &FormDefinitionList{Items: append([]*FormDefinition{}, items...)}
}

// GetFormByID finds a FormDefinition by ID
func (f *FormDefinitionList) GetFormByID(formID string) (*FormDefinition, error) {
	for _, item := range f.Items {
		if item.ID == formID {
			return item, nil
		}
	}
	return nil, fmt.Errorf("form definition with id %s not found", formID)
}

// Len returns the length of the FormDefinitionList
func (f *FormDefinitionList) Len() int {
	return len(f.Items)
}

// IsEmpty returns whether a FormDefinition is empty or not
func (f *FormDefinitionList) IsEmpty() bool {
	return f.Len() == 0
}

// GetAtIndex returns the FormDefinition at the given index
func (f *FormDefinitionList) GetAtIndex(index int) *FormDefinition {
	return f.Items[index]
}

// GetAllFormsAndSubFormIDs returns the IDs of all the FormDefinition and their SubForms (recursively), if any
func (f FormDefinition) GetAllFormsAndSubFormIDs() sets.String {
	ids := sets.NewString(f.ID)
	return addFormsAndSubFormIDsInternal(f.Fields, ids)
}

// addFormsAndSubFormIDsInternal is an internal method used by GetAllFormsAndSubFormIDs to recursively
// add the form and subform ids.
func addFormsAndSubFormIDsInternal(fields []*FieldDefinition, ids sets.String) sets.String {
	for _, field := range fields {
		if field.FieldType.SubForm != nil {
			ids.Insert(field.FieldType.SubForm.ID)
			addFormsAndSubFormIDsInternal(field.FieldType.SubForm.Fields, ids)
		}
	}
	return ids
}
