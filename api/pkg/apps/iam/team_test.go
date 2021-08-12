package iam_test

import "github.com/stretchr/testify/assert"
import . "github.com/nrc-no/core/pkg/apps/iam"

func (s *Suite) TestTeamAPI() {
	// Create team
	name := newUUID()
	created, err := s.client.Teams().Create(s.Ctx, &Team{
		Name: name,
	})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), name, created.Name)
	assert.NotEmpty(s.T(), created.ID)

	// Get team
	get, err := s.client.Teams().Get(s.Ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	if !assert.Equal(s.T(), get, created) {
		return
	}

	// List teams
	list, err := s.client.Teams().List(s.Ctx, TeamListOptions{})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Contains(s.T(), list.Items, get)
}
