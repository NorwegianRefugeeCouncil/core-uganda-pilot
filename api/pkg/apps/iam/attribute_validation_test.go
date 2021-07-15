package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	"github.com/stretchr/testify/assert"
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
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".name")[0].Type, validation.ErrorTypeRequired)
			},
		},
		{
			name:      "invalid name",
			attribute: &Attribute{Name: "&2"},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".name")[0].Type, validation.ErrorTypeInvalid)
			},
		},
		{
			name:      "empty party types",
			attribute: &Attribute{},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".partyTypeIds")[0].Type, validation.ErrorTypeRequired)
			},
		},
		{
			name:      "empty translations",
			attribute: &Attribute{},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".translations")[0].Type, validation.ErrorTypeRequired)
			},
		},
		{
			name: "missing locale",
			attribute: &Attribute{
				Translations: []AttributeTranslation{{}},
			},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".translations[0].locale")[0].Type, validation.ErrorTypeRequired)
			},
		},
		{
			name: "invalid locale",
			attribute: &Attribute{
				Translations: []AttributeTranslation{{Locale: "*"}},
			},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".translations[0].locale")[0].Type, validation.ErrorTypeInvalid)
			},
		},
		{
			name: "missing long formulation",
			attribute: &Attribute{
				Translations: []AttributeTranslation{{}},
			},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".translations[0].long")[0].Type, validation.ErrorTypeRequired)
			},
		},
		{
			name: "missing short formulation",
			attribute: &Attribute{
				Translations: []AttributeTranslation{{}},
			},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".translations[0].short")[0].Type, validation.ErrorTypeRequired)
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
