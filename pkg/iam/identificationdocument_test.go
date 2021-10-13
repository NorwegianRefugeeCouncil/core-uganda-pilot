package iam_test

import (
	. "github.com/nrc-no/core/pkg/iam"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestIdentificationDocument() {
	s.Run("API", func() { s.testIdentificationDocumentAPI() })
}

func (s *Suite) testIdentificationDocumentAPI() {
	// Create
	partyID := newUUID()
	identificationDocumentTypeID := newUUID()
	documentNumber := newUUID()
	created, err := s.client.IdentificationDocuments().Create(s.Ctx, &IdentificationDocument{
		PartyID:                      partyID,
		IdentificationDocumentTypeID: identificationDocumentTypeID,
		DocumentNumber:               documentNumber,
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), partyID, created.PartyID)
	assert.Equal(s.T(), identificationDocumentTypeID, created.IdentificationDocumentTypeID)
	assert.Equal(s.T(), documentNumber, created.DocumentNumber)
	assert.NotEmpty(s.T(), created.ID)

	// Get
	get, err := s.client.IdentificationDocuments().Get(s.Ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), get, created)

	// Update
	newPartyID := newUUID()
	newIdentificationDocumentTypeID := newUUID()
	newDocumentNumber := newUUID()
	updated := *created
	updated.PartyID = newPartyID
	updated.IdentificationDocumentTypeID = newIdentificationDocumentTypeID
	updated.DocumentNumber = newDocumentNumber
	_, err = s.client.IdentificationDocuments().Update(s.Ctx, &updated)
	assert.NoError(s.T(), err)

	// Get
	get, err = s.client.IdentificationDocuments().Get(s.Ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), get, &updated)

	// List
	list, err := s.client.IdentificationDocuments().List(s.Ctx, IdentificationDocumentListOptions{})
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), list.Items, get)

	// Delete
	err = s.client.IdentificationDocuments().Delete(s.Ctx, created.ID)
	assert.NoError(s.T(), err)
}
