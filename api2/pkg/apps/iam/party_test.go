package iam

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func (s *Suite) TestPartyCRUD() {

	// CREATE
	mock := "create"
	created, err := s.client.Parties().Create(s.ctx, &Party{
		PartyTypeIDs: []string{mock},
		Attributes:   PartyAttributes{"mock": []string{mock}},
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), []string{mock}, created.PartyTypeIDs)
	assert.Equal(s.T(), PartyAttributes{"mock": []string{mock}}, created.Attributes)

	// GET
	get, err := s.client.Parties().Get(s.ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created, get)

	// UPDATE
	updatedMock := "update"
	updated, err := s.client.Parties().Update(s.ctx, &Party{
		ID:           created.ID,
		PartyTypeIDs: []string{updatedMock},
		Attributes:   PartyAttributes{"mock": []string{updatedMock}},
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created.ID, updated.ID)
	assert.Equal(s.T(), []string{updatedMock}, updated.PartyTypeIDs)
	assert.Equal(s.T(), PartyAttributes{"mock": []string{updatedMock}}, updated.Attributes)

	// GET
	get, err = s.client.Parties().Get(s.ctx, updated.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.client.Parties().List(s.ctx, PartyListOptions{})
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), list.Items, get)
}

// TestPartyList tests that we can effectively filter the parties by
// - PartyID
// - PartyTypeID
// - Both PartyID and PartyTypeID
func (s *Suite) TestPartyList() {
	partyTypeID1 := uuid.NewV4().String()
	partyTypeID2 := uuid.NewV4().String()
	partyTypeID3 := uuid.NewV4().String()
	partyID1 := uuid.NewV4().String()
	partyID2 := uuid.NewV4().String()

	// Define the parties we're going to create/list
	partyList := []*struct {
		partyTypeID string
		partyID     string
		party       *Party
	}{
		{
			partyTypeID: partyTypeID1,
			partyID:     partyID1,
		}, {
			partyTypeID: partyTypeID2,
			partyID:     partyID2,
		}, {
			partyTypeID: partyTypeID3,
			partyID:     partyID1,
		}, {
			partyTypeID: partyTypeID3,
			partyID:     partyID2,
		},
	}

	// Holds the parties by PartyTypeID
	byPartyTypeID := map[string][]*Party{}

	// Holds the parties by PartyID
	byPartyID := map[string][]*Party{}

	// Holds the parties by PartyTypeID + PartyID
	byPartyTypeIdAndPartyId := map[string][]*Party{}

	// Create the parties
	for _, party := range partyList {
		k, err := s.client.Parties().Create(s.ctx, &Party{
			PartyTypeIDs: []string{party.partyTypeID},
		})
		if !assert.NoError(s.T(), err) {
			return
		}
		party.party = k

		// Add the created parties to the relevant maps
		byPartyTypeID[party.partyTypeID] = append(byPartyTypeID[party.partyTypeID], k)
		byPartyID[party.partyID] = append(byPartyID[party.partyID], k)
		caseAndPartyId := party.partyTypeID + "-" + party.partyID
		byPartyTypeIdAndPartyId[caseAndPartyId] = append(byPartyTypeIdAndPartyId[caseAndPartyId], k)
	}

	s.T().Run("test filter by partyTypeId", func(t *testing.T) {
		for _, party := range partyList {

			t.Logf("listing parties with PartyTypeID: %s", party.partyTypeID)

			list, err := s.client.Parties().List(s.ctx, PartyListOptions{
				PartyTypeID: party.partyTypeID,
			})
			if !assert.NoError(t, err) {
				return
			}

			t.Logf("received %d items", len(list.Items))
			assert.Len(t, list.Items, len(byPartyTypeID[party.partyTypeID]))
			for _, k := range byPartyTypeID[party.partyTypeID] {
				assert.Contains(t, list.Items, k)
			}
		}
	})

}
