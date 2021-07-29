package attachments_test

import (
	. "github.com/nrc-no/core/pkg/apps/attachments"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) testAttachmentAPI() {
	// Create
	attachment := s.aMo
	created, err := s.client.CaseTypes().Create(s.Ctx, caseType)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	caseType.ID = created.ID
	assert.Equal(s.T(), caseType, created)

	// Get

	// Update

	// Get

	// List
}
