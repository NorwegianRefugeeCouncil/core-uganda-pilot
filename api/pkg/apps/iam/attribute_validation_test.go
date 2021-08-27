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
		// Bypassed for now; see comment in attribute_validation.go
		//{
		//	name:      "empty name",
		//	attribute: &Attribute{},
		//	assert:    assertRequired(".name"),
		//},
		//{
		//	name:      "invalid name",
		//	attribute: &Attribute{Name: "&2"},
		//	assert:    assertInvalid(".name"),
		//},
		//{
		//	name:      "empty party types",
		//	attribute: &Attribute{},
		//	assert:    assertRequired(".partyTypeIds"),
		//},
		//{
		//	name:      "empty translations",
		//	attribute: &Attribute{},
		//	assert:    assertRequired(".translations"),
		//},
		//{
		//	name: "missing locale",
		//	attribute: &Attribute{
		//		Translations: []AttributeTranslation{{}},
		//	},
		//	assert: assertRequired(".translations[0].locale"),
		//},
		{
			name: "invalid locale",
			attribute: &Attribute{
				Translations: []AttributeTranslation{{Locale: "*"}},
			},
			assert: assertInvalid(".translations[0].locale"),
		},
		//{
		//	name: "missing long formulation",
		//	attribute: &Attribute{
		//		Translations: []AttributeTranslation{{}},
		//	},
		//	assert: assertRequired(".translations[0].long"),
		//},
		//{
		//	name: "missing short formulation",
		//	attribute: &Attribute{
		//		Translations: []AttributeTranslation{{}},
		//	},
		//	assert: assertRequired(".translations[0].short"),
		//},
		//{
		//	name: "missing required form field",
		//	attribute: &Attribute{
		//		Type: form.Text,
		//		Attributes: form.FormElementAttributes{
		//			Name: "text",
		//		},
		//		Validation: form.FormElementValidation{Required: true},
		//	},
		//	assert: assertRequired(".text"),
		//},
		{
			name: "invalid email",
			attribute: &Attribute{
				Type: form.Email,
				Attributes: form.FormElementAttributes{
					Name:  "field",
					Value: []string{"42"},
				},
			},
			assert: assertInvalid(".field"),
		},
		{
			name: "valid email",
			attribute: &Attribute{
				Type: form.Email,
				Attributes: form.FormElementAttributes{
					Name:  "field",
					Value: []string{"valid@email.com"},
				},
			},
			assert: assertNoError(".field"),
		},
		{
			name: "invalid phone",
			attribute: &Attribute{
				Type: form.Email,
				Attributes: form.FormElementAttributes{
					Name:  "field",
					Value: []string{"42"},
				},
			},
			assert: assertInvalid(".field"),
		},
		{
			name: "valid phone",
			attribute: &Attribute{
				Type: form.Phone,
				Attributes: form.FormElementAttributes{
					Name:  "field",
					Value: []string{"+256-345-939499"},
				},
			},
			assert: assertNoError(".field"),
		},
		{
			name: "valid phone alternate",
			attribute: &Attribute{
				Type: form.Phone,
				Attributes: form.FormElementAttributes{
					Name:  "field",
					Value: []string{"0345 939499"},
				},
			},
			assert: assertNoError(".field"),
		},
		{
			name: "invalid date",
			attribute: &Attribute{
				Type: form.Date,
				Attributes: form.FormElementAttributes{
					Name:  "field",
					Value: []string{"1987-14-02"},
				},
			},
			assert: assertInvalid(".field"),
		},
		{
			name: "invalid date alternate",
			attribute: &Attribute{
				Type: form.Date,
				Attributes: form.FormElementAttributes{
					Name:  "field",
					Value: []string{"1987-03-32"},
				},
			},
			assert: assertInvalid(".field"),
		},
		{
			name: "valid date",
			attribute: &Attribute{
				Type: form.Date,
				Attributes: form.FormElementAttributes{
					Name:  "field",
					Value: []string{"1987-12-30"},
				},
			},
			assert: assertNoError(".field"),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			errList := ValidateAttribute(tc.attribute, validation.NewPath(""))
			tc.assert(t, errList)
		})
	}
}
