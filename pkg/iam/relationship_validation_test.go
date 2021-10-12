package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	"testing"
)

func TestValidateRelationship(t *testing.T) {
	tcs := []struct {
		name         string
		relationship *Relationship
		assert       func(t *testing.T, errList validation.ErrorList)
	}{
		{
			name:         "empty relationshipType",
			relationship: &Relationship{},
			assert:       assertRequired(".relationshipTypeId"),
		}, {
			name:         "empty firstParty",
			relationship: &Relationship{},
			assert:       assertRequired(".firstParty"),
		}, {
			name:         "empty secondParty",
			relationship: &Relationship{},
			assert:       assertRequired(".secondParty"),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			errList := ValidateRelationship(tc.relationship, validation.NewPath(""))
			tc.assert(t, errList)
		})
	}
}
