package iam_test

import (
	. "github.com/nrc-no/core/pkg/apps/iam"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestMembership() {
	s.Run("API", func() { s.testMembershipAPI() })
}

func (s *Suite) testMembershipAPI() {
	membership := mockMembership()
	membership.TeamID = newUUID()
	membership.IndividualID = newUUID()

	// Create
	created, err := s.client.Memberships().Create(s.Ctx, membership)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), membership.TeamID, created.TeamID)
	assert.Equal(s.T(), membership.IndividualID, created.IndividualID)

	// Get
	get, err := s.client.Memberships().Get(s.Ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), created, get)

	// Update
	// NB Membership doesn't implement UPDATE for now

	// List
	list, err := s.client.Memberships().List(s.Ctx, MembershipListOptions{})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Contains(s.T(), list.Items, get)
}
