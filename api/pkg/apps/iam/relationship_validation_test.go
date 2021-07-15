package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	"github.com/stretchr/testify/assert"
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
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".relationshipTypeId")[0].Type, validation.ErrorTypeRequired)
			},
		}, {
			name:         "empty firstParty",
			relationship: &Relationship{},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".firstParty")[0].Type, validation.ErrorTypeRequired)
			},
		}, {
			name:         "empty secondParty",
			relationship: &Relationship{},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.NotEmpty(t, errList)
				assert.Equal(t, errList.Find(".secondParty")[0].Type, validation.ErrorTypeRequired)
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			errList := ValidateRelationship(tc.relationship, validation.NewPath(""))
			tc.assert(t, errList)
		})
	}
}
