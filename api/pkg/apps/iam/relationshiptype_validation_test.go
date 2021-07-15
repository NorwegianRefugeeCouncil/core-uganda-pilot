package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	"testing"
)

func TestValidateRelationshipType(t *testing.T) {
	tcs := []struct {
		name             string
		relationshipType *RelationshipType
		assert           func(t *testing.T, errList validation.ErrorList)
	}{
		{
			name:             "empty name",
			relationshipType: &RelationshipType{},
			assert:           assertRequired(".name"),
		},
		{
			name:             "invalid name",
			relationshipType: &RelationshipType{Name: "^&"},
			assert:           assertInvalid(".name"),
		},
		{
			name:             "empty firstPartyRole",
			relationshipType: &RelationshipType{},
			assert:           assertRequired(".firstPartyRole"),
		},
		{
			name:             "invalid firstPartyRole",
			relationshipType: &RelationshipType{FirstPartyRole: "^&"},
			assert:           assertInvalid(".firstPartyRole"),
		},
		{
			name:             "empty secondPartyRole",
			relationshipType: &RelationshipType{IsDirectional: true},
			assert:           assertRequired(".secondPartyRole"),
		},
		{
			name:             "invalid secondPartyRole",
			relationshipType: &RelationshipType{IsDirectional: true, SecondPartyRole: "^&"},
			assert:           assertInvalid(".secondPartyRole"),
		},
		{
			name:             "empty rules",
			relationshipType: &RelationshipType{},
			assert:           assertRequired(".rules"),
		},
		{
			name: "empty secondPartyTypeId",
			relationshipType: &RelationshipType{Rules: []RelationshipTypeRule{{
				PartyTypeRule: &PartyTypeRule{},
			}}},
			assert: assertRequired(".rules[0].firstPartyTypeId"),
		},
		{
			name: "empty secondPartyTypeId",
			relationshipType: &RelationshipType{IsDirectional: true, Rules: []RelationshipTypeRule{{
				PartyTypeRule: &PartyTypeRule{},
			}}},
			assert: assertRequired(".rules[0].secondPartyTypeId"),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			errList := ValidateRelationshipType(tc.relationshipType, validation.NewPath(""))
			tc.assert(t, errList)
		})
	}
}
