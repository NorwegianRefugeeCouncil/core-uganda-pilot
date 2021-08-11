package iam

import (
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/validation"
	"testing"
)

func TestValidateAttribute(t *testing.T) {
	tcs := []struct {
		name      string
		attribute *Attribute
		assert    func(t *testing.T, errList validation.ErrorList)
	}{
		{
			name:      "empty name",
			attribute: &Attribute{},
			assert:    assertRequired(".name"),
		},
		{
			name:      "invalid name",
			attribute: &Attribute{Name: "&2"},
			assert:    assertInvalid(".name"),
		},
		{
			name:      "empty party types",
			attribute: &Attribute{},
			assert:    assertRequired(".partyTypeIds"),
		},
		{
			name:      "empty translations",
			attribute: &Attribute{},
			assert:    assertRequired(".translations"),
		},
		{
			name: "missing locale",
			attribute: &Attribute{
				Translations: []AttributeTranslation{{}},
			},
			assert: assertRequired(".translations[0].locale"),
		},
		{
			name: "invalid locale",
			attribute: &Attribute{
				Translations: []AttributeTranslation{{Locale: "*"}},
			},
			assert: assertInvalid(".translations[0].locale"),
		},
		{
			name: "missing long formulation",
			attribute: &Attribute{
				Translations: []AttributeTranslation{{}},
			},
			assert: assertRequired(".translations[0].long"),
		},
		{
			name: "missing short formulation",
			attribute: &Attribute{
				Translations: []AttributeTranslation{{}},
			},
			assert: assertRequired(".translations[0].short"),
		},
		{
			name: "missing required field",
			attribute: &Attribute{
				Type: form.Text,
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			errList := ValidateAttribute(tc.attribute, validation.NewPath(""))
			tc.assert(t, errList)
		})
	}
}
