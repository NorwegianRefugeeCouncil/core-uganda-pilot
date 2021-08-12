package cms_test

import (
	"context"
	. "github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/form"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"reflect"
)

func (s *Suite) TestCaseAPI() {
	// Create
	mockCase := aMockCase()
	created, err := s.client.Cases().Create(s.Ctx, mockCase)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	mockCase.ID = created.ID
	mockCase.CreatorID = created.CreatorID
	assert.Equal(s.T(), mockCase, created)

	// GET
	get, err := s.client.Cases().Get(s.Ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), created, get)

	// UPDATE
	get.Template = &CaseTemplate{FormElements: []form.FormElement{{
		Type: "textarea",
		Attributes: form.FormElementAttributes{
			Label: "update",
			Name:  "update",
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

	bunch := newCaseBunch()
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
				PartyIDs: bunch.parties[0:1],
			},
			want: bunch.filterBy("PartyID", bunch.parties[0:1]).getCases(),
		},
		{
			name: "Multiple PartyIDs",
			args: CaseListOptions{
				PartyIDs: bunch.parties[0:2],
			},
			want: bunch.filterBy("PartyID", bunch.parties[0:2]).getCases(),
		},
		{
			name: "Wrong PartyID",
			args: CaseListOptions{
				PartyIDs: []string{"abc"},
			},
			want:    bunch.filterBy("PartyID", []string{"abc"}).getCases(),
			wantErr: true,
		},
		{
			name: "Single TeamID",
			args: CaseListOptions{
				TeamIDs: bunch.teams[0:1],
			},
			want: bunch.filterBy("TeamID", bunch.teams[0:1]).getCases(),
		},
		{
			name: "Multiple TeamIDs",
			args: CaseListOptions{
				TeamIDs: bunch.teams[0:2],
			},
			want: bunch.filterBy("TeamID", bunch.teams[0:2]).getCases(),
		},
		{
			name: "Wrong TeamID",
			args: CaseListOptions{
				TeamIDs: []string{"abc"},
			},
			want:    bunch.filterBy("TeamID", []string{"abc"}).getCases(),
			wantErr: true,
		},
		{
			name: "Single CaseTypeID",
			args: CaseListOptions{
				CaseTypeIDs: bunch.caseTypes[0:1],
			},
			want: bunch.filterBy("CaseTypeID", bunch.caseTypes[0:1]).getCases(),
		},
		{
			name: "Multiple CaseTypeIDs",
			args: CaseListOptions{
				CaseTypeIDs: bunch.caseTypes[0:2],
			},
			want: bunch.filterBy("CaseTypeID", bunch.caseTypes[0:2]).getCases(),
		},
		{
			name: "Wrong CaseTypeID",
			args: CaseListOptions{
				CaseTypeIDs: []string{"abc"},
			},
			want:    bunch.filterBy("CaseTypeID", []string{"abc"}).getCases(),
			wantErr: true,
		},
		{
			name: "Single ParentID",
			args: CaseListOptions{
				ParentID: bunch.parents[0],
			},
			want: bunch.filterBy("ParentID", bunch.parents[0:1]).getCases(),
		},
		{
			name: "Wrong ParentID",
			args: CaseListOptions{
				ParentID: "abc",
			},
			want:    bunch.filterBy("ParentID", []string{"abc"}).getCases(),
			wantErr: true,
		},
		{
			name: "Done",
			args: CaseListOptions{
				Done: boolPtr(true),
			},
			want: bunch.filterByDone(true).getCases(),
		},
		{
			name: "Done",
			args: CaseListOptions{
				Done: boolPtr(false),
			},
			want: bunch.filterByDone(false).getCases(),
		},
		{
			name: "Combo PartyIDs TeamIDs",
			args: CaseListOptions{
				PartyIDs: bunch.parties[0:2],
				TeamIDs:  bunch.teams[0:2],
			},
			want: bunch.filterBy("PartyID", bunch.parties[0:2]).filterBy("TeamID", bunch.teams[0:2]).getCases(),
		},
		{
			name: "Combo PartyIDs TeamIDs CaseTypes",
			args: CaseListOptions{
				PartyIDs:    bunch.parties[0:2],
				TeamIDs:     bunch.teams[0:2],
				CaseTypeIDs: bunch.caseTypes[0:2],
			},
			want: bunch.filterBy("PartyID", bunch.parties[0:2]).filterBy("TeamID", bunch.teams[0:2]).filterBy("CaseTypeID", bunch.caseTypes[0:2]).getCases(),
		},
		{
			name: "Combo PartyIDs TeamIDs CaseTypes Parent",
			args: CaseListOptions{
				PartyIDs:    bunch.parties[0:2],
				TeamIDs:     bunch.teams[0:2],
				CaseTypeIDs: bunch.caseTypes[0:2],
				ParentID:    bunch.parents[0],
			},
			want: bunch.filterBy("PartyID", bunch.parties[0:2]).filterBy("TeamID", bunch.teams[0:2]).filterBy("CaseTypeID", bunch.caseTypes[0:2]).filterBy("ParentID", bunch.parents[0:1]).getCases(),
		},
		{
			name: "Combo PartyIDs TeamIDs CaseTypes Parent Done",
			args: CaseListOptions{
				PartyIDs:    bunch.parties[0:2],
				TeamIDs:     bunch.teams[0:2],
				CaseTypeIDs: bunch.caseTypes[0:2],
				ParentID:    bunch.parents[0],
				Done:        boolPtr(true),
			},
			want: bunch.filterBy("PartyID", bunch.parties[0:2]).filterBy("TeamID", bunch.teams[0:2]).filterBy("CaseTypeID", bunch.caseTypes[0:2]).filterBy("ParentID", bunch.parents[0:1]).filterByDone(true).getCases(),
		},
	}

	for _, tt := range tests {
		tc := tt
		s.Run(tt.name, func() {
			//s.T().Parallel()
			got, err := s.client.Cases().List(s.Ctx, tc.args)
			if !s.NoError(err) {
				s.T().FailNow()
			}
			if tc.wantErr {
				s.Empty(got.Items)
			} else {
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

func newCaseBunch() *caseBunch {

	const (
		nCases     = 6
		nParties   = 5
		nTeams     = 4
		nCaseTypes = 5
		nParents   = 10
	)

	cases := mockCases(nCases)
	parties := uuidSlice(nParties)
	teams := uuidSlice(nTeams)
	caseTypes := uuidSlice(nCaseTypes)
	parents := uuidSlice(nParents)

	// Prepare test data
	for i, c := range cases {
		c.PartyID = parties[i%len(parties)]
		c.TeamID = teams[i%len(teams)]
		c.CaseTypeID = caseTypes[i%len(caseTypes)]
		c.ParentID = parents[i%len(parents)]
		c.Done = i%2 == 0
	}

	return &caseBunch{
		cases,
		parties,
		teams,
		caseTypes,
		parents,
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
			kase.CreatorID = created.CreatorID
			return nil
		})
	}
	return g.Wait()
}

func (c caseBunch) filterBy(fieldName string, uuids []string) caseBunch {
	cb := caseBunch{}
	idMap := map[string]bool{}
	for _, uuid := range uuids {
		idMap[uuid] = true
	}
	for _, kase := range c.cases {
		k := reflect.ValueOf(kase).Elem().FieldByName(fieldName).String()
		if _, ok := idMap[k]; ok {
			cb.cases = append(cb.cases, kase)
		}
	}
	return cb
}

func (c caseBunch) filterByDone(done bool) caseBunch {
	cb := caseBunch{}
	for _, kase := range c.cases {
		if kase.Done == done {
			cb.cases = append(cb.cases, kase)
		}
	}
	return cb
}

func (c caseBunch) getCases() []*Case {
	return c.cases
}

func boolPtr(b bool) *bool {
	return &b
}
