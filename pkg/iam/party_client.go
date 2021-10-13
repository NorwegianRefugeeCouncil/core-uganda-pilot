package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"path"
)

type RESTPartyClient struct {
	c *rest.Client
}

func (r RESTPartyClient) Get(ctx context.Context, id string) (*Party, error) {
	var obj Party
	err := r.c.Get().Path(path.Join(server.PartiesEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyClient) Create(ctx context.Context, create *Party) (*Party, error) {
	var obj Party
	err := r.c.Post().Body(create).Path(server.PartiesEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyClient) Update(ctx context.Context, update *Party) (*Party, error) {
	var obj Party
	err := r.c.Put().Body(update).Path(path.Join(server.PartiesEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyClient) List(ctx context.Context, listOptions PartyListOptions) (*PartyList, error) {
	var obj PartyList
	err := r.c.Get().Path(server.PartiesEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyClient) Search(ctx context.Context, listOptions PartySearchOptions) (*PartyList, error) {
	var obj PartyList
	err := r.c.Post().Path(path.Join(server.PartiesEndpoint, "/search")).Body(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ PartyClient = &RESTPartyClient{}