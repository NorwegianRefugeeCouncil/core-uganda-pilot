package testing

import (
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestCaseTypeCRUD() {
	// CREATE
	mock := "create"
	created, err := s.server.CaseTypeClient.Create(s.ctx, &casetypes.CaseType{
		Name:        mock,
		PartyTypeID: mock,
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mock, created.Name)
	assert.Equal(s.T(), mock, created.PartyTypeID)

	// GET
	get, err := s.server.CaseTypeClient.Get(s.ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created, get)

	// UPDATE
	updatedMock := "update"
	updated, err := s.server.CaseTypeClient.Update(s.ctx, &casetypes.CaseType{
		ID:          created.ID,
		Name:        updatedMock,
		PartyTypeID: updatedMock,
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created.ID, updated.ID)
	assert.Equal(s.T(), updatedMock, updated.Name)
	assert.Equal(s.T(), updatedMock, updated.PartyTypeID)

	// GET
	get, err = s.server.CaseTypeClient.Get(s.ctx, updated.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.server.CaseTypeClient.List(s.ctx, casetypes.ListOptions{})
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), list.Items, get)
}
