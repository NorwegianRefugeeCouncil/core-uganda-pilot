package iam_test

import (
	. "github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestIndividual() {
	s.Run("API", func() { s.testIndividualAPI() })
	s.SetupTest()
	s.Run("List filtering", func() { s.testIndividualListFilter() })
}

func (s *Suite) testIndividualAPI() {
	// CREATE
	individual := s.mockIndividuals(1)[0]
	created, err := s.client.Individuals().Create(s.Ctx, individual)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	individual.ID = created.ID
	assert.Equal(s.T(), created, individual)

	// GET
	get, err := s.client.Individuals().Get(s.Ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), created, get)

	// UPDATE
	individual.Attributes.Set(FirstNameAttribute.ID, "updated")
	individual.Attributes.Set(LastNameAttribute.ID, "updated")
	updated, err := s.client.Individuals().Update(s.Ctx, individual)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), individual, updated)

	// GET
	get, err = s.client.Individuals().Get(s.Ctx, updated.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.client.Individuals().List(s.Ctx, IndividualListOptions{})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Contains(s.T(), list.Items, get)
}

func (s *Suite) testIndividualListFilter() {
	nIndividuals := 10
	nPartyTypes := 3
	nAttributes := 5

	// Make a couple Individuals
	individuals := s.mockIndividuals(nIndividuals)
	// ... and a couple PartyTypes
	partyTypes := make([]string, 0)
	for i := 0; i < nPartyTypes; i++ {
		partyTypes = append(partyTypes, newUUID())
	}
	// ... and a couple Attributes
	attributes := make([]string, 0)
	for i := 0; i < nAttributes; i++ {
		attributes = append(attributes, newUUID())
	}

	// Prepare the test data
	//individualsFromPartyTypes := make(map[string][]string)
	//individualsFromAttributes := make(map[string][]string)
	for i, individual := range individuals {
		// Individuals have the IndividualPartyType by default, lets give them a couple more
		n := i % len(partyTypes)
		individual.PartyTypeIDs = append(individual.PartyTypeIDs, partyTypes[0:n]...)
		//for _, partyType := range individual.PartyTypeIDs {
		//	individualsFromPartyTypes[partyType] = append(individualsFromPartyTypes[partyType], individual.ID)
		//}

		// Individuals have the FirstName and LastName attributes by default, lets give them a couple more
		m := i % len(attributes)
		for _, attribute := range attributes[0:m] {
			individual.Attributes.Set(attribute, "mock")
			//individualsFromAttributes[attribute] = append(individualsFromAttributes[attribute], individual.ID)
		}

		// Save the individual to the DB
		created, err := s.client.Individuals().Create(s.Ctx, individual)
		assert.NoError(s.T(), err)
		individual.ID = created.ID
	}

	s.Run("by party type", func() {
		s.testIndividualFilterByPartyType(individuals, partyTypes)
	})
	s.Run("by attribute", func() {
		s.testIndividualFilterByAttribute(individuals, attributes)
	})
	s.Run("by party type and attribute", func() {
		s.testIndividualFilterByPartyTypeAndAttribute(individuals, partyTypes, attributes)
	})
}

func (s *Suite) testIndividualFilterByPartyType(individuals []*Individual, partyTypeIds []string) {
	for i := 1; i <= len(partyTypeIds); i++ {
		types := append([]string{IndividualPartyType.ID}, partyTypeIds[0:i]...)
		list, err := s.client.Individuals().List(s.Ctx, IndividualListOptions{
			PartyTypeIDs: types,
		})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		// Get expected items
		var expected []string
		for _, individual := range individuals {
			var include = true
			for _, wantedPartyTypeId := range types {
				if !utils.Contains(individual.PartyTypeIDs, wantedPartyTypeId) {
					include = false
					break
				}
			}
			if include {
				expected = append(expected, individual.ID)
			}
		}

		// Compare list lengths
		assert.Len(s.T(), list.Items, len(expected))
		// ... and contents
		for _, item := range list.Items {
			assert.Contains(s.T(), expected, item.ID)
		}
	}
}

func (s *Suite) testIndividualFilterByAttribute(individuals []*Individual, attributes []string) {
	for i := 1; i <= len(attributes); i++ {
		wantedAttributes := attributes[0:i]
		attributesOptions := make(map[string]string)
		for _, attribute := range wantedAttributes {
			attributesOptions[attribute] = "mock"
		}
		options := IndividualListOptions{Attributes: attributesOptions}
		list, err := s.client.Individuals().List(s.Ctx, options)
		assert.NoError(s.T(), err)

		// Get expected items
		var expected []string
		for _, individual := range individuals {
			var include = true
			for id := range attributesOptions {
				if !individual.HasAttribute(id) {
					include = false
					break
				}
			}
			if include {
				expected = append(expected, individual.ID)
			}
		}

		// Compare list lengths
		assert.Len(s.T(), list.Items, len(expected))
		// ... and contents
		for _, item := range list.Items {
			assert.Contains(s.T(), expected, item.ID)
		}
	}
}

func (s *Suite) testIndividualFilterByPartyTypeAndAttribute(individuals []*Individual, types []string, attributes []string) {
	for i := 1; i <= len(attributes); i++ {
		for j := 1; i <= len(types); i++ {
			partyTypeIds := append([]string{IndividualPartyType.ID}, types[0:j]...)
			attributeOptions := make(map[string]string)
			for _, attribute := range attributes[0:i] {
				attributeOptions[attribute] = "mock"
			}
			list, err := s.client.Individuals().List(s.Ctx, IndividualListOptions{
				PartyTypeIDs: partyTypeIds,
				Attributes:   attributeOptions,
			})
			if !assert.NoError(s.T(), err) {
				s.T().FailNow()
			}
			// Get expected items
			var expected []string
			for _, individual := range individuals {
				var include = true
				for _, partyType := range partyTypeIds {
					if !utils.Contains(individual.PartyTypeIDs, partyType) {
						include = false
						break
					}
				}
				if include {
					for id := range attributeOptions {
						if !individual.HasAttribute(id) {
							include = false
							break
						}
					}
				}
				if include {
					expected = append(expected, individual.ID)
				}
			}

			// Compare list lengths
			assert.Len(s.T(), expected, len(list.Items))
			// ... and contents
			for _, item := range list.Items {
				assert.Contains(s.T(), expected, item.ID)
			}
		}
	}
}
