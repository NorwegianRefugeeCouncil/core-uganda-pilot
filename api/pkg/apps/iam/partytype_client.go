package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
)

type RESTPartyTypeClient struct {
	c *rest.Client
}

var partyTypesEP = server.Endpoints["partytypes"]

func (r RESTPartyTypeClient) Get(ctx context.Context, id string) (*PartyType, error) {
	var obj PartyType
	err := r.c.Get().Path(partyTypesEP + "/" + id).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyTypeClient) Create(ctx context.Context, create *PartyType) (*PartyType, error) {
	var obj PartyType
	err := r.c.Post().Body(create).Path(partyTypesEP).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyTypeClient) Update(ctx context.Context, update *PartyType) (*PartyType, error) {
	var obj PartyType
	err := r.c.Put().Body(update).Path(partyTypesEP + "/" + update.ID).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyTypeClient) List(ctx context.Context, listOptions PartyTypeListOptions) (*PartyTypeList, error) {
	var obj PartyTypeList
	err := r.c.Get().Path(partyTypesEP).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ PartyTypeClient = &RESTPartyTypeClient{}
