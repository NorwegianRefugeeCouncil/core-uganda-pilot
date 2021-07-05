// +build integration

package iam

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestTeamCrud() {

	// Create team
	name := uuid.NewV4().String()
	created, err := s.client.Teams().Create(s.ctx, &Team{
		Name: name,
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), name, created.Name)
	assert.NotEmpty(s.T(), created.ID)

	// Get team
	get, err := s.client.Teams().Get(s.ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Equal(s.T(), get, created) {
		return
	}

	// List teams
	list, err := s.client.Teams().List(s.ctx, TeamListOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, get)

}
