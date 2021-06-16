package testing

import (
	"github.com/nrc-no/core-kafka/pkg/cases/cases"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestCaseCRUD() {
	// CREATE
	mock := "create"
	created, err := s.server.CaseClient.Create(s.ctx, &cases.Case{
		CaseTypeID:  mock,
		PartyID:     mock,
		Description: mock,
		Done:        false,
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mock, created.CaseTypeID)
	assert.Equal(s.T(), mock, created.PartyID)
	assert.Equal(s.T(), mock, created.Description)
	assert.False(s.T(), created.Done)

	// GET
	get, err := s.server.CaseClient.Get(s.ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created, get)

	// UPDATE
	updatedMock := "update"
	updated, err := s.server.CaseClient.Update(s.ctx, &cases.Case{
		ID:          created.ID,
		CaseTypeID:  updatedMock,
		PartyID:     updatedMock,
		Description: updatedMock,
		Done:        !created.Done,
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created.ID, updated.ID)
	assert.Equal(s.T(), updatedMock, updated.CaseTypeID)
	assert.Equal(s.T(), updatedMock, updated.PartyID)
	assert.Equal(s.T(), updatedMock, updated.Description)
	assert.False(s.T(), created.Done == updated.Done)

	// GET
	get, err = s.server.CaseClient.Get(s.ctx, updated.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.server.CaseClient.List(s.ctx, cases.ListOptions{})
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), list.Items, get)
}
