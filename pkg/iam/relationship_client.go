package iam

import (
	"context"
	"github.com/nrc-no/core/internal/generic/server"
	"github.com/nrc-no/core/internal/rest"
	"path"
)

type RESTRelationshipClient struct {
	c *rest.Client
}

func (r RESTRelationshipClient) Get(ctx context.Context, id string) (*Relationship, error) {
	var obj Relationship
	err := r.c.Get().Path(path.Join(server.RelationshipsEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipClient) Create(ctx context.Context, create *Relationship) (*Relationship, error) {
	var obj Relationship
	err := r.c.Post().Body(create).Path(server.RelationshipsEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipClient) Update(ctx context.Context, update *Relationship) (*Relationship, error) {
	var obj Relationship
	err := r.c.Put().Body(update).Path(path.Join(server.RelationshipsEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipClient) List(ctx context.Context, listOptions RelationshipListOptions) (*RelationshipList, error) {
	var obj RelationshipList
	err := r.c.Get().Path(server.RelationshipsEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipClient) Delete(ctx context.Context, id string) error {
	return r.c.Delete().Path(server.RelationshipsEndpoint + "/" + id).Do(ctx).Error()
}

var _ RelationshipClient = &RESTRelationshipClient{}
