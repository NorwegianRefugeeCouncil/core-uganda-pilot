package testing

import (
	"github.com/nrc-no/core-kafka/pkg/expressions"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestAttributesCRUD() {
	// CREATE
	mock := "create"
	created, err := s.server.AttributeClient.Create(s.ctx, &attributes.Attribute{
		Name:                         mock,
		ValueType:                    expressions.ValueType{},
		PartyTypeIDs:                 []string{mock},
		IsPersonallyIdentifiableInfo: false,
		Translations: []attributes.AttributeTranslation{{
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
	assert.Empty(s.T(), created.ValueType)
	assert.Equal(s.T(), []string{mock}, created.PartyTypeIDs)
	assert.False(s.T(), created.IsPersonallyIdentifiableInfo)
	assert.Equal(s.T(), []attributes.AttributeTranslation{{
		Locale:           mock,
		LongFormulation:  mock,
		ShortFormulation: mock,
	},
	}, created.Translations)

	// GET
	get, err := s.server.AttributeClient.Get(s.ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), get, created)

	// UPDATE
	updatedMock := "update"

	updated, err := s.server.AttributeClient.Update(s.ctx, &attributes.Attribute{
		ID:                           created.ID,
		Name:                         updatedMock,
		ValueType:                    expressions.ValueType{},
		PartyTypeIDs:                 []string{updatedMock},
		IsPersonallyIdentifiableInfo: false,
		Translations: []attributes.AttributeTranslation{{
			Locale:           updatedMock,
			LongFormulation:  updatedMock,
			ShortFormulation: updatedMock,
		},
		},
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created.ID, updated.ID)
	assert.Equal(s.T(), updatedMock, updated.Name)
	assert.Empty(s.T(), updated.ValueType)
	assert.Equal(s.T(), []string{updatedMock}, updated.PartyTypeIDs)
	assert.False(s.T(), updated.IsPersonallyIdentifiableInfo)
	assert.Equal(s.T(), []attributes.AttributeTranslation{{
		Locale:           updatedMock,
		LongFormulation:  updatedMock,
		ShortFormulation: updatedMock,
	},
	}, updated.Translations)

	// GET
	get, err = s.server.AttributeClient.Get(s.ctx, updated.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.server.AttributeClient.List(s.ctx)
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), list.Items, get)

}
