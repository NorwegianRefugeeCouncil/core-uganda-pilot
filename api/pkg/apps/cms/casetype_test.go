package cms_test

import (
	. "github.com/nrc-no/core/pkg/apps/cms"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestCaseType() {
	s.Run("API", func() { s.testCaseTypeAPI() })
	s.SetupTest()
	s.Run("List filtering", func() { s.testCaseTypeListFilter() })
}

func (s *Suite) testCaseTypeAPI() {
	// Create
	caseType := s.mockCaseTypes(1)[0]
	created, err := s.client.CaseTypes().Create(s.Ctx, caseType)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	caseType.ID = created.ID
	assert.Equal(s.T(), caseType, created)

	// GET
	get, err := s.client.CaseTypes().Get(s.Ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), created, get)

	// UPDATE
	caseType.Name = "updated"
	caseType.PartyTypeID = newUUID()
	caseType.TeamID = newUUID()
	caseType.Template = &CaseTemplate{FormElements: []CaseTemplateFormElement{{Type: "textarea", Attributes: CaseTemplateFormElementAttribute{Label: "updated"}}}}
	updated, err := s.client.CaseTypes().Update(s.Ctx, caseType)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), caseType, updated)

	// GET
	get, err = s.client.CaseTypes().Get(s.Ctx, updated.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.client.CaseTypes().List(s.Ctx, CaseTypeListOptions{})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Contains(s.T(), list.Items, get)
}

func (s *Suite) testCaseTypeListFilter() {
	nCaseTypes := 20
	nPartyTypes := 5

	// Make some CaseTypes
	caseTypes := s.mockCaseTypes(nCaseTypes)

	// Mase some PartyTypes
	partyTypes := []string{}
	for i := 0; i < nPartyTypes; i++ {
		partyTypes = append(partyTypes, newUUID())
	}

	// Prepare test data
	for i, caseType := range caseTypes {
		n := i % len(partyTypes)
		caseType.PartyTypeID = partyTypes[n]
		created, err := s.client.CaseTypes().Create(s.Ctx, caseType)
		assert.NoError(s.T(), err)
		caseType.ID = created.ID
	}

	s.Run("by party type", func() { s.testCaseTypeFilterByPartyType(caseTypes, partyTypes) })
}

func (s *Suite) testCaseTypeFilterByPartyType(caseTypes []*CaseType, partyTypes []string) {
	for i := 1; i <= len(partyTypes); i++ {
		types := partyTypes[0:i]
		list, err := s.client.CaseTypes().List(s.Ctx, CaseTypeListOptions{types})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		expected := []string{}
		for _, caseType := range caseTypes {
			if contains(types, caseType.PartyTypeID) {
				expected = append(expected, caseType.ID)
			}
		}
		assert.Len(s.T(), list.Items, len(expected))
		for _, item := range list.Items {
			assert.Contains(s.T(), expected, item.ID)
		}
	}
}
