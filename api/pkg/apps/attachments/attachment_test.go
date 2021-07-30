package attachments_test

import (
	"github.com/nrc-no/core/pkg/apps/attachments"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestAttachments() {
	s.Run("API", func() { s.testAttachmentAPI() })
	s.SetupTest()
	s.Run("List filtering", func() { s.testAttachmentListFilter() })
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
	assert.Equal(s.T(), updated, get)

	// List
	list, err := s.client.List(s.Ctx, attachments.AttachmentListOptions{})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Contains(s.T(), list.Items, get)
}

func (s *Suite) testAttachmentListFilter() {
	const nOwners = 10
	const nAttachmentsPerOwner = 10

	var owners []string
	// create owners
	for i := 0; i < nOwners; i++ {
		owners = append(owners, s.NewUUID())
	}

	for _, owner := range owners {
		attachmentsForOwner := []*attachments.Attachment{}

		// create attachments for each owner
		for i := 0; i < nAttachmentsPerOwner; i++ {
			attachment := s.aMockAttachment(owner)
			created, err := s.client.Create(s.Ctx, attachment)
			if !assert.NoError(s.T(), err) {
				s.T().FailNow()
			}
			attachment.ID = created.ID
			assert.Equal(s.T(), attachment, created)
			attachmentsForOwner = append(attachmentsForOwner, attachment)
		}

		// test filter for each owner returns expected amount of records
		list, err := s.client.List(s.Ctx, attachments.AttachmentListOptions{
			AttachedToID: owner,
		})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		assert.Equal(s.T(), len(list.Items), nAttachmentsPerOwner)

		// test all expected attachments are in returned list
		for _, attachment := range attachmentsForOwner {
			assert.Contains(s.T(), list.Items, attachment)
		}
	}
}
