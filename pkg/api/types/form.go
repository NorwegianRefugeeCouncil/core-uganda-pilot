package types

import (
	"fmt"
)

// FormDefinition represents the definition of a Form for data collection.
type FormDefinition struct {
	// ID is the unique ID of the FormDefinition
	ID string `json:"id"`
	// Code of the FormDefinition
	// TODO remove this. It's not used yet.
	Code string `json:"code"`
	// DatabaseID of the FormDefinition
	DatabaseID string `json:"databaseId,omitempty"`
	// FolderID of the FormDefinition. If the FolderID is empty,
	// this means that the FormDefinition exists at the root
	// of the DatabaseID
	FolderID string `json:"folderId"`
	// Name of the FormDefinition
	Name string `json:"name,omitempty"`
	// Fields that constitute the FormDefinition
	Fields FieldDefinitions `json:"fields"`
}

type FormInterface interface {
	GetID() string
	GetFields() FieldDefinitions
	GetParentID() *string
}

func (f *FormDefinition) GetID() string {
	return f.ID
}

func (f *FormDefinition) GetFields() FieldDefinitions {
	return f.Fields
}

func (f *FormDefinition) GetParentID() *string {
	return nil
}

type formIntf struct {
	parentId string
	subForm  *FieldTypeSubForm
}

func (f *formIntf) GetID() string {
	return f.subForm.ID
}

func (f *formIntf) GetFields() FieldDefinitions {
	return f.subForm.Fields
}

func (f *formIntf) GetParentID() *string {
	return &f.parentId
}

func findFormIntf(parentId, id string, fields []*FieldDefinition) FormInterface {
	for _, field := range fields {
		subForm := field.FieldType.SubForm
		if subForm != nil {
			if subForm.ID == id {
				return &formIntf{
					parentId: parentId,
					subForm:  subForm,
				}
			}
			var childF = findFormIntf(subForm.ID, id, subForm.Fields)
			if childF != nil {
				return childF
			}
		}
	}
	return nil
}

func (f *FormDefinition) GetFormInterface(formOrSubFormID string) (FormInterface, error) {
	if f.ID == formOrSubFormID {
		return f, nil
	}
	childF := findFormIntf(f.ID, formOrSubFormID, f.Fields)
	if childF == nil {
		return nil, fmt.Errorf("could not find form or subform with id %s", formOrSubFormID)
	}
	return childF, nil
}

type FormDefinitionList struct {
	Items []*FormDefinition `json:"items"`
}

func NewFormDefinitionList(items ...*FormDefinition) *FormDefinitionList {
	return &FormDefinitionList{Items: append([]*FormDefinition{}, items...)}
}

func (f *FormDefinitionList) GetForm(formID string) (*FormDefinition, error) {
	for _, item := range f.Items {
		if item.ID == formID {
			return item, nil
		}
	}
	return nil, fmt.Errorf("form definition with id %s not found", formID)
}

func (f *FormDefinitionList) Len() int {
	return len(f.Items)
}

func (f *FormDefinitionList) Empty() bool {
	return f.Len() == 0
}

func (f *FormDefinitionList) GetAtIndex(index int) *FormDefinition {
	return f.Items[index]
}

func (f FormDefinition) GetFieldByID(fieldID string) *FieldDefinition {
	for _, field := range f.Fields {
		if field.ID == fieldID {
			return field
		}
	}
	return &FieldDefinition{}
}

func (f FormDefinition) GetAllFormsAndSubFormIDs() []string {
	ids := []string{f.ID}
	return getAllFormsAndSubFormIDsInternal(f.Fields, ids)
}

func getAllFormsAndSubFormIDsInternal(fields []*FieldDefinition, ids []string) []string {
	for _, field := range fields {
		if field.FieldType.SubForm != nil {
			ids = append(ids, field.FieldType.SubForm.ID)
			ids = append(ids, getAllFormsAndSubFormIDsInternal(field.FieldType.SubForm.Fields, ids)...)
		}
	}
	return ids
}

func (f FormDefinition) GetFieldByName(fieldName string) *FieldDefinition {
	for _, field := range f.Fields {
		if field.Name == fieldName {
			return field
		}
	}
	return &FieldDefinition{}
}

func (f *FormDefinition) RemoveField(fieldName string) *FieldDefinition {
	var result []*FieldDefinition
	var fld *FieldDefinition
	for _, field := range f.Fields {
		if field.Name == fieldName {
			fld = field
			continue
		}
		result = append(result, field)
	}
	f.Fields = result
	return fld
}
