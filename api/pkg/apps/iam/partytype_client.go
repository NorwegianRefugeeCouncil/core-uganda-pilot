package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"path"
)

type RESTPartyTypeClient struct {
	c *rest.Client
}

func (r RESTPartyTypeClient) Get(ctx context.Context, id string) (*PartyType, error) {
	var obj PartyType
	err := r.c.Get().Path(path.Join(server.PartyTypesEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyTypeClient) Create(ctx context.Context, create *PartyType) (*PartyType, error) {
	var obj PartyType
	err := r.c.Post().Body(create).Path(server.PartyTypesEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyTypeClient) Update(ctx context.Context, update *PartyType) (*PartyType, error) {
	var obj PartyType
	err := r.c.Put().Body(update).Path(path.Join(server.PartyTypesEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyTypeClient) List(ctx context.Context, listOptions PartyTypeListOptions) (*PartyTypeList, error) {
	var obj PartyTypeList
	err := r.c.Get().Path(server.PartyTypesEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ PartyTypeClient = &RESTPartyTypeClient{}
