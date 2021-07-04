// +build integration

package iam

import (
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestAttributeCRUD() {

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
