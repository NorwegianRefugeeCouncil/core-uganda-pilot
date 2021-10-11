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
	partyId := newUUID()
	identificationDocumentTypeId := newUUID()
	documentNumber := newUUID()
	created, err := s.client.IdentificationDocuments().Create(s.Ctx, &IdentificationDocument{
		PartyID:                      partyId,
		IdentificationDocumentTypeID: identificationDocumentTypeId,
		DocumentNumber:               documentNumber,
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), partyId, created.PartyID)
	assert.Equal(s.T(), identificationDocumentTypeId, created.IdentificationDocumentTypeID)
	assert.Equal(s.T(), documentNumber, created.DocumentNumber)
	assert.NotEmpty(s.T(), created.ID)

	// Get
	get, err := s.client.IdentificationDocuments().Get(s.Ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), get, created)

	// Update
	newPartyId := newUUID()
	newIdentificationDocumentTypeId := newUUID()
	newDocumentNumber := newUUID()
	updated := *created
	updated.PartyID = newPartyId
	updated.IdentificationDocumentTypeID = newIdentificationDocumentTypeId
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
