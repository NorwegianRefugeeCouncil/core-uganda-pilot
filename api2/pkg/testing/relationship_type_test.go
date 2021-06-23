package testing

import (
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

// add a test to test the filter of relationship types, refer to case_test.go as an example

func (s *Suite) TestRelationShipTypeCRUD() {
	// CREATE relationship type
	mock := "create"
	created, err := s.server.RelationshipTypeClient.Create(s.ctx, &relationshiptypes.RelationshipType{
		IsDirectional:   true,
		Name:            mock,
		FirstPartyRole:  mock,
		SecondPartyRole: mock,
		Rules:           []relationshiptypes.RelationshipTypeRule{{&relationshiptypes.PartyTypeRule{FirstPartyTypeID: mock, SecondPartyTypeID: mock}}},
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.NotEmpty(s.T(), created.ID)
	assert.True(s.T(), created.IsDirectional)
	assert.Equal(s.T(), mock, created.Name)
	assert.Equal(s.T(), mock, created.FirstPartyRole)
	assert.Equal(s.T(), mock, created.SecondPartyRole)
	assert.IsType(s.T(), []relationshiptypes.RelationshipTypeRule{{&relationshiptypes.PartyTypeRule{FirstPartyTypeID: mock, SecondPartyTypeID: mock}}}, created.Rules)

	// GET relationship type
	get, err := s.server.RelationshipTypeClient.Get(s.ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Equal(s.T(), get, created) {
		return
	}

	// UPDATE relationship type
	updatedMock := "update"

	updated, err := s.server.RelationshipTypeClient.Update(s.ctx, &relationshiptypes.RelationshipType{
		ID:              created.ID,
		IsDirectional:   !created.IsDirectional,
		Name:            updatedMock,
		FirstPartyRole:  updatedMock,
		SecondPartyRole: updatedMock,
		Rules: []relationshiptypes.RelationshipTypeRule{
			{
				&relationshiptypes.PartyTypeRule{
					FirstPartyTypeID:  updatedMock,
					SecondPartyTypeID: updatedMock,
				},
			},
		},
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), created.ID, updated.ID)
	assert.Equal(s.T(), updatedMock, updated.Name)
	assert.False(s.T(), created.IsDirectional == updated.IsDirectional)
	assert.Equal(s.T(), updatedMock, updated.FirstPartyRole)
	assert.Equal(s.T(), updatedMock, updated.SecondPartyRole)
	assert.IsType(s.T(), []relationshiptypes.RelationshipTypeRule{{&relationshiptypes.PartyTypeRule{FirstPartyTypeID: updatedMock, SecondPartyTypeID: updatedMock}}}, updated.Rules)

	// GET relationship type
	get, err = s.server.RelationshipTypeClient.Get(s.ctx, updated.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Equal(s.T(), get, updated) {
		return
	}

	// LIST relationship types
	list, err := s.server.RelationshipTypeClient.List(s.ctx, relationshiptypes.ListOptions{
		PartyType: updatedMock,
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, get)
}

// TestRelationshipTypeList tests that we can effectively filter relationship types by
// - PartyType = IndividualPartyType
// - PartyType = HouseholdPartyType
func (s *Suite) TestRelationshipTypeList() {

	s.T().Run("test filter by IndividualPartyType", func(t *testing.T) {
		t.Logf("listing relationship types with party type: IndividualPartyType")

		list, err := s.server.RelationshipTypeClient.List(s.ctx, relationshiptypes.ListOptions{
			PartyType: partytypes.IndividualPartyType.ID,
		})
		if !assert.NoError(t, err) {
			return
		}

		for _, rt := range list.Items {
			valid := false
			for _, r := range rt.Rules {
				t.Logf("checking rule for type %s \nwith 1st %s \nand 2nd %s \nto see if it contains %s", rt.Name, r.PartyTypeRule.FirstPartyTypeID, r.PartyTypeRule.SecondPartyTypeID, partytypes.IndividualPartyType.ID)
				if r.PartyTypeRule.FirstPartyTypeID == partytypes.IndividualPartyType.ID || r.PartyTypeRule.SecondPartyTypeID == partytypes.IndividualPartyType.ID {
					valid = true
				}
			}
			assert.True(t, valid, "asserting that there is at least one rule with the individual party type")
		}
	})

	s.T().Run("test filter by HouseholdPartyType", func(t *testing.T) {
		t.Logf("listing relationship types with party type: HouseholdPartyType")

		list2, err := s.server.RelationshipTypeClient.List(s.ctx, relationshiptypes.ListOptions{
			PartyType: partytypes.HouseholdPartyType.ID,
		})
		if !assert.NoError(t, err) {
			return
		}

		for _, rt := range list2.Items {
			valid := false
			for _, r := range rt.Rules {
				t.Logf("checking rule for type %s \nwith 1st %s \nand 2nd %s \nto see if it contains %s", rt.Name, r.PartyTypeRule.FirstPartyTypeID, r.PartyTypeRule.SecondPartyTypeID, partytypes.HouseholdPartyType.ID)
				if r.PartyTypeRule.FirstPartyTypeID == partytypes.HouseholdPartyType.ID || r.PartyTypeRule.SecondPartyTypeID == partytypes.HouseholdPartyType.ID {
					valid = true
				}
			}
			assert.True(t, valid, "asserting that there is at least one rule with the individual party type")
		}
	})

}
