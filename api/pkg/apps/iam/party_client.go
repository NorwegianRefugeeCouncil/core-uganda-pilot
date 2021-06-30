package iam

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/rest"
)

type RESTPartyClient struct {
	c *rest.Client
}

func (r RESTPartyClient) Get(ctx context.Context, id string) (*Party, error) {
	var obj Party
	err := r.c.Get().Path(fmt.Sprintf("/apis/iam/v1/parties/%s", id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyClient) Create(ctx context.Context, create *Party) (*Party, error) {
	var obj Party
	err := r.c.Post().Body(create).Path("/apis/iam/v1/parties").Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyClient) Update(ctx context.Context, update *Party) (*Party, error) {
	var obj Party
	err := r.c.Put().Body(update).Path(fmt.Sprintf("/apis/iam/v1/parties/%s", update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyClient) List(ctx context.Context, listOptions PartyListOptions) (*PartyList, error) {
	var obj PartyList
	err := r.c.Get().Path("/apis/iam/v1/parties").WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTPartyClient) Search(ctx context.Context, listOptions PartySearchOptions) (*PartyList, error) {
	var obj PartyList
	err := r.c.Post().Path("/apis/iam/v1/parties/search").Body(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ PartyClient = &RESTPartyClient{}
