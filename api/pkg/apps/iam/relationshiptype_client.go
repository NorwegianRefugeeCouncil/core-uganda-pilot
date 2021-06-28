package iam

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/rest"
)

type RESTRelationshipTypeClient struct {
	c *rest.Client
}

func (r RESTRelationshipTypeClient) Get(ctx context.Context, id string) (*RelationshipType, error) {
	var obj RelationshipType
	err := r.c.Get().Path(fmt.Sprintf("/apis/iam/v1/relationshiptypes/%s", id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipTypeClient) Create(ctx context.Context, create *RelationshipType) (*RelationshipType, error) {
	var obj RelationshipType
	err := r.c.Post().Body(create).Path("/apis/iam/v1/relationshiptypes").Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipTypeClient) Update(ctx context.Context, update *RelationshipType) (*RelationshipType, error) {
	var obj RelationshipType
	err := r.c.Put().Body(update).Path(fmt.Sprintf("/apis/iam/v1/relationshiptypes/%s", update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTRelationshipTypeClient) List(ctx context.Context, listOptions RelationshipTypeListOptions) (*RelationshipTypeList, error) {
	var obj RelationshipTypeList
	err := r.c.Get().Path("/apis/iam/v1/relationshiptypes").WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ RelationshipTypeClient = &RESTRelationshipTypeClient{}
