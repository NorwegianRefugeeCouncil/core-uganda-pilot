package testing

import (
	"github.com/nrc-no/core/pkg/teams"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestTeamCrud() {

	// Create team
	name := uuid.NewV4().String()
	created, err := s.server.TeamClient.Create(s.ctx, &teams.Team{
		Name: name,
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), name, created.Name)
	assert.NotEmpty(s.T(), created.ID)

	// Get team
	get, err := s.server.TeamClient.Get(s.ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Equal(s.T(), get, created) {
		return
	}

	// List teams
	list, err := s.server.TeamClient.List(s.ctx, teams.ListOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, get)

}
