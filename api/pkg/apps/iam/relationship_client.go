package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
)

type RESTRelationshipClient struct {
	c *rest.Client
}

var relationshipsEP = server.Endpoints["relationships"]

func (r RESTRelationshipClient) Get(ctx context.Context, id string) (*Relationship, error) {
	var obj Relationship
	err := r.c.Get().Path(relationshipsEP + "/" + id).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipClient) Create(ctx context.Context, create *Relationship) (*Relationship, error) {
	var obj Relationship
	err := r.c.Post().Body(create).Path(relationshipsEP).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipClient) Update(ctx context.Context, update *Relationship) (*Relationship, error) {
	var obj Relationship
	err := r.c.Put().Body(update).Path(relationshipsEP + "/" + update.ID).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipClient) List(ctx context.Context, listOptions RelationshipListOptions) (*RelationshipList, error) {
	var obj RelationshipList
	err := r.c.Get().Path(relationshipsEP).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipClient) Delete(ctx context.Context, id string) error {
	return r.c.Delete().Path(relationshipsEP + "/" + id).Do(ctx).Error()
}

var _ RelationshipClient = &RESTRelationshipClient{}
