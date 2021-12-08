package types

import (
	"fmt"
	"reflect"
)

// FormInterface is a common interface between a FormDefinition and a FieldTypeSubForm.
// This way, we can treat both the "root" FormDefinition and the child SubForm as
// a single type.
type FormInterface interface {
	FormReference

	// GetFields returns the form/subform fields
	GetFields() FieldDefinitions

	// IsSubForm returns whether the FormInterface is a sub form or not
	IsSubForm() bool

	// FindSubForm will recursively try to find a form or subform
	// with the given ID
	// TODO remove this from the FormInterface. Perhaps add a SubFormInterface
	FindSubForm(subFormId string) (FormInterface, error)
}

// formInterface is the implementation of FormInterface
type formInterface struct {
	parent     FormInterface
	id         string
	databaseId string
	fields     FieldDefinitions
}

// GetFormID implements FormInterface.GetID
func (f *formInterface) GetFormID() string {
	return f.id
}

// GetFields implements FormInterface.GetFields
func (f *formInterface) GetFields() FieldDefinitions {
	return f.fields
}

// HasOwner implements FormInterface.HasOwner
func (f *formInterface) IsSubForm() bool {
	parentValue := reflect.ValueOf(f.parent)
	return parentValue.Kind() == reflect.Ptr && !parentValue.IsNil()
}

// GetDatabaseID implements FormInterface.GetDatabaseID
func (f *formInterface) GetDatabaseID() string {
	return f.databaseId
}

// FindSubForm implements FormInterface.FindSubForm
func (f *formInterface) FindSubForm(subFormId string) (FormInterface, error) {
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
func (f *formInterface) childFormInterface(field *FieldDefinition) FormInterface {
	return newFormInterface(f, f.databaseId, field.ID, field.FieldType.SubForm.Fields)
}

// newFormInterface returns a new instance of a FormInterface
func newFormInterface(parent FormInterface, databaseId, id string, fields FieldDefinitions) FormInterface {
	return &formInterface{
		databaseId: databaseId,
		parent:     parent,
		id:         id,
		fields:     fields,
	}
}
