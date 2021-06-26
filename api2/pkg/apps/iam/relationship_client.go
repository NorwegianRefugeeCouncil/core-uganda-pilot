package iam

import (
	"context"
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/rest"
)

type RESTRelationshipClient struct {
	c *rest.Client
}

func (r RESTRelationshipClient) Get(ctx context.Context, id string) (*Relationship, error) {
	var obj Relationship
	err := r.c.Get().Path(fmt.Sprintf("/apis/iam/v1/relationships/%s", id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipClient) Create(ctx context.Context, create *Relationship) (*Relationship, error) {
	var obj Relationship
	err := r.c.Post().Body(create).Path("/apis/iam/v1/relationships").Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipClient) Update(ctx context.Context, update *Relationship) (*Relationship, error) {
	var obj Relationship
	err := r.c.Put().Body(update).Path(fmt.Sprintf("/apis/iam/v1/relationships/%s", update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipClient) List(ctx context.Context, listOptions RelationshipListOptions) (*RelationshipList, error) {
	var obj RelationshipList
	err := r.c.Get().Path("/apis/iam/v1/relationships").WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipClient) Delete(ctx context.Context, id string) error {
	return r.c.Delete().Path(fmt.Sprintf("/apis/iam/v1/relationships/%s", id)).Do(ctx).Error()
}

var _ RelationshipClient = &RESTRelationshipClient{}
