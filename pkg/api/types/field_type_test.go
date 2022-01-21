package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const duplicateAccessorMsg = `it seems that the accessor for %s returned a value that was returned by another accessor.
Check that the accessor returns the proper value`

const nilAccessorMSg = `The accessor for FieldKind %s returned a nil value. That should not happen`

const unregisteredAccessor = `The field kind %s is not in the list of all field kinds. Run make gen to regenerate.`

func TestGetFieldKind(t *testing.T) {

	tests := []struct {
		name       string
		fieldType  FieldType
		expectErr  bool
		expectKind FieldKind
	}{
		{
			name:       "text",
			fieldType:  FieldType{Text: &FieldTypeText{}},
			expectKind: FieldKindText,
		}, {
			name:       "multilineText",
			fieldType:  FieldType{MultilineText: &FieldTypeMultilineText{}},
			expectKind: FieldKindMultilineText,
		}, {
			name:       "month",
			fieldType:  FieldType{Month: &FieldTypeMonth{}},
			expectKind: FieldKindMonth,
		}, {
			name:       "week",
			fieldType:  FieldType{Week: &FieldTypeWeek{}},
			expectKind: FieldKindWeek,
		}, {
			name:       "date",
			fieldType:  FieldType{Date: &FieldTypeDate{}},
			expectKind: FieldKindDate,
		}, {
			name:       "reference",
			fieldType:  FieldType{Reference: &FieldTypeReference{}},
			expectKind: FieldKindReference,
		}, {
			name:       "quantity",
			fieldType:  FieldType{Quantity: &FieldTypeQuantity{}},
			expectKind: FieldKindQuantity,
		}, {
			name:       "subform",
			fieldType:  FieldType{SubForm: &FieldTypeSubForm{}},
			expectKind: FieldKindSubForm,
		}, {
			name:       "singleSelect",
			fieldType:  FieldType{SingleSelect: &FieldTypeSingleSelect{}},
			expectKind: FieldKindSingleSelect,
		}, {
			name:       "multiSelect",
			fieldType:  FieldType{MultiSelect: &FieldTypeMultiSelect{}},
			expectKind: FieldKindMultiSelect,
		},
		{
			name:       "boolean",
			fieldType:  FieldType{Boolean: &FieldTypeBoolean{}},
			expectKind: FieldKindBoolean,
		},
	}

	var handledFieldKinds []FieldKind
	for _, tc := range tests {
		test := tc
		handledFieldKinds = append(handledFieldKinds, test.expectKind)
		t.Run(test.name, func(t *testing.T) {
			kind, err := test.fieldType.GetFieldKind()
			if test.expectErr && !assert.Error(t, err) {
				return
			}
			if !test.expectErr && !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, test.expectKind, kind)
		})
	}

	for _, kind := range GetAllFieldKinds() {
		found := false
		for _, handled := range handledFieldKinds {
			if handled == kind {
				found = true
				break
			}
		}
		if kind != FieldKindUnknown {
			assert.True(t, found, "FieldKind %s does not have a test for FieldType.GetFieldKind", kind)
		}
	}

}

func TestAccessor(t *testing.T) {

	text := &FieldTypeText{}
	reference := &FieldTypeReference{}
	form := &FieldTypeSubForm{}
	multilineText := &FieldTypeMultilineText{}
	date := &FieldTypeDate{}
	quantity := &FieldTypeQuantity{}
	month := &FieldTypeMonth{}
	week := &FieldTypeWeek{}
	singleSelect := &FieldTypeSingleSelect{}
	multiSelect := &FieldTypeMultiSelect{}
	boolean := &FieldTypeBoolean{}

	ft := FieldType{
		Text:          text,
		Reference:     reference,
		SubForm:       form,
		MultilineText: multilineText,
		Date:          date,
		Quantity:      quantity,
		Month:         month,
		Week:          week,
		SingleSelect:  singleSelect,
		MultiSelect:   multiSelect,
		Boolean:       boolean,
	}

	var foundValues []interface{}
	allKinds := GetAllFieldKinds()
	for _, kind := range allKinds {
		t.Run(kind.String(), func(t *testing.T) {
			k := kind
			field, err := ft.GetFieldType(k)
			if !assert.NoError(t, err) {
				return
			}
			if kind != FieldKindUnknown && !assert.NotNil(t, field, nilAccessorMSg, k) {
				return
			}
			for _, foundValue := range foundValues {
				if field == foundValue {
					assert.Fail(t, duplicateAccessorMsg, k)
				}
			}
			foundValues = append(foundValues, field)
		})
	}

	for kind := range fieldAccessors {
		found := false
		for _, registeredKind := range allKinds {
			if kind == registeredKind {
				found = true
				break
			}
		}
		if !assert.True(t, found, unregisteredAccessor, kind) {
			return
		}
	}

}
