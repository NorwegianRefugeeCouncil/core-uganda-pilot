// +build integration

package iam_test

import (
	. "github.com/nrc-no/core/pkg/apps/iam"
	"github.com/stretchr/testify/assert"
	"math/rand"
)

func (s *Suite) TestAttribute() {
	s.Run("API", func() { s.testAttributeAPI() })
	s.SetupTest()
	s.Run("List filtering", func() { s.testAttributeListFilter() })
}

func (s *Suite) testAttributeAPI() {
	// CREATE
	mock := "create"
	created, err := s.client.Attributes().Create(s.ctx, &Attribute{
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
	get, err := s.client.Attributes().Get(s.ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), get, created)

	// UPDATE
	updatedMock := "update"

	updated, err := s.client.Attributes().Update(s.ctx, &Attribute{
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
	get, err = s.client.Attributes().Get(s.ctx, updated.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.client.Attributes().List(s.ctx, AttributeListOptions{})
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), list.Items, get)
}

func (s *Suite) testAttributeListFilter() {

	nPartyTypeIds := 30
	nAttributes := 200
	// Make a couple PartyTypeIDs
	var partyTypeIds []string
	for i := 0; i < nPartyTypeIds; i++ {
		partyTypeIds = append(partyTypeIds, newUUID())
	}
	// Make a couple Attributes
	var attributes []Attribute
	for i := 0; i < nAttributes; i++ {
		attributes = append(attributes, newRandomAttribute(partyTypeIds))
	}
	// Map PartyTypeIDs to the Attribute IDs
	attributesFromPartyTypes := map[string][]string{}
	for _, a := range attributes {
		for _, p := range a.PartyTypeIDs {
			attributesFromPartyTypes[p] = append(attributesFromPartyTypes[p], a.ID)
		}
	}
	// Save the attributes to the DB
	for _, attribute := range attributes {
		_, err := s.client.Attributes().Create(s.ctx, &attribute)
		assert.NoError(s.T(), err)
	}
	// Test list filtering with different PartyTypeID combinations
	for i := 1; i < len(partyTypeIds); i++ {
		partyTypeIdSlice := partyTypeIds[0:i]
		options := AttributeListOptions{PartyTypeIDs: partyTypeIdSlice}
		list, err := s.client.Attributes().List(ctx, options)
		assert.NoError(s.T(), err)

		// Get expected items
		var expected []string
		for _, a := range attributes {
			var include = true
			for _, p := range partyTypeIdSlice {
				if !contains(a.PartyTypeIDs, p) {
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

func newRandomAttribute(partyTypeIds []string) Attribute {
	if len(partyTypeIds) != 0 {
		n := rand.Intn(len(partyTypeIds)-1) + 1
		return Attribute{
			ID:                           newUUID(),
			Name:                         "",
			PartyTypeIDs:                 partyTypeIds[0:n],
			IsPersonallyIdentifiableInfo: false,
			Translations:                 nil,
		}
	}
	return Attribute{}
}
