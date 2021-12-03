package types

// FormInterface is a common interface between a FormDefinition and a FieldTypeSubForm.
// This way, we can treat both the "root" FormDefinition and the child SubForm as
// a single type.
type FormInterface interface {
	// GetID returns the form or subform ID
	GetID() string
	// GetFields returns the form/subform fields
	GetFields() FieldDefinitions
	// GetParentID returns the parentID (if the form is a SubForm) or nil
	// This also supports multi-level sub forms
	GetParentID() *string
}

// subFormInterface is the implementation of FormInterface
type subFormInterface struct {
	parentId string
	subForm  *FieldTypeSubForm
}

// GetID implements FormInterface.GetID
func (f *subFormInterface) GetID() string {
	return f.subForm.ID
}

// GetFields implements FormInterface.GetFields
func (f *subFormInterface) GetFields() FieldDefinitions {
	return f.subForm.Fields
}

// GetParentID implements FormInterface.GetParentID
func (f *subFormInterface) GetParentID() *string {
	return &f.parentId
}

// findSubFormInterface will recursively iterate through the fields of a form and
// find a sub form with the given ID, and return a FormInterface
func findSubFormInterface(parentId, id string, fields []*FieldDefinition) FormInterface {
	for _, field := range fields {
		subForm := field.FieldType.SubForm
		if subForm != nil {
			if subForm.ID == id {
				return &subFormInterface{
					parentId: parentId,
					subForm:  subForm,
				}
			}
			var childF = findSubFormInterface(subForm.ID, id, subForm.Fields)
			if childF != nil {
				return childF
			}
		}
	}
	return nil
}
