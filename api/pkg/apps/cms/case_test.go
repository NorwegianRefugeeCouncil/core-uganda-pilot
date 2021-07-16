// +build integration

package cms_test

import (
	. "github.com/nrc-no/core/pkg/apps/cms"
	"github.com/stretchr/testify/assert"
	"reflect"
)

func (s *Suite) TestCase() {
	s.Run("API", func() { s.testCaseAPI() })
	s.SetupTest()
	s.Run("List filtering", func() { s.testCaseListFilter() })
}

func (s *Suite) testCaseAPI() {
	// Create
	kase := s.mockCases(1)[0]
	created, err := s.client.Cases().Create(s.ctx, kase)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), kase, created)

	// GET
	get, err := s.client.Cases().Get(s.ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), created, get)

	// UPDATE
	kase.Done = true
	kase.FormData = &CaseTemplate{FormElements: []CaseTemplateFormElement{{
		Type: "textarea",
		Attributes: CaseTemplateFormElementAttribute{
			Label: "mock",
			Value: []string{"mock"},
		},
	}}}
	updated, err := s.client.Cases().Update(s.ctx, kase)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), kase, updated)

	// GET
	get, err = s.client.Cases().Get(s.ctx, updated.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.client.Cases().List(s.ctx, CaseListOptions{})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Contains(s.T(), list.Items, get)
}

func (s *Suite) testCaseListFilter() {
	nCases := 20
	const (
		nParties   = 5
		nTeams     = 4
		nCaseTypes = 5
		nParents   = 30
	)

	parties := []string{}
	teams := []string{}
	caseTypes := []string{}
	parents := []string{}

	cases := s.mockCases(nCases)
	for i := 0; i < nParties; i++ {
		parties = append(parties, newUUID())
	}
	for i := 0; i < nTeams; i++ {
		teams = append(teams, newUUID())
	}
	for i := 0; i < nCaseTypes; i++ {
		caseTypes = append(caseTypes, newUUID())
	}
	for i := 0; i < nParents; i++ {
		parents = append(parents, newUUID())
	}

	// Prepare test data
	for i, c := range cases {
		c.PartyID = parties[i%len(parties)]
		c.TeamID = teams[i%len(teams)]
		c.CaseTypeID = caseTypes[i%len(caseTypes)]
		c.ParentID = parents[i%len(parents)]
		c.Done = i%2 == 0
		_, err := s.client.Cases().Create(s.ctx, c)
		if err != nil {
			s.T().FailNow()
		}
	}

	s.Run("by party type", func() { s.testCaseFilterBy("PartyIDs")(cases, parties) })
	s.Run("by team", func() { s.testCaseFilterBy("TeamIDs")(cases, teams) })
	s.Run("by case type", func() { s.testCaseFilterBy("CaseTypeIDs")(cases, caseTypes) })
	s.Run("by parent", func() { s.testCaseFilterBy("ParentID")(cases, parents) })
	s.Run("by done", func() { s.testCaseFilterByDone(cases) })
}

func (s *Suite) testCaseFilterBy(field string) func(kases []*Case, search []string) {
	return func(kases []*Case, search []string) {
		hasMany := field[len(field)-1:] == "s"
		for i := 1; i <= len(search); i++ {
			searchOpts := search[0:i]
			var opt = CaseListOptions{}
			f := reflect.ValueOf(&opt).Elem().FieldByName(field)
			if hasMany {
				f.Set(reflect.ValueOf(searchOpts))
			} else {
				f.Set(reflect.ValueOf(search[i-1]))
			}
			list, err := s.client.Cases().List(s.ctx, opt)
			if err != nil {
				s.T().FailNow()
			}
			expected := []string{}
			for _, kase := range kases {
				name := field
				if hasMany {
					name = field[0 : len(field)-1]
				}
				f := reflect.ValueOf(kase).Elem().FieldByName(name)
				if (hasMany && contains(searchOpts, f.String())) || search[i-1] == f.String() {
					expected = append(expected, kase.ID)
				}
			}
			assert.Len(s.T(), list.Items, len(expected))
			for _, item := range list.Items {
				assert.Contains(s.T(), expected, item.ID)
			}
		}
	}
}

func (s *Suite) testCaseFilterByDone(kases []*Case) {
	d := true
	for i := 0; i < 2; i++ {
		list, err := s.client.Cases().List(s.ctx, CaseListOptions{Done: &d})
		if err != nil {
			s.T().FailNow()
		}
		expected := []string{}
		for _, kase := range kases {
			if kase.Done == d {
				expected = append(expected, kase.ID)
			}
		}
		assert.Len(s.T(), list.Items, len(expected))
		for _, item := range list.Items {
			assert.Contains(s.T(), expected, item.ID)
		}
		d = false
	}
}
