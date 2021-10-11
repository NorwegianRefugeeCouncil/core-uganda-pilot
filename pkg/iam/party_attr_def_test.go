package iam_test

import (
	"github.com/nrc-no/core/internal/utils"
	. "github.com/nrc-no/core/pkg/iam"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestPartyAttributeDefinition() {
	s.Run("API", func() { s.testPartyAttributeDefinitionAPI() })
	s.SetupTest()
	s.Run("List filtering", func() { s.testPartyAttributeDefinitionListFilter() })
}

func (s *Suite) testPartyAttributeDefinitionAPI() {
	// CREATE
	create := mockPartyAttributeDefinition()
	created, err := s.client.PartyAttributeDefinitions().Create(s.Ctx, create)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.NotEmpty(s.T(), created.ID)
	create.ID = created.ID
	assert.Equal(s.T(), create, created)
	// GET
	get, err := s.client.PartyAttributeDefinitions().Get(s.Ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), get, created)

	// UPDATE
	update := created
	update.PartyTypeIDs = []string{newUUID()}
	update.CountryID = newUUID()
	update.IsPersonallyIdentifiableInfo = !update.IsPersonallyIdentifiableInfo
	updated, err := s.client.PartyAttributeDefinitions().Update(s.Ctx, update)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created.ID, updated.ID)
	assert.Equal(s.T(), update, updated)

	// GET
	get, err = s.client.PartyAttributeDefinitions().Get(s.Ctx, updated.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.client.PartyAttributeDefinitions().List(s.Ctx, PartyAttributeDefinitionListOptions{})
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), list.Items, get)
}

func (s *Suite) testPartyAttributeDefinitionListFilter() {

	const nAttributes = 10
	const nPartyTypeIds = 6

	// Make a couple Attributes
	attributes := mockPartyAttributeDefinitions(nAttributes)

	// Make a couple PartyTypeIDs
	var partyTypeIds [nPartyTypeIds]string
	for i := range partyTypeIds {
		partyTypeIds[i] = newUUID()
	}

	// Save the attributes to the DB
	for i, attribute := range attributes {
		attribute.PartyTypeIDs = partyTypeIds[0 : 1+(i%len(partyTypeIds))]
		created, err := s.client.PartyAttributeDefinitions().Create(s.Ctx, attribute)
		assert.NoError(s.T(), err)
		attribute.ID = created.ID
	}

	// Test list filtering with different PartyTypeID combinations
	for i := 1; i <= len(partyTypeIds); i++ {
		partyTypes := partyTypeIds[0:i]
		list, err := s.client.PartyAttributeDefinitions().List(s.Ctx, PartyAttributeDefinitionListOptions{PartyTypeIDs: partyTypes})
		assert.NoError(s.T(), err)

		// Get expected items
		var expected []string
		for _, a := range attributes {
			var include = true
			for _, p := range partyTypes {
				if !utils.Contains(a.PartyTypeIDs, p) {
					include = false
					break
				}
			}
			if include {
				expected = append(expected, a.ID)
			}
		}

		// Check length
		assert.Equal(s.T(), len(expected), len(list.Items))

		// Check contents
		for _, item := range list.Items {
			assert.Contains(s.T(), expected, item.ID)
		}
	}
}
