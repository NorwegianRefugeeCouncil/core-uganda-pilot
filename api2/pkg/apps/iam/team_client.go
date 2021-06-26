package iam

import (
	"context"
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/rest"
)

type RESTTeamClient struct {
	c *rest.Client
}

func (r RESTTeamClient) Get(ctx context.Context, id string) (*Team, error) {
	var obj Team
	err := r.c.Get().Path(fmt.Sprintf("/apis/iam/v1/teams/%s", id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTTeamClient) Create(ctx context.Context, create *Team) (*Team, error) {
	var obj Team
	err := r.c.Post().Body(create).Path("/apis/iam/v1/teams").Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTTeamClient) Update(ctx context.Context, update *Team) (*Team, error) {
	var obj Team
	err := r.c.Put().Body(update).Path(fmt.Sprintf("/apis/iam/v1/teams/%s", update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTTeamClient) List(ctx context.Context, listOptions TeamListOptions) (*TeamList, error) {
	var obj TeamList
	err := r.c.Get().Path("/apis/iam/v1/teams").WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ TeamClient = &RESTTeamClient{}
