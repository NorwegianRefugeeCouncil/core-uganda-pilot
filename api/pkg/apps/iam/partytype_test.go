// +build integration

package iam

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestPartyTypeCRUD() {

	// Create party type
	name := uuid.NewV4().String()
	created, err := s.client.PartyTypes().Create(s.ctx, &PartyType{
		Name:      name,
		IsBuiltIn: false,
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), name, created.Name)
	assert.NotEmpty(s.T(), created.ID)
	assert.False(s.T(), created.IsBuiltIn)

	// Get party type
	get, err := s.client.PartyTypes().Get(s.ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Equal(s.T(), get, created) {
		return
	}

	// List party types
	list, err := s.client.PartyTypes().List(s.ctx, PartyTypeListOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, get)

}
