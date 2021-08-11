package iam_test

import (
	. "github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestAttribute() {
	s.Run("API", func() { s.testAttributeAPI() })
	s.SetupTest()
	s.Run("List filtering", func() { s.testAttributeListFilter() })
}

func (s *Suite) testAttributeAPI() {
	// CREATE
	mock := "create"
	created, err := s.client.Attributes().Create(s.Ctx, &Attribute{
		Name:                         mock,
		PartyTypeIDs:                 []string{mock},
		IsPersonallyIdentifiableInfo: false,
		Translations: []AttributeTranslation{{
			Locale:           mock,
			LongFormulation:  mock,
			ShortFormulation: mock,
		},
		},
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.NotEmpty(s.T(), created.ID)
	assert.Equal(s.T(), mock, created.Name)
	assert.Equal(s.T(), []string{mock}, created.PartyTypeIDs)
	assert.False(s.T(), created.IsPersonallyIdentifiableInfo)
	assert.Equal(s.T(), []AttributeTranslation{{
		Locale:           mock,
		LongFormulation:  mock,
		ShortFormulation: mock,
	},
	}, created.Translations)

	// GET
	get, err := s.client.Attributes().Get(s.Ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), get, created)

	// UPDATE
	updatedMock := "update"

	updated, err := s.client.Attributes().Update(s.Ctx, &Attribute{
		ID:                           created.ID,
		Name:                         updatedMock,
		PartyTypeIDs:                 []string{updatedMock},
		IsPersonallyIdentifiableInfo: false,
		Translations: []AttributeTranslation{{
			Locale:           updatedMock,
			LongFormulation:  updatedMock,
			ShortFormulation: updatedMock,
		},
		},
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created.ID, updated.ID)
	assert.Equal(s.T(), updatedMock, updated.Name)
	assert.Equal(s.T(), []string{updatedMock}, updated.PartyTypeIDs)
	assert.False(s.T(), updated.IsPersonallyIdentifiableInfo)
	assert.Equal(s.T(), []AttributeTranslation{{
		Locale:           updatedMock,
		LongFormulation:  updatedMock,
		ShortFormulation: updatedMock,
	},
	}, updated.Translations)

	// GET
	get, err = s.client.Attributes().Get(s.Ctx, updated.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.client.Attributes().List(s.Ctx, AttributeListOptions{})
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), list.Items, get)
}

func (s *Suite) testAttributeListFilter() {

	const nAttributes = 30
	const nPartyTypeIds = 6

	// Make a couple Attributes
	attributes := s.mockAttributes(nAttributes)

	// Make a couple PartyTypeIDs
	var partyTypeIds [nPartyTypeIds]string
	for i := range partyTypeIds {
		partyTypeIds[i] = newUUID()
	}

	// Save the attributes to the DB
	for i, attribute := range attributes {
		attribute.PartyTypeIDs = partyTypeIds[0 : 1+(i%len(partyTypeIds))]
		created, err := s.client.Attributes().Create(s.Ctx, attribute)
		assert.NoError(s.T(), err)
		attribute.ID = created.ID
	}

	// Test list filtering with different PartyTypeID combinations
	for i := 1; i <= len(partyTypeIds); i++ {
		partyTypes := partyTypeIds[0:i]
		list, err := s.client.Attributes().List(ctx, AttributeListOptions{PartyTypeIDs: partyTypes})
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
