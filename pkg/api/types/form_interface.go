package types

import (
	"fmt"
)

// FormInterface is a common interface between a FormDefinition and a FieldTypeSubForm.
// This way, we can treat both the "root" FormDefinition and the child SubForm as
// a single type.
type FormInterface interface {
	FormReference

	// GetFields returns the form/subform fields
	GetFields() FieldDefinitions

	// FindSubForm will recursively try to find a form or subform
	// with the given ID
	// TODO remove this from the FormInterface. Add a method without received FindSubForm(f FormInterface) SubFormInterface
	// Also dont make it recursive
	FindSubForm(subFormId string) (FormInterface, error)
}

type SubFormInterface interface {
	FormInterface
	GetOwnerFormID() string
}

// subFormInterface is the implementation of FormInterface
type subFormInterface struct {
	ownerFormID string
	id          string
	databaseId  string
	fields      FieldDefinitions
}

// GetFormID implements FormInterface.GetID
func (f *subFormInterface) GetFormID() string {
	return f.id
}

// GetFields implements FormInterface.GetFields
func (f *subFormInterface) GetFields() FieldDefinitions {
	return f.fields
}

// GetDatabaseID implements FormInterface.GetDatabaseID
func (f *subFormInterface) GetDatabaseID() string {
	return f.databaseId
}

// GetOwnerFormID implements SubFormInterface.GetOwnerFormID
func (f *subFormInterface) GetOwnerFormID() string {
	return f.ownerFormID
}

// FindSubForm implements FormInterface.FindSubForm
func (f *subFormInterface) FindSubForm(subFormId string) (FormInterface, error) {
	var foundInterface FormInterface
	for _, field := range f.fields {
		isSubForm, err := field.FieldType.IsKind(FieldKindSubForm)
		if err != nil {
			return nil, err
		}
		if field.ID == subFormId {
			if !isSubForm {
				return nil, fmt.Errorf("field '%s' is not of kind SubForm", subFormId)
			}
			subFormInterface := f.childFormInterface(field)
			return subFormInterface, nil
		}
		if !isSubForm {
			continue
		}
		subFormInterface := f.childFormInterface(field)
		foundInterface, err = subFormInterface.FindSubForm(subFormId)
		if err != nil {
			return nil, err
		}
		if foundInterface != nil {
			return foundInterface, nil
		}
	}
	return nil, nil
}

// childFormInterface returns a FormInterface for the given SubForm field
func (f *subFormInterface) childFormInterface(field *FieldDefinition) FormInterface {
	return newFormInterface(f.GetFormID(), f.databaseId, field.ID, field.FieldType.SubForm.Fields)
}

// newFormInterface returns a new instance of a FormInterface
func newFormInterface(ownerId string, databaseId, id string, fields FieldDefinitions) FormInterface {
	return &subFormInterface{
		databaseId:  databaseId,
		ownerFormID: ownerId,
		id:          id,
		fields:      fields,
	}
}
