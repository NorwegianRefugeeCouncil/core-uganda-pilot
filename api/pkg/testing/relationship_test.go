package testing

import (
	"github.com/nrc-no/core/pkg/parties/relationships"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

func (s *Suite) TestRelationshipCRUD() {
	// CREATE relationship
	mock := "create"
	created, err := s.server.RelationshipClient.Create(s.ctx, &relationships.Relationship{
		RelationshipTypeID: mock,
		FirstParty:         mock,
		SecondParty:        mock,
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.NotEmpty(s.T(), created.ID)
	assert.Equal(s.T(), mock, created.RelationshipTypeID)
	assert.Equal(s.T(), mock, created.FirstParty)
	assert.Equal(s.T(), mock, created.SecondParty)

	// GET relationship
	get, err := s.server.RelationshipClient.Get(s.ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), created, get)

	// UPDATE relationships type
	updatedMock := "update"

	updated, err := s.server.RelationshipClient.Update(s.ctx, &relationships.Relationship{
		ID:                 created.ID,
		RelationshipTypeID: updatedMock,
		FirstParty:         updatedMock,
		SecondParty:        updatedMock,
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), created.ID, updated.ID)
	assert.Equal(s.T(), updatedMock, updated.RelationshipTypeID)
	assert.Equal(s.T(), updatedMock, updated.FirstParty)
	assert.Equal(s.T(), updatedMock, updated.SecondParty)

	// GET relationships type
	get, err = s.server.RelationshipClient.Get(s.ctx, updated.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Equal(s.T(), updated, get) {
		return
	}
	// LIST relationships types
	list, err := s.server.RelationshipClient.List(s.ctx, relationships.ListOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, get)
}

func (s *Suite) TestRelationshipListFilterByListOptions() {
	s.T().SkipNow()
	// TODO test filtering by start- and endOfRelationship values

	// generate mock data
	mockRelationshipTypeIDs := [3]string{"a", "b", "c"}

	var mockParties []string
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
	var byRelationshipTypeID map[string][]*relationships.Relationship
	var byFirstParty map[string][]*relationships.Relationship
	var bySecondParty map[string][]*relationships.Relationship
	var byEitherParty map[string][]*relationships.Relationship
	var byParties map[string][]*relationships.Relationship
	var byRFS map[string][]*relationships.Relationship
	var byRTIDAndEitherParty map[string][]*relationships.Relationship

	// Create the relationships
	for _, mock := range mockRelationships {
		r, err := s.server.RelationshipClient.Create(s.ctx, &relationships.Relationship{
			RelationshipTypeID: mock.RelationshipTypeID,
			FirstParty:         mock.FirstParty,
			SecondParty:        mock.SecondParty,
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
	list, err := s.server.RelationshipClient.List(s.ctx, relationships.ListOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}
	if l := len(list.Items); l != numMockRelationships {
		s.T().Logf("Incorrect number of DB entries created. expected: %d actual: %d", numMockRelationships, l)
		return
	}

	// Actual tests
	s.T().Run("Filter by RelationshipTypeID", func(t *testing.T) {
		for _, mockRTID := range mockRelationshipTypeIDs {
			list, err := s.server.RelationshipClient.List(s.ctx, relationships.ListOptions{
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
			list, err := s.server.RelationshipClient.List(s.ctx, relationships.ListOptions{
				FirstPartyId: mockParty,
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
			list, err := s.server.RelationshipClient.List(s.ctx, relationships.ListOptions{
				SecondParty: mockParty,
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
			list, err := s.server.RelationshipClient.List(s.ctx, relationships.ListOptions{
				EitherParty: mockParty,
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
				list, err := s.server.RelationshipClient.List(s.ctx, relationships.ListOptions{
					RelationshipTypeID: mockRTID,
					EitherParty:        mockParty,
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
				list, err := s.server.RelationshipClient.List(s.ctx, relationships.ListOptions{
					FirstPartyId: mockParty1,
					SecondParty:  mockParty2,
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
					list, err := s.server.RelationshipClient.List(s.ctx, relationships.ListOptions{
						RelationshipTypeID: mockRTID,
						FirstPartyId:       mockParty1,
						SecondParty:        mockParty2,
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
