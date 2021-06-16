package testing

import (
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestRelationShipTypeCRUD() {
	// CREATE relationship type
	mock := uuid.NewV4().String()
	created, err := s.server.RelationshipTypeClient.Create(s.ctx, &relationshiptypes.RelationshipType{
		//ID:              mock,
		IsDirectional:   true,
		Name:            mock,
		FirstPartyRole:  mock,
		SecondPartyRole: mock,
		Rules:           []relationshiptypes.RelationshipTypeRule{{relationshiptypes.PartyTypeRule{FirstPartyType: mock, SecondPartyType: mock}}},
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	//assert.Equal(s.T(), mock, created.ID)
	assert.True(s.T(), created.IsDirectional)
	assert.Equal(s.T(), mock, created.Name)
	assert.Equal(s.T(), mock, created.FirstPartyRole)
	assert.Equal(s.T(), mock, created.SecondPartyRole)
	assert.IsType(s.T(), []relationshiptypes.RelationshipTypeRule{{relationshiptypes.PartyTypeRule{FirstPartyType: mock, SecondPartyType: mock}}}, created.Rules)

	// GET relationship type
	get, err := s.server.RelationshipTypeClient.Get(s.ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Equal(s.T(), get, created) {
		return
	}

	// UPDATE relationship type
	updatedMock := "update"

	updated, err := s.server.RelationshipTypeClient.Update(s.ctx, &relationshiptypes.RelationshipType{
		ID:              created.ID,
		IsDirectional:   !created.IsDirectional,
		Name:            updatedMock,
		FirstPartyRole:  updatedMock,
		SecondPartyRole: updatedMock,
		Rules:           []relationshiptypes.RelationshipTypeRule{{relationshiptypes.PartyTypeRule{FirstPartyType: updatedMock, SecondPartyType: updatedMock}}},
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), created.ID, updated.ID)
	assert.Equal(s.T(), updatedMock, updated.Name)
	assert.False(s.T(), created.IsDirectional == updated.IsDirectional)
	assert.Equal(s.T(), updatedMock, updated.FirstPartyRole)
	assert.Equal(s.T(), updatedMock, updated.SecondPartyRole)
	assert.IsType(s.T(), []relationshiptypes.RelationshipTypeRule{{relationshiptypes.PartyTypeRule{FirstPartyType: updatedMock, SecondPartyType: updatedMock}}}, updated.Rules)

	// GET relationship type
	get, err = s.server.RelationshipTypeClient.Get(s.ctx, updated.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Equal(s.T(), get, updated) {
		return
	}
	// LIST relationship types
	list, err := s.server.RelationshipTypeClient.List(s.ctx)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, get)
}
