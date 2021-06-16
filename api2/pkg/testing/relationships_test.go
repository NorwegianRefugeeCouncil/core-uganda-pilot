package testing

import (
	"github.com/nrc-no/core-kafka/pkg/parties/relationships"
	"github.com/stretchr/testify/assert"
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
	if !assert.Equal(s.T(), get, created) {
		return
	}

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
	if !assert.Equal(s.T(), get, updated) {
		return
	}
	// LIST relationships types
	list, err := s.server.RelationshipClient.List(s.ctx, relationships.ListOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, get)
}
