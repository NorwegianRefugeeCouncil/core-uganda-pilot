package iam

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/rest"
)

type RESTMembershipClient struct {
	c *rest.Client
}

func (r RESTMembershipClient) Get(ctx context.Context, id string) (*Membership, error) {
	var obj Membership
	err := r.c.Get().Path(fmt.Sprintf("/apis/iam/v1/memberships/%s", id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTMembershipClient) Create(ctx context.Context, create *Membership) (*Membership, error) {
	var obj Membership
	err := r.c.Post().Body(create).Path("/apis/iam/v1/memberships").Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTMembershipClient) Update(ctx context.Context, update *Membership) (*Membership, error) {
	var obj Membership
	err := r.c.Put().Body(update).Path(fmt.Sprintf("/apis/iam/v1/memberships/%s", update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTMembershipClient) List(ctx context.Context, listOptions MembershipListOptions) (*MembershipList, error) {
	var obj MembershipList
	err := r.c.Get().Path("/apis/iam/v1/memberships").WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ MembershipClient = &RESTMembershipClient{}
