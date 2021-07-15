package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateParty(t *testing.T) {
	tcs := []struct {
		name   string
		party  *Party
		assert func(t *testing.T, errList validation.ErrorList)
	}{
		{
			name:  "empty partyTypeIds",
			party: &Party{},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".partyTypeIds")[0].Type, validation.ErrorTypeRequired)
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			errList := ValidateParty(tc.party, validation.NewPath(""))
			tc.assert(t, errList)
		})
	}
}
