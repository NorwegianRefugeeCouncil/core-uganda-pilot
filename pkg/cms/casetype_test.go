package cms_test

import (
	. "github.com/nrc-no/core/pkg/cms"
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/i18n"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var FormWithAllControlTypes = form.Form{
	Controls: []form.Control{
		{
			Name:       "text",
			Type:       form.Text,
			Validation: form.ControlValidation{Required: true},
		},
		{
			Name:       "email",
			Type:       form.Email,
			Validation: form.ControlValidation{Required: true},
		},
		{
			Name:       "phone",
			Type:       form.Phone,
			Validation: form.ControlValidation{Required: true},
		},
		{
			Name:       "url",
			Type:       form.URL,
			Validation: form.ControlValidation{Required: true},
		},
		{
			Name:       "date",
			Type:       form.Date,
			Validation: form.ControlValidation{Required: true},
		},
		{
			Name:       "textarea",
			Type:       form.Textarea,
			Validation: form.ControlValidation{Required: true},
		},
		{
			Name: "dropdown",
			Type: form.Dropdown,
			Options: []i18n.Strings{
				{
					{"en", "a"},
				},
				{
					{"en", "b"},
				}},
			Validation: form.ControlValidation{Required: true},
		},
		{
			Name: "checkbox",
			Type: form.Checkbox,
			CheckboxOptions: []form.CheckboxOption{
				{Label: i18n.Strings{{"en", "a"}}},
				{Label: i18n.Strings{{"en", "b"}}},
			},
			Validation: form.ControlValidation{Required: true},
		},
		{
			Name: "radio",
			Type: form.Radio,
			CheckboxOptions: []form.CheckboxOption{
				{Label: i18n.Strings{{"en", "a"}}},
				{Label: i18n.Strings{{"en", "b"}}},
			},
			Validation: form.ControlValidation{Required: true},
		},
		{
			Name:       "taxonomy",
			Type:       form.Taxonomy,
			Validation: form.ControlValidation{Required: true},
		},
		{
			Name:       "file",
			Type:       form.File,
			Validation: form.ControlValidation{Required: true},
		},
	},
}

func (s *Suite) TestCaseType() {
	s.Run("API", func() { s.testCaseTypeAPI() })
	s.SetupTest()
	s.Run("List filtering", func() { s.testCaseTypeListFilter() })
}

func (s *Suite) testCaseTypeAPI() {
	// Create
	caseType := mockCaseType()
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
	caseType.Form = mockForm()
	updated, err := s.client.CaseTypes().Update(s.Ctx, caseType)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), caseType.Name, updated.Name)
	assert.Equal(s.T(), caseType.PartyTypeID, updated.PartyTypeID)
	assert.Equal(s.T(), caseType.TeamID, updated.TeamID)
	assert.Equal(s.T(), caseType.Form, updated.Form)

	// GET
	get, err = s.client.CaseTypes().Get(s.Ctx, updated.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), updated.Name, get.Name)
	assert.Equal(s.T(), updated.PartyTypeID, get.PartyTypeID)
	assert.Equal(s.T(), updated.TeamID, get.TeamID)
	assert.Equal(s.T(), updated.Form, get.Form)

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

	// Make some PartyTypes
	var partyTypes []string
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
		list, err := s.client.CaseTypes().List(s.Ctx, CaseTypeListOptions{PartyTypeIDs: types, TeamIDs: []string{}})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		var expected []string
		for _, caseType := range caseTypes {
			if utils.Contains(types, caseType.PartyTypeID) {
				expected = append(expected, caseType.ID)
			}
		}
		assert.Len(s.T(), list.Items, len(expected))
		for _, item := range list.Items {
			assert.Contains(s.T(), expected, item.ID)
		}
	}
}
