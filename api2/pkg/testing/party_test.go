package testing

import (
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func (s *Suite) TestPartyCRUD() {
	// CREATE
	mock := "create"
	created, err := s.server.PartyClient.Create(s.ctx, &parties.Party{
		PartyTypeIDs: []string{mock},
		Attributes:   parties.PartyAttributes{"mock": []string{mock}},
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), []string{mock}, created.PartyTypeIDs)
	assert.Equal(s.T(), parties.PartyAttributes{"mock": []string{mock}}, created.Attributes)

	// GET
	get, err := s.server.PartyClient.Get(s.ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created, get)

	// UPDATE
	updatedMock := "update"
	updated, err := s.server.PartyClient.Update(s.ctx, &parties.Party{
		ID:           created.ID,
		PartyTypeIDs: []string{updatedMock},
		Attributes:   parties.PartyAttributes{"mock": []string{updatedMock}},
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created.ID, updated.ID)
	assert.Equal(s.T(), []string{updatedMock}, updated.PartyTypeIDs)
	assert.Equal(s.T(), parties.PartyAttributes{"mock": []string{updatedMock}}, updated.Attributes)

	// GET
	get, err = s.server.PartyClient.Get(s.ctx, updated.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.server.PartyClient.List(s.ctx, parties.ListOptions{})
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), list.Items, get)
}

// TestPartyList tests that we can effectively filter the parties by
// - PartyID
// - PartyTypeID
// - Both PartyID and PartyTypeID
func (s *Suite) TestPartyList() {
	s.T().SkipNow() // TODO
	partyTypeID1 := uuid.NewV4().String()
	partyTypeID2 := uuid.NewV4().String()
	partyTypeID3 := uuid.NewV4().String()
	partyID1 := uuid.NewV4().String()
	partyID2 := uuid.NewV4().String()

	// Define the parties we're going to create/list
	partyList := []*struct {
		partyTypeID string
		partyID     string
		party       *parties.Party
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
	byPartyTypeID := map[string][]*parties.Party{}
	// Holds the parties by PartyID
	byPartyID := map[string][]*parties.Party{}
	// Holds the parties by PartyTypeID + PartyID
	byPartyTypeIdAndPartyId := map[string][]*parties.Party{}

	// Create the parties
	for _, party := range partyList {
		k, err := s.server.PartyClient.Create(s.ctx, &parties.Party{})
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

			list, err := s.server.PartyClient.List(s.ctx, parties.ListOptions{})
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

	s.T().Run("test filter by partyId", func(t *testing.T) {
		for _, party := range partyList {

			t.Logf("listing parties with PartyID: %s", party.partyID)

			list, err := s.server.PartyClient.List(s.ctx, parties.ListOptions{})
			if !assert.NoError(t, err) {
				return
			}

			t.Logf("received %d items", len(list.Items))
			assert.Len(t, list.Items, len(byPartyID[party.partyID]))
			for _, k := range byPartyID[party.partyID] {
				assert.Contains(t, list.Items, k)
			}
		}
	})

	s.T().Run("test filter by partyTypeId and partyId", func(t *testing.T) {

		for _, party := range partyList {

			t.Logf("listing parties with PartyID: %s and PartyTypeID: %s", party.partyID, party.partyTypeID)

			list, err := s.server.PartyClient.List(s.ctx, parties.ListOptions{})
			if !assert.NoError(t, err) {
				return
			}

			t.Logf("received %d items", len(list.Items))
			caseAndPartyId := party.partyTypeID + "-" + party.partyID
			assert.Len(t, list.Items, len(byPartyTypeIdAndPartyId[caseAndPartyId]))
			for _, k := range byPartyID[caseAndPartyId] {
				assert.Contains(t, list.Items, k)
			}

		}

	})
}
