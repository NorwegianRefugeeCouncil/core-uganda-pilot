// +build integration

package iam_test

import "github.com/stretchr/testify/assert"
import . "github.com/nrc-no/core/pkg/apps/iam"

func (s *Suite) TestTeam() {
	s.Run("API", func() { s.testTeamAPI() })
}

func (s *Suite) testTeamAPI() {
	// Create team
	name := newUUID()
	created, err := s.client.Teams().Create(s.ctx, &Team{
		Name: name,
	})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), name, created.Name)
	assert.NotEmpty(s.T(), created.ID)

	// Get team
	get, err := s.client.Teams().Get(s.ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	if !assert.Equal(s.T(), get, created) {
		return
	}

	// List teams
	list, err := s.client.Teams().List(s.ctx, TeamListOptions{})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Contains(s.T(), list.Items, get)
}
