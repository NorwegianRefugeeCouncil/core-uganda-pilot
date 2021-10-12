package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	"testing"
)

func TestValidateAttribute(t *testing.T) {
	tcs := []struct {
		name      string
		attribute *PartyAttributeDefinition
		assert    func(t *testing.T, errList validation.ErrorList)
	}{}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			errList := ValidatePartyAttributeDefinition(tc.attribute, validation.NewPath(""))
			tc.assert(t, errList)
		})
	}
}
