package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidatePartyType(t *testing.T) {
	tcs := []struct {
		name      string
		partyType *PartyType
		assert    func(t *testing.T, errList validation.ErrorList)
	}{
		{
			name:      "empty name",
			partyType: &PartyType{},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".name")[0].Type, validation.ErrorTypeRequired)
			},
		},
		{
			name:      "invalid name",
			partyType: &PartyType{Name: "*(^"},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".name")[0].Type, validation.ErrorTypeInvalid)
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			errList := ValidatePartyType(tc.partyType, validation.NewPath(""))
			tc.assert(t, errList)
		})
	}
}
