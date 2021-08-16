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
	created, err := s.client.Relationships().Create(s.Ctx, relationship)
	if !assert.NoError(s.T(), err) {
		return
	}
	relationship.ID = created.ID
	assert.Equal(s.T(), relationship.RelationshipTypeID, created.RelationshipTypeID)
	assert.Equal(s.T(), relationship.FirstPartyID, created.FirstPartyID)
	assert.Equal(s.T(), relationship.SecondPartyID, created.SecondPartyID)

	// GET relationship
	get, err := s.client.Relationships().Get(s.Ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), created, get)

	// UPDATE relationships type
	relationship.RelationshipTypeID = newUUID()
	relationship.FirstPartyID = newUUID()
	relationship.SecondPartyID = newUUID()

	updated, err := s.client.Relationships().Update(s.Ctx, relationship)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), relationship.ID, updated.ID)
	assert.Equal(s.T(), relationship.RelationshipTypeID, updated.RelationshipTypeID)
	assert.Equal(s.T(), relationship.FirstPartyID, updated.FirstPartyID)
	assert.Equal(s.T(), relationship.SecondPartyID, updated.SecondPartyID)

	// GET relationships type
	get, err = s.client.Relationships().Get(s.Ctx, updated.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Equal(s.T(), updated, get) {
		return
	}
	// LIST relationships types
	list, err := s.client.Relationships().List(s.Ctx, RelationshipListOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, get)
}

func (s *Suite) testRelationshipListFilter() {
	const nRelationships = 10000
	const nRelationshipTypes = 3
	const nParties = 4

	relationships := s.mockRelationships(nRelationships)
	relationshipTypes := []string{}
	parties := []string{}

	for i := 0; i < nRelationshipTypes; i++ {
		relationshipTypes = append(relationshipTypes, newUUID())
	}

	for i := 0; i < nParties; i++ {
		parties = append(parties, newUUID())
	}

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
		rType := relationshipTypes[i%len(relationshipTypes)]
		party1 := parties[i%len(parties)]
		party2 := parties[(i+1)%len(parties)]
		relationship.RelationshipTypeID = rType
		relationship.FirstPartyID = party1
		relationship.SecondPartyID = party2

		// Save the relationship to the DB
		created, err := s.client.Relationships().Create(s.Ctx, relationship)
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		relationship.ID = created.ID

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

func (s *Suite) testRelationshipFilterByType(expected map[string][]string, types []string) {
	for _, relationshipType := range types {
		rType := relationshipType
		list, err := s.client.Relationships().List(s.Ctx, RelationshipListOptions{
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

func (s *Suite) testRelationshipFilterByFirstParty(expected map[string][]string, parties []string) {
	for _, party := range parties {
		list, err := s.client.Relationships().List(s.Ctx, RelationshipListOptions{
			FirstPartyID: party,
		})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		// map entry and returned list should have the same length
		assert.Len(s.T(), expected[party], len(list.Items))
		// map entry should contain same items as list
		for _, l := range list.Items {
			assert.Contains(s.T(), expected[party], l.ID)
		}
	}
}

func (s *Suite) testRelationshipFilterBySecondParty(expected map[string][]string, parties []string) {
	for _, party := range parties {
		list, err := s.client.Relationships().List(s.Ctx, RelationshipListOptions{
			SecondPartyID: party,
		})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		// map entry and returned list should have the same length
		assert.Len(s.T(), expected[party], len(list.Items))
		// map entry should contain same items as list
		for _, l := range list.Items {
			assert.Contains(s.T(), expected[party], l.ID)
		}
	}
}

func (s *Suite) testRelationshipFilterByEitherParty(expected map[string][]string, parties []string) {
	for _, party := range parties {
		list, err := s.client.Relationships().List(s.Ctx, RelationshipListOptions{
			EitherPartyID: party,
		})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		// map entry and returned list should have the same length
		assert.Len(s.T(), list.Items, len(expected[party]))
		// map entry should contain same items as list
		for _, l := range list.Items {
			assert.Contains(s.T(), expected[party], l.ID)
		}
	}
}

func (s *Suite) testRelationshipFilterByBothParties(expected map[string][]string, parties []string) {
	for _, party1 := range parties {
		for _, party2 := range parties {
			if party1 == party2 {
				continue
			}
			list, err := s.client.Relationships().List(s.Ctx, RelationshipListOptions{
				FirstPartyID:  party1,
				SecondPartyID: party2,
			})
			if !assert.NoError(s.T(), err) {
				s.T().FailNow()
			}
			key := party1 + party2
			// map entry and returned list should have the same length
			assert.Len(s.T(), expected[key], len(list.Items))
			// map entry should contain same items as list
			for _, l := range list.Items {
				assert.Contains(s.T(), expected[key], l.ID)
			}
		}
	}
}

func (s *Suite) testRelationshipFilterByTypeAndEither(expected map[string][]string, types []string, parties []string) {
	for _, relationshipType := range types {
		rType := relationshipType
		for _, party := range parties {
			list, err := s.client.Relationships().List(s.Ctx, RelationshipListOptions{
				RelationshipTypeID: rType,
				EitherPartyID:      party,
			})
			if !assert.NoError(s.T(), err) {
				s.T().FailNow()
			}
			key := rType + party
			// map entry and returned list should have the same length
			assert.Len(s.T(), expected[key], len(list.Items))
			// map entry should contain same items as list
			for _, l := range list.Items {
				assert.Contains(s.T(), expected[key], l.ID)
			}
		}
	}
}
func (s *Suite) testRelationshipFilterByTypeAndBoth(expected map[string][]string, types []string, parties []string) {
	for _, relationshipType := range types {
		rType := relationshipType
		for _, party1 := range parties {
			for _, party2 := range parties {
				list, err := s.client.Relationships().List(s.Ctx, RelationshipListOptions{
					RelationshipTypeID: rType,
					FirstPartyID:       party1,
					SecondPartyID:      party2,
				})
				if !assert.NoError(s.T(), err) {
					s.T().FailNow()
				}
				key := rType + party1 + party2
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
