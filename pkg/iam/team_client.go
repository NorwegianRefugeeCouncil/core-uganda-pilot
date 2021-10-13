package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"path"
)

type RESTTeamClient struct {
	c *rest.Client
}

func (r RESTTeamClient) Get(ctx context.Context, id string) (*Team, error) {
	var obj Team
	err := r.c.Get().Path(path.Join(server.TeamsEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTTeamClient) Create(ctx context.Context, create *Team) (*Team, error) {
	var obj Team
	err := r.c.Post().Body(create).Path(server.TeamsEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTTeamClient) Update(ctx context.Context, update *Team) (*Team, error) {
	var obj Team
	err := r.c.Put().Body(update).Path(path.Join(server.TeamsEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTTeamClient) List(ctx context.Context, listOptions TeamListOptions) (*TeamList, error) {
	var obj TeamList
	err := r.c.Get().Path(server.TeamsEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ TeamClient = &RESTTeamClient{}