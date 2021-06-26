package iam

import (
	"context"
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/rest"
)

type RESTPartyTypeClient struct {
	c *rest.Client
}

func (r RESTPartyTypeClient) Get(ctx context.Context, id string) (*PartyType, error) {
	var obj PartyType
	err := r.c.Get().Path(fmt.Sprintf("/apis/iam/v1/partytypes/%s", id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyTypeClient) Create(ctx context.Context, create *PartyType) (*PartyType, error) {
	var obj PartyType
	err := r.c.Post().Body(create).Path("/apis/iam/v1/partytypes").Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyTypeClient) Update(ctx context.Context, update *PartyType) (*PartyType, error) {
	var obj PartyType
	err := r.c.Put().Body(update).Path(fmt.Sprintf("/apis/iam/v1/partytypes/%s", update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyTypeClient) List(ctx context.Context, listOptions PartyTypeListOptions) (*PartyTypeList, error) {
	var obj PartyTypeList
	err := r.c.Get().Path("/apis/iam/v1/partytypes").WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ PartyTypeClient = &RESTPartyTypeClient{}
