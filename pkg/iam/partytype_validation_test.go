package iam

import (
	"github.com/nrc-no/core/internal/validation"
	"testing"
)

func TestValidateParty(t *testing.T) {
	tcs := []struct {
		name   string
		party  *Party
		assert func(t *testing.T, errList validation.ErrorList)
	}{
		{
			name:   "empty partyTypeIds",
			party:  &Party{},
			assert: assertRequired(".partyTypeIds"),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			errList := ValidateParty(tc.party, validation.NewPath(""))
			tc.assert(t, errList)
		})
	}
}
