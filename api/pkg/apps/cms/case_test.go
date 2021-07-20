package cms_test

import (
	"context"
	. "github.com/nrc-no/core/pkg/apps/cms"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func (s *Suite) TestCaseAPI() {
	// Create
	mockCase := aMockCase()
	created, err := s.client.Cases().Create(s.Ctx, mockCase)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	mockCase.ID = created.ID
	assert.Equal(s.T(), mockCase, created)

	// GET
	get, err := s.client.Cases().Get(s.Ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), created, get)

	// UPDATE
	get.Done = true
	get.FormData = &CaseTemplate{FormElements: []CaseTemplateFormElement{{
		Type: "textarea",
		Attributes: CaseTemplateFormElementAttribute{
			Label: "mock",
			Value: []string{"mock"},
		},
	}}}
	updated, err := s.client.Cases().Update(s.Ctx, get)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), get, updated)

	// GET
	get, err = s.client.Cases().Get(s.Ctx, updated.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.client.Cases().List(s.Ctx, CaseListOptions{})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Contains(s.T(), list.Items, get)
}

func (s *Suite) TestCaseListFilter() {

	bunch := newCaseBunch(20)
	if !assert.NoError(s.T(), bunch.create(s.Ctx, s.client.Cases())) {
		return
	}

	type tc struct {
		name    string
		args    CaseListOptions
		want    []*Case
		wantErr bool
	}

	tests := []tc{
		{
			name: "Single PartyID",
			args: CaseListOptions{
				PartyIDs: []string{
					bunch.parties[0],
				},
			},
			want: bunch.filterByPartyIDs(bunch.parties[0]).getCases(),
		},
		{
			name: "Multiple PartyIDs",
			args: CaseListOptions{
				PartyIDs: []string{
					bunch.parties[0],
					bunch.parties[1],
				},
			},
			want: bunch.filterByPartyIDs(bunch.parties[0], bunch.parties[1]).getCases(),
		},
	}

	for _, tt := range tests {
		tc := tt
		s.Run(tt.name, func() {
			s.T().Parallel()
			got, err := s.client.Cases().List(s.Ctx, tc.args)
			if tc.wantErr {
				s.Error(err)
				return
			} else {
				if !s.NoError(err) {
					return
				}
				s.ElementsMatch(tc.want, got.Items)
			}
		})
	}

}

type caseBunch struct {
	cases     []*Case
	parties   []string
	teams     []string
	caseTypes []string
	parents   []string
}

func newCaseBunch(nCases int) *caseBunch {

	var cases = mockCases(nCases)

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
	}

	return &caseBunch{
		cases:     cases,
		parties:   parties,
		teams:     teams,
		caseTypes: caseTypes,
		parents:   parents,
	}
}

func (c caseBunch) create(ctx context.Context, cli CaseClient) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, c := range c.cases {
		kase := c
		g.Go(func() error {
			created, err := cli.Create(ctx, kase)
			if err != nil {
				return err
			}
			kase.ID = created.ID
			return nil
		})
	}
	return g.Wait()
}

func (c caseBunch) filterByPartyIDs(partyID ...string) caseBunch {
	cb := caseBunch{}

	partyIDMap := map[string]bool{}
	for _, partyID := range partyID {
		partyIDMap[partyID] = true
	}

	for _, kase := range c.cases {
		if _, ok := partyIDMap[kase.PartyID]; ok {
			cb.cases = append(cb.cases, kase)
		}
	}
	return cb
}

func (c caseBunch) getCases() []*Case {
	return c.cases
}
