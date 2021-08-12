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
			Template: &CaseTemplate{
				FormElements: []FormElement{
					{
						Type: Dropdown,
						Attributes: FormElementAttribute{
							Name: "dropdown",
						},
						Validation: FormElementValidation{
							Required: true,
						},
					},
					{
						Type: Textarea,
						Attributes: FormElementAttribute{
							Name: "textarea",
						},
						Validation: FormElementValidation{
							Required: true,
						},
					},
					{
						Type: TextInput,
						Attributes: FormElementAttribute{
							Name: "textinput",
						},
						Validation: FormElementValidation{
							Required: true,
						},
					},
					{
						Type: Checkbox,
						Attributes: FormElementAttribute{
							Name: "checkbox",
						},
						Validation: FormElementValidation{
							Required: true,
						},
					},
				},
			},
		},
		assert: func(t *testing.T, errList validation.ErrorList) {
			assert.NotEmpty(t, errList)
			dd := errList.Find(".dropdown")
			ta := errList.Find(".textarea")
			ti := errList.Find(".textinput")
			cb := errList.Find(".checkbox")
			for _, list := range []*validation.ErrorList{dd, ta, ti, cb} {
				assert.NotNil(t, list)
				assert.NotEmpty(t, list)
				assert.Len(t, *list, 1)
				l := *list
				err := l[0]
				assert.Equal(t, err.Type, validation.ErrorTypeRequired)
			}
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
