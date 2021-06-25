package iam

import (
	"context"
	"fmt"
)

type RESTAttributeClient struct {
	c *RESTClient
}

func (r RESTAttributeClient) Get(ctx context.Context, id string) (*Attribute, error) {
	var obj Attribute
	err := r.c.Get().Path(fmt.Sprintf("/apis/iam/v1/attributes/%s", id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTAttributeClient) Create(ctx context.Context, create *Attribute) (*Attribute, error) {
	var obj Attribute
	err := r.c.Post().Body(create).Path("/apis/iam/v1/attributes").Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTAttributeClient) Update(ctx context.Context, update *Attribute) (*Attribute, error) {
	var obj Attribute
	err := r.c.Put().Body(update).Path(fmt.Sprintf("/apis/iam/v1/attributes/%s", update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTAttributeClient) List(ctx context.Context, listOptions AttributeListOptions) (*AttributeList, error) {
	var obj AttributeList
	err := r.c.Get().Path("/apis/iam/v1/attributes").WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ AttributeClient = &RESTAttributeClient{}
