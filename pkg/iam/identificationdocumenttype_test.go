package iam_test

import (
	"github.com/nrc-no/core/pkg/iam"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestIdentificationDocumentType() {
	s.Run("API", func() { s.testIdentificationDocumentTypeAPI() })
}

func (s *Suite) testIdentificationDocumentTypeAPI() {
	// Create
	name := newUUID()
	created, err := s.client.IdentificationDocumentTypes().Create(s.Ctx, &iam.IdentificationDocumentType{
		Name: name,
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), name, created.Name)
	assert.NotEmpty(s.T(), created.ID)

	// Get
	get, err := s.client.IdentificationDocumentTypes().Get(s.Ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), get, created)

	// Update
	newName := newUUID()
	updated := *created
	updated.Name = newName
	_, err = s.client.IdentificationDocumentTypes().Update(s.Ctx, &updated)
	assert.NoError(s.T(), err)

	// Get
	get, err = s.client.IdentificationDocumentTypes().Get(s.Ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), get, &updated)

	// List
	list, err := s.client.IdentificationDocumentTypes().List(s.Ctx, iam.IdentificationDocumentTypeListOptions{})
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), list.Items, get)
}
