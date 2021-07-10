// +build integration

package iam_test

import (
	. "github.com/nrc-no/core/pkg/apps/iam"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestMembership() {
	s.Run("API", func() { s.testMembershipAPI() })
}

func (s *Suite) testMembershipAPI() {
	membership := s.mockMemberships(1)[0]
	membership.TeamID = newUUID()
	membership.IndividualID = newUUID()

	// Create
	created, err := s.client.Memberships().Create(s.ctx, membership)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), membership, created)

	// Get
	get, err := s.client.Memberships().Get(s.ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), created, get)

	// NB Membership doesn't implement UPDATE for now
	//// Update
	//membership.TeamID = newUUID()
	//membership.IndividualID = newUUID()
	//updated, err := s.client.Memberships().Update(s.ctx, &membership)
	//if !assert.NoError(s.T(), err) {
	//	s.T().FailNow()
	//}
	//assert.Equal(s.T(), membership, *updated)
	//
	//// Get
	//get, err = s.client.Memberships().Get(s.ctx, updated.ID)
	//if !assert.NoError(s.T(), err) {
	//	s.T().FailNow()
	//}
	//assert.Equal(s.T(), updated, get)

	// List
	list, err := s.client.Memberships().List(s.ctx, MembershipListOptions{})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Contains(s.T(), list.Items, get)
}
