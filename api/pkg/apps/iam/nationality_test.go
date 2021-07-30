package iam_test

import (
	. "github.com/nrc-no/core/pkg/apps/iam"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestNationality() {
	s.Run("API", func() { s.testNationalityAPI() })
}

func (s *Suite) testNationalityAPI() {
	nationality := s.mockNationalitys(1)[0]
	nationality.CountryID = newUUID()
	nationality.TeamID = newUUID()

	// Create
	created, err := s.client.Nationalitys().Create(s.Ctx, nationality)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), nationality.CountryID, created.CountryID)
	assert.Equal(s.T(), nationality.TeamID, created.TeamID)

	// Get
	get, err := s.client.Nationalitys().Get(s.Ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), created, get)

	// Update
	// NB Nationality doesn't implement UPDATE for now

	// List
	list, err := s.client.Nationalitys().List(s.Ctx, NationalityListOptions{})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Contains(s.T(), list.Items, get)
}
