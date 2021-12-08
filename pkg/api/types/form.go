package types

import (
	"fmt"
	"github.com/nrc-no/core/pkg/utils/sets"
)

// FormDefinition represents the definition of a Form for data collection.
type FormDefinition struct {
	// ID is the unique ID of the FormDefinition
	ID string `json:"id" yaml:"id"`
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

// GetDatabaseID implements FormInterface.GetDatabaseID
func (f FormDefinition) GetDatabaseID() string {
	return f.DatabaseID
}

// GetFormID implements FormInterface.GetFormID
func (f FormDefinition) GetFormID() string {
	return f.ID
}

// GetFields implements FormInterface.GetFields
func (f *FormDefinition) GetFields() FieldDefinitions {
	return f.Fields
}

// RemoveFieldByID removes a field from the form
func (f *FormDefinition) RemoveFieldByID(fieldID string) (*FieldDefinition, error) {
	var result FieldDefinitions
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

// FindSubForm implements FormInterface.FindFormInterface
func (f *FormDefinition) FindSubForm(subFormID string) (FormInterface, error) {
	formIntf := newFormInterface(f.ID, f.DatabaseID, f.ID, f.Fields)
	return formIntf.FindSubForm(subFormID)
}

// GetFormOrSubForm gets the form or subForm for the given id
func (f *FormDefinition) GetFormOrSubForm(formOrSubFormID string) (FormInterface, error) {
	if f.GetFormID() == formOrSubFormID {
		return f, nil
	}
	subForm, err := f.FindSubForm(formOrSubFormID)
	if err != nil {
		return nil, err
	}
	if subForm == nil {
		return nil, fmt.Errorf("failed to get form with id %s", formOrSubFormID)
	}
	return subForm, nil
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
func addFormsAndSubFormIDsInternal(fields FieldDefinitions, ids sets.String) sets.String {
	for _, field := range fields {
		if field.FieldType.SubForm != nil {
			ids.Insert(field.ID)
			addFormsAndSubFormIDsInternal(field.FieldType.SubForm.Fields, ids)
		}
	}
	return ids
}
