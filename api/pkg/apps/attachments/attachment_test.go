package attachments_test

import (
	"github.com/nrc-no/core/pkg/apps/attachments"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestAttachments() {
	s.Run("API", func() { s.testAttachmentAPI() })
}

func (s *Suite) testAttachmentAPI() {
	var imaginaryOwner = s.NewUUID()
	var anotherImaginaryOwner = s.NewUUID()

	// Create
	attachment := s.aMockAttachment(imaginaryOwner)
	created, err := s.client.Create(s.Ctx, attachment)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	attachment.ID = created.ID
	assert.Equal(s.T(), attachment, created)

	// Get
	get, err := s.client.Get(s.Ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), created, get)

	// Update
	attachment.Body = "{\"data\":\"updated\"}"
	attachment.AttachedToID = anotherImaginaryOwner
	updated, err := s.client.Update(s.Ctx, attachment)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), attachment, updated)

	// Get
	get, err = s.client.Get(s.Ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), created, get)

	// List
	list, err := s.client.List(s.Ctx, attachments.AttachmentListOptions{})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Contains(s.T(), list.Items, get)
}
