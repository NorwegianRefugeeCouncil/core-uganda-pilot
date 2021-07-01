package iam

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
)

type RESTAttributeClient struct {
	c *rest.Client
}

var attributesEP = server.Endpoints["attributes"]

func (r RESTAttributeClient) Get(ctx context.Context, id string) (*Attribute, error) {
	var obj Attribute
	err := r.c.Get().Path(fmt.Sprintf(attributesEP + id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTAttributeClient) Create(ctx context.Context, create *Attribute) (*Attribute, error) {
	var obj Attribute
	err := r.c.Post().Body(create).Path(attributesEP).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTAttributeClient) Update(ctx context.Context, update *Attribute) (*Attribute, error) {
	var obj Attribute
	err := r.c.Put().Body(update).Path(fmt.Sprintf(attributesEP + update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTAttributeClient) List(ctx context.Context, listOptions AttributeListOptions) (*AttributeList, error) {
	var obj AttributeList
	err := r.c.Get().Path(attributesEP).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ AttributeClient = &RESTAttributeClient{}
