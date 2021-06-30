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

	type MockRelationship struct {
		RelationshipTypeID string
		FirstParty         string
		SecondParty        string
	}

	var mockRelationships []*MockRelationship
	numMockRelationships := 20

	for i := 0; i < numMockRelationships; i++ {
		p1 := mockParties[rand.Intn(len(mockParties))]
		p2 := p1
		for p2 == p1 {
			p2 = mockParties[rand.Intn(len(mockParties))]
		}
		mockRelationships = append(mockRelationships, &MockRelationship{
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
	var byRelTypeIDFirstAndSecondParties = map[string][]*Relationship{}
	var byRelTypeIDAndEitherParty = map[string][]*Relationship{}

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
		byRelTypeIDAndEitherParty[mock.RelationshipTypeID+mock.FirstParty] = append(byRelTypeIDAndEitherParty[mock.RelationshipTypeID+mock.FirstParty], r)
		byRelTypeIDAndEitherParty[mock.RelationshipTypeID+mock.SecondParty] = append(byRelTypeIDAndEitherParty[mock.RelationshipTypeID+mock.SecondParty], r)
		firstAndSecond := mock.FirstParty + mock.SecondParty
		byParties[firstAndSecond] = append(byParties[firstAndSecond], r)
		rfs := mock.RelationshipTypeID + firstAndSecond
		byRelTypeIDFirstAndSecondParties[rfs] = append(byRelTypeIDFirstAndSecondParties[rfs], r)
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
		for _, mockRelTypeID := range mockRelationshipTypeIDs {
			list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
				RelationshipTypeID: mockRelTypeID,
			})
			if !assert.NoError(t, err) {
				return
			}
			// map entry and returned list should have the same length
			assert.Len(t, byRelationshipTypeID[mockRelTypeID], len(list.Items))
			// map entry should contain same items as list
			for _, l := range list.Items {
				assert.Contains(t, byRelationshipTypeID[mockRelTypeID], l)
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

	s.T().Run("Filter by RelTypeID and EitherParty", func(t *testing.T) {
		for _, mockRelTypeID := range mockRelationshipTypeIDs {
			for _, mockParty := range mockParties {
				list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
					RelationshipTypeID: mockRelTypeID,
					EitherPartyID:      mockParty,
				})
				if !assert.NoError(t, err) {
					return
				}
				// map entry and returned list should have the same length
				assert.Len(t, byRelTypeIDAndEitherParty[mockRelTypeID+mockParty], len(list.Items))
				// map entry should contain same items as list
				for _, l := range list.Items {
					assert.Contains(t, byRelTypeIDAndEitherParty[mockRelTypeID+mockParty], l)
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

	s.T().Run("Filter by combinations of RelTypeID, First- and SecondPartyID", func(t *testing.T) {
		for _, mockRelTypeID := range mockRelationshipTypeIDs {
			for _, mockParty1 := range mockParties {
				for _, mockParty2 := range mockParties {
					if mockParty1 == mockParty2 {
						continue
					}
					list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
						RelationshipTypeID: mockRelTypeID,
						FirstPartyID:       mockParty1,
						SecondPartyID:      mockParty2,
					})
					if !assert.NoError(t, err) {
						return
					}
					key := mockRelTypeID + mockParty1 + mockParty2
					// map entry and returned list should have the same length
					assert.Len(t, byRelTypeIDFirstAndSecondParties[key], len(list.Items))
					// map entry should contain same items as list
					for _, l := range list.Items {
						assert.Contains(t, byRelTypeIDFirstAndSecondParties[key], l)
					}
				}
			}
		}
	})
}
