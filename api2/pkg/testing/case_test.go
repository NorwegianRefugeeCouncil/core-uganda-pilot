package testing

import (
	"github.com/nrc-no/core-kafka/pkg/cases/cases"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

// TestCaseCRUD tests that we can do basic CRUD operation
// on the Cases API
func (s *Suite) TestCaseCRUD() {
	// CREATE
	mock := "create"
	created, err := s.server.CaseClient.Create(s.ctx, &cases.Case{
		CaseTypeID:  mock,
		PartyID:     mock,
		Description: mock,
		Done:        false,
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mock, created.CaseTypeID)
	assert.Equal(s.T(), mock, created.PartyID)
	assert.Equal(s.T(), mock, created.Description)
	assert.False(s.T(), created.Done)

	// GET
	get, err := s.server.CaseClient.Get(s.ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created, get)

	// UPDATE
	updatedMock := "update"
	updated, err := s.server.CaseClient.Update(s.ctx, &cases.Case{
		ID:          created.ID,
		CaseTypeID:  updatedMock,
		PartyID:     updatedMock,
		Description: updatedMock,
		Done:        !created.Done,
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created.ID, updated.ID)
	assert.Equal(s.T(), updatedMock, updated.CaseTypeID)
	assert.Equal(s.T(), updatedMock, updated.PartyID)
	assert.Equal(s.T(), updatedMock, updated.Description)
	assert.False(s.T(), created.Done == updated.Done)

	// GET
	get, err = s.server.CaseClient.Get(s.ctx, updated.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.server.CaseClient.List(s.ctx, cases.ListOptions{})
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), list.Items, get)
}

// TestCaseList tests that we can effectively filter the cases by
// - PartyID
// - CaseTypeID
// - Both PartyID and CaseTypeID
func (s *Suite) TestCaseList() {

	caseTypeID1 := uuid.NewV4().String()
	caseTypeID2 := uuid.NewV4().String()
	caseTypeID3 := uuid.NewV4().String()
	partyID1 := uuid.NewV4().String()
	partyID2 := uuid.NewV4().String()

	// Define the cases we're going to create/list
	kases := []*struct {
		caseTypeID string
		partyID    string
		kase       *cases.Case
	}{
		{
			caseTypeID: caseTypeID1,
			partyID:    partyID1,
		}, {
			caseTypeID: caseTypeID2,
			partyID:    partyID2,
		}, {
			caseTypeID: caseTypeID3,
			partyID:    partyID1,
		}, {
			caseTypeID: caseTypeID3,
			partyID:    partyID2,
		},
	}

	// Holds the cases by CaseTypeID
	byCaseTypeID := map[string][]*cases.Case{}
	// Holds the cases by PartyID
	byPartyID := map[string][]*cases.Case{}
	// Holds the cases by CaseTypeID + PartyID
	byCaseTypeIdAndPartyId := map[string][]*cases.Case{}

	// Create the cases
	for i, kase := range kases {
		k, err := s.server.CaseClient.Create(s.ctx, &cases.Case{
			CaseTypeID:  kase.caseTypeID,
			PartyID:     kase.partyID,
			Description: "TestCaseListByCaseTypeID-" + strconv.Itoa(i),
			Done:        false,
		})
		if !assert.NoError(s.T(), err) {
			return
		}
		kase.kase = k

		// Add the created cases to the relevant maps
		byCaseTypeID[kase.caseTypeID] = append(byCaseTypeID[kase.caseTypeID], k)
		byPartyID[kase.partyID] = append(byPartyID[kase.partyID], k)
		caseAndPartyId := kase.caseTypeID + "-" + kase.partyID
		byCaseTypeIdAndPartyId[caseAndPartyId] = append(byCaseTypeIdAndPartyId[caseAndPartyId], k)
	}

	s.T().Run("test filter by caseTypeId", func(t *testing.T) {
		for _, kase := range kases {

			t.Logf("listing cases with CaseTypeID: %s", kase.caseTypeID)

			list, err := s.server.CaseClient.List(s.ctx, cases.ListOptions{
				CaseTypeID: kase.caseTypeID,
			})
			if !assert.NoError(t, err) {
				return
			}

			t.Logf("received %d items", len(list.Items))
			assert.Len(t, list.Items, len(byCaseTypeID[kase.caseTypeID]))
			for _, k := range byCaseTypeID[kase.caseTypeID] {
				assert.Contains(t, list.Items, k)
			}
		}
	})

	s.T().Run("test filter by partyId", func(t *testing.T) {
		for _, kase := range kases {

			t.Logf("listing cases with PartyID: %s", kase.partyID)

			list, err := s.server.CaseClient.List(s.ctx, cases.ListOptions{
				PartyID: kase.partyID,
			})
			if !assert.NoError(t, err) {
				return
			}

			t.Logf("received %d items", len(list.Items))
			assert.Len(t, list.Items, len(byPartyID[kase.partyID]))
			for _, k := range byPartyID[kase.partyID] {
				assert.Contains(t, list.Items, k)
			}
		}
	})

	s.T().Run("test filter by caseTypeId and partyId", func(t *testing.T) {

		for _, kase := range kases {

			t.Logf("listing cases with PartyID: %s and CaseTypeID: %s", kase.partyID, kase.caseTypeID)

			list, err := s.server.CaseClient.List(s.ctx, cases.ListOptions{
				PartyID:    kase.partyID,
				CaseTypeID: kase.caseTypeID,
			})
			if !assert.NoError(t, err) {
				return
			}

			t.Logf("received %d items", len(list.Items))
			caseAndPartyId := kase.caseTypeID + "-" + kase.partyID
			assert.Len(t, list.Items, len(byCaseTypeIdAndPartyId[caseAndPartyId]))
			for _, k := range byPartyID[caseAndPartyId] {
				assert.Contains(t, list.Items, k)
			}

		}

	})
}
