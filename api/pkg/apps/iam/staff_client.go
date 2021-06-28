package iam

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/rest"
)

type RESTStaffClient struct {
	c *rest.Client
}

func (r RESTStaffClient) Get(ctx context.Context, id string) (*Staff, error) {
	var obj Staff
	err := r.c.Get().Path(fmt.Sprintf("/apis/iam/v1/staff/%s", id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTStaffClient) Create(ctx context.Context, create *Staff) (*Staff, error) {
	var obj Staff
	err := r.c.Post().Body(create).Path("/apis/iam/v1/staff").Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTStaffClient) Update(ctx context.Context, update *Staff) (*Staff, error) {
	var obj Staff
	err := r.c.Put().Body(update).Path(fmt.Sprintf("/apis/iam/v1/staff/%s", update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTStaffClient) List(ctx context.Context, listOptions StaffListOptions) (*StaffList, error) {
	var obj StaffList
	err := r.c.Get().Path("/apis/iam/v1/staff").WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ StaffClient = &RESTStaffClient{}
