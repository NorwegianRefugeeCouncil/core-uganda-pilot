package iam

import (
	"context"
	"github.com/nrc-no/core/internal/generic/server"
	"github.com/nrc-no/core/internal/rest"
	"path"
)

type RESTAttributeClient struct {
	c *rest.Client
}

func (r RESTAttributeClient) Get(ctx context.Context, id string) (*PartyAttributeDefinition, error) {
	var obj PartyAttributeDefinition
	err := r.c.Get().Path(path.Join(server.AttributesEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTAttributeClient) Create(ctx context.Context, create *PartyAttributeDefinition) (*PartyAttributeDefinition, error) {
	var obj PartyAttributeDefinition
	err := r.c.Post().Body(create).Path(server.AttributesEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTAttributeClient) Update(ctx context.Context, update *PartyAttributeDefinition) (*PartyAttributeDefinition, error) {
	var obj PartyAttributeDefinition
	err := r.c.Put().Body(update).Path(path.Join(server.AttributesEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTAttributeClient) List(ctx context.Context, listOptions PartyAttributeDefinitionListOptions) (*PartyAttributeDefinitionList, error) {
	var obj PartyAttributeDefinitionList
	err := r.c.Get().Path(server.AttributesEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ PartyAttributeDefinitionClient = &RESTAttributeClient{}
