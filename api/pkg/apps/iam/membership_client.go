package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
)

type RESTMembershipClient struct {
	c *rest.Client
}

var membershipsEP = server.Endpoints["memberships"]

func (r RESTMembershipClient) Get(ctx context.Context, id string) (*Membership, error) {
	var obj Membership
	err := r.c.Get().Path(membershipsEP + "/" + id).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTMembershipClient) Create(ctx context.Context, create *Membership) (*Membership, error) {
	var obj Membership
	err := r.c.Post().Body(create).Path(membershipsEP).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTMembershipClient) Update(ctx context.Context, update *Membership) (*Membership, error) {
	var obj Membership
	err := r.c.Put().Body(update).Path(membershipsEP + "/" + update.ID).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTMembershipClient) List(ctx context.Context, listOptions MembershipListOptions) (*MembershipList, error) {
	var obj MembershipList
	err := r.c.Get().Path(membershipsEP).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ MembershipClient = &RESTMembershipClient{}
