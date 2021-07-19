// +build integration

package iam_test

import (
	. "github.com/nrc-no/core/pkg/apps/iam"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestRelationship() {
	s.Run("API", func() { s.testRelationshipAPI() })
	s.SetupTest()
	s.Run("List filtering", func() { s.testRelationshipListFilter() })
}

func (s *Suite) testRelationshipAPI() {
	// CREATE relationship
	relationship := s.mockRelationships(1)[0]
	relationship.RelationshipTypeID = newUUID()
	relationship.FirstPartyID = newUUID()
	relationship.SecondPartyID = newUUID()
	created, err := s.client.Relationships().Create(s.ctx, relationship)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), relationship.ID, created.ID)
	assert.Equal(s.T(), relationship.RelationshipTypeID, created.RelationshipTypeID)
	assert.Equal(s.T(), relationship.FirstPartyID, created.FirstPartyID)
	assert.Equal(s.T(), relationship.SecondPartyID, created.SecondPartyID)

	// GET relationship
	get, err := s.client.Relationships().Get(s.ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), created, get)

	// UPDATE relationships type
	relationship.RelationshipTypeID = newUUID()
	relationship.FirstPartyID = newUUID()
	relationship.SecondPartyID = newUUID()

	updated, err := s.client.Relationships().Update(s.ctx, relationship)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), relationship.ID, updated.ID)
	assert.Equal(s.T(), relationship.RelationshipTypeID, updated.RelationshipTypeID)
	assert.Equal(s.T(), relationship.FirstPartyID, updated.FirstPartyID)
	assert.Equal(s.T(), relationship.SecondPartyID, updated.SecondPartyID)

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

func (s *Suite) testRelationshipListFilter() {
	nRelationships := 10
	nRelationshipTypes := 3
	nParties := 4

	relationships := s.mockRelationships(nRelationships)
	relationshipTypes := s.mockRelationshipTypes(nRelationshipTypes)
	parties := s.mockParties(nParties)

	relationshipsFromTypes := make(map[string][]string)
	relationshipsFromFirstParties := make(map[string][]string)
	relationshipsFromSecondParties := make(map[string][]string)
	relationshipsFromEitherParties := make(map[string][]string)
	relationshipsFromBoth := make(map[string][]string)
	relationshipsFromTypeAndEither := make(map[string][]string)
	relationshipsFromTypeAndBoth := make(map[string][]string)

	// Prepare test data
	for i, relationship := range relationships {
		// Assign field values to the empty relationships
		rType := relationshipTypes[i%len(relationshipTypes)].ID
		party1 := parties[i%len(parties)].ID
		party2 := parties[(i+1)%len(parties)].ID
		relationship.RelationshipTypeID = rType
		relationship.FirstPartyID = party1
		relationship.SecondPartyID = party2

		// Save the relationship to the DB
		_, err := s.client.Relationships().Create(s.ctx, relationship)
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}

		// Maps the fields to the relationship
		id := relationship.ID
		relationshipsFromTypes[rType] = append(relationshipsFromTypes[rType], id)
		relationshipsFromFirstParties[party1] = append(relationshipsFromFirstParties[party1], id)
		relationshipsFromSecondParties[party2] = append(relationshipsFromSecondParties[party2], id)
		relationshipsFromEitherParties[party1] = append(relationshipsFromEitherParties[party1], id)
		relationshipsFromEitherParties[party2] = append(relationshipsFromEitherParties[party2], id)
		relationshipsFromBoth[party1+party2] = append(relationshipsFromBoth[party1+party2], id)
		relationshipsFromTypeAndEither[rType+party1] = append(relationshipsFromTypeAndEither[rType+party1], id)
		relationshipsFromTypeAndEither[rType+party2] = append(relationshipsFromTypeAndEither[rType+party2], id)
		relationshipsFromTypeAndBoth[rType+party1+party2] = append(relationshipsFromTypeAndBoth[rType+party1+party2], id)
	}

	s.Run("by type", func() { s.testRelationshipFilterByType(relationshipsFromTypes, relationshipTypes) })
	s.Run("by first party", func() { s.testRelationshipFilterByFirstParty(relationshipsFromFirstParties, parties) })
	s.Run("by second party", func() { s.testRelationshipFilterBySecondParty(relationshipsFromSecondParties, parties) })
	s.Run("by either party", func() { s.testRelationshipFilterByEitherParty(relationshipsFromEitherParties, parties) })
	s.Run("by both parties", func() { s.testRelationshipFilterByBothParties(relationshipsFromBoth, parties) })
	s.Run("by type and either party", func() {
		s.testRelationshipFilterByTypeAndEither(relationshipsFromTypeAndEither, relationshipTypes, parties)
	})
	s.Run("by type and both parties", func() {
		s.testRelationshipFilterByTypeAndBoth(relationshipsFromTypeAndBoth, relationshipTypes, parties)
	})
}

func (s *Suite) testRelationshipFilterByType(expected map[string][]string, types []*RelationshipType) {
	for _, relationshipType := range types {
		rType := relationshipType.ID
		list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
			RelationshipTypeID: rType,
		})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		// map entry and returned list should have the same length
		assert.Len(s.T(), expected[rType], len(list.Items))
		// map entry should contain same items as list
		for _, l := range list.Items {
			assert.Contains(s.T(), expected[rType], l.ID)
		}
	}
}

func (s *Suite) testRelationshipFilterByFirstParty(expected map[string][]string, parties []*Party) {
	for _, party := range parties {
		list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
			FirstPartyID: party.ID,
		})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		// map entry and returned list should have the same length
		assert.Len(s.T(), expected[party.ID], len(list.Items))
		// map entry should contain same items as list
		for _, l := range list.Items {
			assert.Contains(s.T(), expected[party.ID], l.ID)
		}
	}
}

func (s *Suite) testRelationshipFilterBySecondParty(expected map[string][]string, parties []*Party) {
	for _, party := range parties {
		list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
			SecondPartyID: party.ID,
		})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		// map entry and returned list should have the same length
		assert.Len(s.T(), expected[party.ID], len(list.Items))
		// map entry should contain same items as list
		for _, l := range list.Items {
			assert.Contains(s.T(), expected[party.ID], l.ID)
		}
	}
}

func (s *Suite) testRelationshipFilterByEitherParty(expected map[string][]string, parties []*Party) {
	for _, party := range parties {
		list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
			EitherPartyID: party.ID,
		})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		// map entry and returned list should have the same length
		assert.Len(s.T(), expected[party.ID], len(list.Items))
		// map entry should contain same items as list
		for _, l := range list.Items {
			assert.Contains(s.T(), expected[party.ID], l.ID)
		}
	}
}

func (s *Suite) testRelationshipFilterByBothParties(expected map[string][]string, parties []*Party) {
	for _, party1 := range parties {
		for _, party2 := range parties {
			if party1.ID == party2.ID {
				continue
			}
			list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
				FirstPartyID:  party1.ID,
				SecondPartyID: party2.ID,
			})
			if !assert.NoError(s.T(), err) {
				s.T().FailNow()
			}
			key := party1.ID + party2.ID
			// map entry and returned list should have the same length
			assert.Len(s.T(), expected[key], len(list.Items))
			// map entry should contain same items as list
			for _, l := range list.Items {
				assert.Contains(s.T(), expected[key], l.ID)
			}
		}
	}
}

func (s *Suite) testRelationshipFilterByTypeAndEither(expected map[string][]string, types []*RelationshipType, parties []*Party) {
	for _, relationshipType := range types {
		rType := relationshipType.ID
		for _, party := range parties {
			list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
				RelationshipTypeID: rType,
				EitherPartyID:      party.ID,
			})
			if !assert.NoError(s.T(), err) {
				s.T().FailNow()
			}
			key := rType + party.ID
			// map entry and returned list should have the same length
			assert.Len(s.T(), expected[key], len(list.Items))
			// map entry should contain same items as list
			for _, l := range list.Items {
				assert.Contains(s.T(), expected[key], l.ID)
			}
		}
	}
}
func (s *Suite) testRelationshipFilterByTypeAndBoth(expected map[string][]string, types []*RelationshipType, parties []*Party) {
	for _, relationshipType := range types {
		rType := relationshipType.ID
		for _, party1 := range parties {
			for _, party2 := range parties {
				list, err := s.client.Relationships().List(s.ctx, RelationshipListOptions{
					RelationshipTypeID: rType,
					FirstPartyID:       party1.ID,
					SecondPartyID:      party2.ID,
				})
				if !assert.NoError(s.T(), err) {
					s.T().FailNow()
				}
				key := rType + party1.ID + party2.ID
				// map entry and returned list should have the same length
				assert.Len(s.T(), expected[key], len(list.Items))
				// map entry should contain same items as list
				for _, l := range list.Items {
					assert.Contains(s.T(), expected[key], l.ID)
				}
			}
		}
	}
}
