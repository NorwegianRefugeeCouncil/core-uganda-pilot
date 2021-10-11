package iam

import (
	"context"
	"github.com/nrc-no/core/internal/generic/server"
	"github.com/nrc-no/core/internal/rest"
	"path"
)

type RESTMembershipClient struct {
	c *rest.Client
}

func (r RESTMembershipClient) Get(ctx context.Context, id string) (*Membership, error) {
	var obj Membership
	err := r.c.Get().Path(path.Join(server.MembershipsEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTMembershipClient) Create(ctx context.Context, create *Membership) (*Membership, error) {
	var obj Membership
	err := r.c.Post().Body(create).Path(server.MembershipsEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTMembershipClient) Update(ctx context.Context, update *Membership) (*Membership, error) {
	var obj Membership
	err := r.c.Put().Body(update).Path(path.Join(server.MembershipsEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTMembershipClient) List(ctx context.Context, listOptions MembershipListOptions) (*MembershipList, error) {
	var obj MembershipList
	err := r.c.Get().Path(server.MembershipsEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ MembershipClient = &RESTMembershipClient{}
