package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const duplicateAccessorMsg = `it seems that the accessor for %s returned a value that was returned by another accessor.
Check that the accessor returns the proper value`

const nilAccessorMSg = `The accessor for FieldKind %s returned a nil value. That should not happen`

func TestAccessor(t *testing.T) {

	text := &FieldTypeText{}
	reference := &FieldTypeReference{}
	form := &FieldTypeSubForm{}
	multilineText := &FieldTypeMultilineText{}
	date := &FieldTypeDate{}
	quantity := &FieldTypeQuantity{}
	month := &FieldTypeMonth{}
	singleSelect := &FieldTypeSingleSelect{}

	ft := FieldType{
		Text:          text,
		Reference:     reference,
		SubForm:       form,
		MultilineText: multilineText,
		Date:          date,
		Quantity:      quantity,
		Month:         month,
		SingleSelect:  singleSelect,
	}

	var foundValues []interface{}
	for _, kind := range GetAllFieldKinds() {
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
}
