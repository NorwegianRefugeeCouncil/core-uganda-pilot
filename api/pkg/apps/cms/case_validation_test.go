package cms

import (
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/validation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateCase(t *testing.T) {

	tcs := []struct {
		name     string
		caseType *Case
		assert   func(t *testing.T, errList validation.ErrorList)
	}{{
		name: "emptyFields",
		caseType: &Case{
			Template: &CaseTemplate{
				FormElements: []form.FormElement{
					{
						Type: form.Dropdown,
						Attributes: form.FormElementAttribute{
							Name: "dropdown",
						},
						Validation: form.FormElementValidation{
							Required: true,
						},
					},
					{
						Type: form.Textarea,
						Attributes: form.FormElementAttribute{
							Name: "textarea",
						},
						Validation: form.FormElementValidation{
							Required: true,
						},
					},
					{
						Type: form.TextInput,
						Attributes: form.FormElementAttribute{
							Name: "textinput",
						},
						Validation: form.FormElementValidation{
							Required: true,
						},
					},
					{
						Type: form.Checkbox,
						Attributes: form.FormElementAttribute{
							Name: "checkbox",
						},
						Validation: form.FormElementValidation{
							Required: true,
						},
					},
				},
			},
		},
		assert: func(t *testing.T, errList validation.ErrorList) {
			assert.NotEmpty(t, errList)
			assert.Equal(t, errList.Find(".dropdown")[0].Type, validation.ErrorTypeRequired)
			assert.Equal(t, errList.Find(".textarea")[0].Type, validation.ErrorTypeRequired)
			assert.Equal(t, errList.Find(".textinput")[0].Type, validation.ErrorTypeRequired)
			assert.Equal(t, errList.Find(".checkbox")[0].Type, validation.ErrorTypeRequired)
		},
	}}

	for _, tc := range tcs {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			errList := ValidateCase(testCase.caseType, validation.NewPath(""))
			testCase.assert(t, errList)
		})
	}

}
