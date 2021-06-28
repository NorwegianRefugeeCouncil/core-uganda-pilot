package iam

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/rest"
)

type RESTOrganizationClient struct {
	c *rest.Client
}

func (r RESTOrganizationClient) Get(ctx context.Context, id string) (*Organization, error) {
	var obj Organization
	err := r.c.Get().Path(fmt.Sprintf("/apis/iam/v1/organizations/%s", id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTOrganizationClient) Create(ctx context.Context, create *Organization) (*Organization, error) {
	var obj Organization
	err := r.c.Post().Body(create).Path("/apis/iam/v1/organizations").Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTOrganizationClient) Update(ctx context.Context, update *Organization) (*Organization, error) {
	var obj Organization
	err := r.c.Put().Body(update).Path(fmt.Sprintf("/apis/iam/v1/organizations/%s", update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTOrganizationClient) List(ctx context.Context, listOptions OrganizationListOptions) (*OrganizationList, error) {
	var obj OrganizationList
	err := r.c.Get().Path("/apis/iam/v1/organizations").WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ OrganizationClient = &RESTOrganizationClient{}
