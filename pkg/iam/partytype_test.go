package iam_test

import (
	. "github.com/nrc-no/core/pkg/iam"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestPartyType() {

	// Create party type
	name := newUUID()
	created, err := s.client.PartyTypes().Create(s.Ctx, &PartyType{
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
	get, err := s.client.PartyTypes().Get(s.Ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Equal(s.T(), get, created) {
		return
	}

	// List party types
	list, err := s.client.PartyTypes().List(s.Ctx, PartyTypeListOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, get)

}
