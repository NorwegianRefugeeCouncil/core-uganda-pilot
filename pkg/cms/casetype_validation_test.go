package cms

import (
	"github.com/nrc-no/core/internal/validation"
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
				assert.NotNil(t, errList)
				assert.NotEmpty(t, errList)
				l := *errList.Find(".name")
				err := l[0]
				assert.Equal(t, err.Type, validation.ErrorTypeRequired)
			},
		},
		{
			name:     "invalid name",
			caseType: &CaseType{Name: "%^&"},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotNil(t, errList)
				assert.NotEmpty(t, errList)
				l := *errList.Find(".name")
				err := l[0]
				assert.Equal(t, err.Type, validation.ErrorTypeInvalid)
			},
		},
		{
			name:     "valid name",
			caseType: &CaseType{Name: "Test(54) one_two"},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Empty(t, errList.Find(".name"))
			},
		},
		{
			name:     "empty partyTypeId",
			caseType: &CaseType{},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotNil(t, errList)
				assert.NotEmpty(t, errList)
				l := *errList.Find(".name")
				err := l[0]
				assert.Equal(t, err.Type, validation.ErrorTypeRequired)
			},
		},
		{
			name:     "empty teamId",
			caseType: &CaseType{},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotNil(t, errList)
				assert.NotEmpty(t, errList)
				l := *errList.Find(".name")
				err := l[0]
				assert.Equal(t, err.Type, validation.ErrorTypeRequired)
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
