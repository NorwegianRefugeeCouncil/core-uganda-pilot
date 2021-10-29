package types

import (
	"fmt"
	"github.com/nrc-no/core/pkg/sets"
)

type FieldDefinition struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Options     []string  `json:"options"`
	Key         bool      `json:"key"`
	Required    bool      `json:"required"`
	FieldType   FieldType `json:"fieldType"`
}

func (f FieldDefinition) IsKind(kind FieldKind) bool {
	return f.GetKind() == kind
}

func (f FieldDefinition) IsReferenceField() bool {
	return f.IsKind(FieldKindReference)
}

func (f FieldDefinition) IsSubFormField() bool {
	return f.IsKind(FieldKindSubForm)
}

func (f FieldDefinition) GetKind() FieldKind {
	if f.FieldType.Reference != nil {
		return FieldKindReference
	}
	if f.FieldType.Text != nil {
		return FieldKindText
	}
	if f.FieldType.SubForm != nil {
		return FieldKindSubForm
	}
	return FieldKindUnknown
}

type FieldDefinitions []*FieldDefinition

func (f FieldDefinitions) GetByID(fieldId string) (*FieldDefinition, error) {
	for _, definition := range f {
		if definition.ID == fieldId {
			return definition, nil
		}
	}
	return nil, fmt.Errorf("field definition with id %s not found", fieldId)
}

func (f FieldDefinitions) OfKind(fieldKind FieldKind) FieldDefinitions {
	var result FieldDefinitions
	for _, definition := range f {
		if definition.GetKind() == fieldKind {
			result = append(result, definition)
		}
	}
	return result
}

func (f FieldDefinitions) ThatAreKeys(isKey bool) FieldDefinitions {
	var result FieldDefinitions
	for _, definition := range f {
		if definition.Key == isKey {
			result = append(result, definition)
		}
	}
	return result
}

func (f FieldDefinitions) FieldIDs() sets.String {
	result := sets.NewString()
	for _, definition := range f {
		result.Insert(definition.ID)
	}
	return result
}

func (f FieldDefinitions) Expand(referencedForms *FormDefinitionList) (FieldDefinitions, error) {

	result := append(f[:])
	walk := 0
	for {
		if walk == len(result) {
			break
		}
		field := result[walk]
		if !field.IsReferenceField() {
			walk++
			continue
		}

		referencedForm, err := referencedForms.GetForm(field.FieldType.Reference.FormID)
		if err != nil {
			return nil, err
		}

		referencedFormKeyFields := referencedForm.GetFields().ThatAreKeys(true)

		newFieldLen := walk - 1
		if newFieldLen < 0 {
			newFieldLen = 0
		}
		newFields := make(FieldDefinitions, newFieldLen)
		if walk > 0 {
			copy(newFields, result[:walk])
		}
		newFields = append(newFields, referencedFormKeyFields...)
		if walk < len(result) {
			for _, definition := range result[walk+1:] {
				newFields = append(newFields, definition)
			}
		}
		result = newFields

		for _, definition := range result {
			fmt.Println(definition.ID)
		}
		fmt.Println("\n==========")

	}

	return result, nil
}
