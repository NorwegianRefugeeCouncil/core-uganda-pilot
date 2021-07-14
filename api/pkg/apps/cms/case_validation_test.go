package cms

import (
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
			FormData: &CaseTemplate{
				FormElements: []CaseTemplateFormElement{
					{
						Type: Dropdown,
						Attributes: CaseTemplateFormElementAttribute{
							ID: "dropdown",
						},
						Validation: CaseTemplateFormElementValidation{
							Required: true,
						},
					},
					{
						Type: Textarea,
						Attributes: CaseTemplateFormElementAttribute{
							ID: "textarea",
						},
						Validation: CaseTemplateFormElementValidation{
							Required: true,
						},
					},
					{
						Type: TextInput,
						Attributes: CaseTemplateFormElementAttribute{
							ID: "textinput",
						},
						Validation: CaseTemplateFormElementValidation{
							Required: true,
						},
					},
					{
						Type: Checkbox,
						Attributes: CaseTemplateFormElementAttribute{
							ID: "checkbox",
						},
						Validation: CaseTemplateFormElementValidation{
							Required: true,
						},
					},
				},
			},
		},
		assert: func(t *testing.T, errList validation.ErrorList) {
			assert.NotEmpty(t, errList)
			assert.Equal(t, errList.Find("dropdown")[0].Type, validation.ErrorTypeRequired)
			assert.Equal(t, errList.Find("textarea")[0].Type, validation.ErrorTypeRequired)
			assert.Equal(t, errList.Find("textinput")[0].Type, validation.ErrorTypeRequired)
			assert.Equal(t, errList.Find("checkbox")[0].Type, validation.ErrorTypeRequired)
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
