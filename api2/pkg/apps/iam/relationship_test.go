package iam

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

func (s *Suite) TestRelationshipCRUD() {

	// CREATE relationship
	mock := "create"
	created, err := s.client.Relationships().Create(s.ctx, &Relationship{
		RelationshipTypeID: mock,
		FirstPartyID:       mock,
		SecondPartyID:      mock,
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.NotEmpty(s.T(), created.ID)
	assert.Equal(s.T(), mock, created.RelationshipTypeID)
	assert.Equal(s.T(), mock, created.FirstPartyID)
	assert.Equal(s.T(), mock, created.SecondPartyID)

	// GET relationship
	get, err := s.client.Relationships().Get(s.ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), created, get)

	// UPDATE relationships type
	updatedMock := "update"

	updated, err := s.client.Relationships().Update(s.ctx, &Relationship{
		ID:                 created.ID,
		RelationshipTypeID: updatedMock,
		FirstPartyID:       updatedMock,
		SecondPartyID:      updatedMock,
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), created.ID, updated.ID)
	assert.Equal(s.T(), updatedMock, updated.RelationshipTypeID)
	assert.Equal(s.T(), updatedMock, updated.FirstPartyID)
	assert.Equal(s.T(), updatedMock, updated.SecondPartyID)

	// GET relationships type
	get, err = s.client.Relationships().Get(s.ctx, updated.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Equal(s.T(), updated, get) {
		return
	}
	// LIST relationships types
	list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, get)
}

func (s *Suite) TestRelationshipListFilterByListOptions() {
	s.ResetDB()

	// generate mock data
	mockRelationshipTypeIDs := [3]string{"a", "b", "c"}

	var mockParties = make([]string, 0)
	numMockParties := 10
	for i := 0; i < numMockParties; i++ {
		mockParties = append(mockParties, strconv.Itoa(i))
	}

	type MockRelationships struct {
		RelationshipTypeID string
		FirstParty         string
		SecondParty        string
	}

	var mockRelationships []MockRelationships
	numMockRelationships := 20

	for i := 0; i < numMockRelationships; i++ {
		p1 := mockParties[rand.Intn(len(mockParties))]
		p2 := p1
		for p2 == p1 {
			p2 = mockParties[rand.Intn(len(mockParties))]
		}
		mockRelationships = append(mockRelationships, MockRelationships{
			mockRelationshipTypeIDs[rand.Intn(len(mockRelationshipTypeIDs))],
			p1,
			p2,
		})
	}

	// maps of fields to relationships
	var byRelationshipTypeID = map[string][]*Relationship{}
	var byFirstParty = map[string][]*Relationship{}
	var bySecondParty = map[string][]*Relationship{}
	var byEitherParty = map[string][]*Relationship{}
	var byParties = map[string][]*Relationship{}
	var byRFS = map[string][]*Relationship{}
	var byRTIDAndEitherParty = map[string][]*Relationship{}

	// Create the relationships
	for _, mock := range mockRelationships {
		r, err := s.client.Relationships().Create(s.ctx, &Relationship{
			RelationshipTypeID: mock.RelationshipTypeID,
			FirstPartyID:       mock.FirstParty,
			SecondPartyID:      mock.SecondParty,
		})
		if !assert.NoError(s.T(), err) {
			return
		}
		byRelationshipTypeID[mock.RelationshipTypeID] = append(byRelationshipTypeID[mock.RelationshipTypeID], r)
		byFirstParty[mock.FirstParty] = append(byFirstParty[mock.FirstParty], r)
		bySecondParty[mock.SecondParty] = append(bySecondParty[mock.SecondParty], r)
		byEitherParty[mock.FirstParty] = append(byEitherParty[mock.FirstParty], r)
		byEitherParty[mock.SecondParty] = append(byEitherParty[mock.SecondParty], r)
		byRTIDAndEitherParty[mock.RelationshipTypeID+mock.FirstParty] = append(byEitherParty[mock.RelationshipTypeID+mock.FirstParty], r)
		byRTIDAndEitherParty[mock.RelationshipTypeID+mock.SecondParty] = append(byEitherParty[mock.RelationshipTypeID+mock.SecondParty], r)
		firstAndSecond := mock.FirstParty + mock.SecondParty
		byParties[firstAndSecond] = append(byParties[firstAndSecond], r)
		rfs := mock.RelationshipTypeID + firstAndSecond
		byRFS[rfs] = append(byRFS[rfs], r)
	}

	// Ensure we created enough entries
	list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}
	if l := len(list.Items); l != numMockRelationships {
		s.T().Logf("Incorrect number of DB entries created. expected: %d actual: %d", numMockRelationships, l)
		s.T().Fatal()
	}

	// Actual tests
	s.T().Run("Filter by RelationshipTypeID", func(t *testing.T) {
		for _, mockRTID := range mockRelationshipTypeIDs {
			list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
				RelationshipTypeID: mockRTID,
			})
			if !assert.NoError(t, err) {
				return
			}
			// map entry and returned list should have the same length
			assert.Len(t, byRelationshipTypeID[mockRTID], len(list.Items))
			// map entry should contain same items as list
			for _, l := range list.Items {
				assert.Contains(t, byRelationshipTypeID[mockRTID], l)
			}
		}
	})

	s.T().Run("Filter by FirstPartyID", func(t *testing.T) {
		for _, mockParty := range mockParties {
			list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
				FirstPartyID: mockParty,
			})
			if !assert.NoError(t, err) {
				return
			}
			// map entry and returned list should have the same length
			assert.Len(t, byFirstParty[mockParty], len(list.Items))
			// map entry should contain same items as list
			for _, l := range list.Items {
				assert.Contains(t, byFirstParty[mockParty], l)
			}
		}
	})

	s.T().Run("Filter by SecondPartyID", func(t *testing.T) {
		for _, mockParty := range mockParties {
			list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
				SecondPartyID: mockParty,
			})
			if !assert.NoError(t, err) {
				return
			}
			// map entry and returned list should have the same length
			assert.Len(t, bySecondParty[mockParty], len(list.Items))
			// map entry should contain same items as list
			for _, l := range list.Items {
				assert.Contains(t, bySecondParty[mockParty], l)
			}
		}
	})

	s.T().Run("Filter by EitherParty", func(t *testing.T) {
		for _, mockParty := range mockParties {
			list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
				EitherPartyID: mockParty,
			})
			if !assert.NoError(t, err) {
				return
			}
			// map entry and returned list should have the same length
			assert.Len(t, byEitherParty[mockParty], len(list.Items))
			// map entry should contain same items as list
			for _, l := range list.Items {
				assert.Contains(t, byEitherParty[mockParty], l)
			}
		}
	})

	s.T().Run("Filter by RTID and EitherParty", func(t *testing.T) {
		for _, mockRTID := range mockRelationshipTypeIDs {
			for _, mockParty := range mockParties {
				list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
					RelationshipTypeID: mockRTID,
					EitherPartyID:      mockParty,
				})
				if !assert.NoError(t, err) {
					return
				}
				// map entry and returned list should have the same length
				assert.Len(t, byRTIDAndEitherParty[mockRTID+mockParty], len(list.Items))
				// map entry should contain same items as list
				for _, l := range list.Items {
					assert.Contains(t, byRTIDAndEitherParty[mockRTID+mockParty], l)
				}
			}
		}
	})

	s.T().Run("Filter by combinations of First- and SecondPartyID", func(t *testing.T) {
		for _, mockParty1 := range mockParties {
			for _, mockParty2 := range mockParties {
				if mockParty1 == mockParty2 {
					continue
				}
				list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
					FirstPartyID:  mockParty1,
					SecondPartyID: mockParty2,
				})
				if !assert.NoError(t, err) {
					return
				}
				key := mockParty1 + mockParty2
				// map entry and returned list should have the same length
				assert.Len(t, byParties[key], len(list.Items))
				// map entry should contain same items as list
				for _, l := range list.Items {
					assert.Contains(t, byParties[key], l)
				}
			}
		}
	})

	s.T().Run("Filter by combinations of RTID, First- and SecondPartyID", func(t *testing.T) {
		for _, mockRTID := range mockRelationshipTypeIDs {
			for _, mockParty1 := range mockParties {
				for _, mockParty2 := range mockParties {
					if mockParty1 == mockParty2 {
						continue
					}
					list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
						RelationshipTypeID: mockRTID,
						FirstPartyID:       mockParty1,
						SecondPartyID:      mockParty2,
					})
					if !assert.NoError(t, err) {
						return
					}
					key := mockRTID + mockParty1 + mockParty2
					// map entry and returned list should have the same length
					assert.Len(t, byRFS[key], len(list.Items))
					// map entry should contain same items as list
					for _, l := range list.Items {
						assert.Contains(t, byRFS[key], l)
					}
				}
			}
		}
	})
}
