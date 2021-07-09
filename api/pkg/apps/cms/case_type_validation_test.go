package cms

import (
	"github.com/nrc-no/core/pkg/validation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateCaseType(t *testing.T) {

	tcs := []struct {
		name     string
		caseType *CaseType
		assert   func(t *testing.T, errList validation.ErrorList)
	}{
		{
			name:     "empty name",
			caseType: &CaseType{},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Len(t, errList.FindErrors("name"), 1)
				assert.Equal(t, errList.FindErrors("name")[0].Type, validation.ErrorTypeRequired)
			},
		},
		{
			name:     "invalid name",
			caseType: &CaseType{Name: "%^&"},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Len(t, errList.FindErrors("name"), 1)
				assert.Equal(t, errList.FindErrors("name")[0].Type, validation.ErrorTypeInvalid)
			},
		},
		{
			name:     "empty partyTypeId",
			caseType: &CaseType{},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Len(t, errList.FindErrors("partyTypeId"), 1)
				assert.Equal(t, errList.FindErrors("partyTypeId")[0].Type, validation.ErrorTypeRequired)
			},
		},
		{
			name:     "empty teamId",
			caseType: &CaseType{},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Len(t, errList.FindErrors("partyTypeId"), 1)
				assert.Equal(t, errList.FindErrors("partyTypeId")[0].Type, validation.ErrorTypeRequired)
			},
		},
	}

	for _, tc := range tcs {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			errList := ValidateCaseType(testCase.caseType, validation.NewPath(""))
			testCase.assert(t, errList)
		})
	}

}
