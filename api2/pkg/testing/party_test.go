package testing

import (
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestPartyCRUD() {
	// CREATE
	mock := "create"
	created, err := s.server.PartyClient.Create(s.ctx, &parties.Party{
		PartyTypes: []string{mock},
		Attributes: parties.PartyAttributes{"mock": []string{mock}},
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), []string{mock}, created.PartyTypes)
	assert.Equal(s.T(), parties.PartyAttributes{"mock": []string{mock}}, created.Attributes)

	// GET
	get, err := s.server.PartyClient.Get(s.ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created, get)

	// UPDATE
	updatedMock := "update"
	updated, err := s.server.PartyClient.Update(s.ctx, &parties.Party{
		ID:         created.ID,
		PartyTypes: []string{updatedMock},
		Attributes: parties.PartyAttributes{"mock": []string{updatedMock}},
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created.ID, updated.ID)
	assert.Equal(s.T(), []string{updatedMock}, updated.PartyTypes)
	assert.Equal(s.T(), parties.PartyAttributes{"mock": []string{updatedMock}}, updated.Attributes)

	// GET
	get, err = s.server.PartyClient.Get(s.ctx, updated.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.server.PartyClient.List(s.ctx, parties.ListOptions{})
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), list.Items, get)
}
