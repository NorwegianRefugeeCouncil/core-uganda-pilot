package iam

import (
	"github.com/nrc-no/core/internal/validation"
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
			assert:    assertRequired(".name"),
		},
		{
			name:      "invalid name",
			partyType: &PartyType{Name: "*(^"},
			assert:    assertInvalid(".name"),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			errList := ValidatePartyType(tc.partyType, validation.NewPath(""))
			tc.assert(t, errList)
		})
	}
}
