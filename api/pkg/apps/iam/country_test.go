package iam_test

import "github.com/stretchr/testify/assert"
import . "github.com/nrc-no/core/pkg/apps/iam"

func (s *Suite) TestCountry() {
	s.Run("API", func() { s.testCountryAPI() })
}

func (s *Suite) testCountryAPI() {
	// Create country
	name := newUUID()
	created, err := s.client.Countrys().Create(s.Ctx, &Country{
		Name: name,
	})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), name, created.Name)
	assert.NotEmpty(s.T(), created.ID)

	// Get country
	get, err := s.client.Countrys().Get(s.Ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	if !assert.Equal(s.T(), get, created) {
		return
	}

	// List countrys
	list, err := s.client.Countrys().List(s.Ctx, CountryListOptions{})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Contains(s.T(), list.Items, get)
}
