package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
)

type RESTTeamClient struct {
	c *rest.Client
}

var teamsEP = server.Endpoints["teams"]

func (r RESTTeamClient) Get(ctx context.Context, id string) (*Team, error) {
	var obj Team
	err := r.c.Get().Path(teamsEP + "/" + id).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTTeamClient) Create(ctx context.Context, create *Team) (*Team, error) {
	var obj Team
	err := r.c.Post().Body(create).Path(teamsEP).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTTeamClient) Update(ctx context.Context, update *Team) (*Team, error) {
	var obj Team
	err := r.c.Put().Body(update).Path(teamsEP + "/" + update.ID).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTTeamClient) List(ctx context.Context, listOptions TeamListOptions) (*TeamList, error) {
	var obj TeamList
	err := r.c.Get().Path(teamsEP).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ TeamClient = &RESTTeamClient{}
