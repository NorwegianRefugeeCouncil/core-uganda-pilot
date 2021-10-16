package types

import "fmt"

type FormDefinition struct {
	ID         string             `json:"id"`
	Code       string             `json:"code"`
	DatabaseID string             `json:"databaseId,omitempty"`
	FolderID   string             `json:"folderId"`
	Name       string             `json:"name,omitempty"`
	Fields     []*FieldDefinition `json:"fields,omitempty"`
}

type FormInterface interface {
	GetID() string
	GetFields() []*FieldDefinition
	GetParentID() *string
}

func (f *FormDefinition) GetID() string {
	return f.ID
}

func (f *FormDefinition) GetFields() []*FieldDefinition {
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

func (f *formIntf) GetFields() []*FieldDefinition {
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

func (f FormDefinition) GetFieldByID(fieldID string) *FieldDefinition {
	for _, field := range f.Fields {
		if field.ID == fieldID {
			return field
		}
	}
	return &FieldDefinition{}
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
