package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	"github.com/stretchr/testify/assert"
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
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".name")[0].Type, validation.ErrorTypeRequired)
			},
		},
		{
			name:             "invalid name",
			relationshipType: &RelationshipType{Name: "^&"},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".name")[0].Type, validation.ErrorTypeInvalid)
			},
		},
		{
			name:             "empty firstPartyRole",
			relationshipType: &RelationshipType{},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".firstPartyRole")[0].Type, validation.ErrorTypeRequired)
			},
		},
		{
			name:             "invalid firstPartyRole",
			relationshipType: &RelationshipType{FirstPartyRole: "^&"},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".firstPartyRole")[0].Type, validation.ErrorTypeInvalid)
			},
		},
		{
			name:             "empty secondPartyRole",
			relationshipType: &RelationshipType{IsDirectional: true},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".secondPartyRole")[0].Type, validation.ErrorTypeRequired)
			},
		},
		{
			name:             "invalid secondPartyRole",
			relationshipType: &RelationshipType{IsDirectional: true, SecondPartyRole: "^&"},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".secondPartyRole")[0].Type, validation.ErrorTypeInvalid)
			},
		},
		{
			name:             "empty rules",
			relationshipType: &RelationshipType{},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".rules")[0].Type, validation.ErrorTypeRequired)
			},
		},
		{
			name: "empty secondPartyTypeId",
			relationshipType: &RelationshipType{Rules: []RelationshipTypeRule{{
				PartyTypeRule: &PartyTypeRule{},
			}}},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".rules[0].firstPartyTypeId")[0].Type, validation.ErrorTypeRequired)
			},
		},
		{
			name: "empty secondPartyTypeId",
			relationshipType: &RelationshipType{IsDirectional: true, Rules: []RelationshipTypeRule{{
				PartyTypeRule: &PartyTypeRule{},
			}}},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".rules[0].secondPartyTypeId")[0].Type, validation.ErrorTypeRequired)
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			errList := ValidateRelationshipType(tc.relationshipType, validation.NewPath(""))
			tc.assert(t, errList)
		})
	}
}
