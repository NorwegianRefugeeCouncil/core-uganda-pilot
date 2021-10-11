package iam

import (
	"context"
	"github.com/nrc-no/core/internal/generic/server"
	"github.com/nrc-no/core/internal/rest"
	"path"
)

type RESTRelationshipTypeClient struct {
	c *rest.Client
}

func (r RESTRelationshipTypeClient) Get(ctx context.Context, id string) (*RelationshipType, error) {
	var obj RelationshipType
	err := r.c.Get().Path(path.Join(server.RelationshipTypesEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipTypeClient) Create(ctx context.Context, create *RelationshipType) (*RelationshipType, error) {
	var obj RelationshipType
	err := r.c.Post().Body(create).Path(server.RelationshipTypesEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipTypeClient) Update(ctx context.Context, update *RelationshipType) (*RelationshipType, error) {
	var obj RelationshipType
	err := r.c.Put().Body(update).Path(path.Join(server.RelationshipTypesEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipTypeClient) List(ctx context.Context, listOptions RelationshipTypeListOptions) (*RelationshipTypeList, error) {
	var obj RelationshipTypeList
	err := r.c.Get().Path(server.RelationshipTypesEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ RelationshipTypeClient = &RESTRelationshipTypeClient{}
