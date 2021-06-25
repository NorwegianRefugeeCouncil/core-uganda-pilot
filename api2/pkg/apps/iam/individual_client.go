package iam

import (
	"context"
	"fmt"
)

type RESTIndividualClient struct {
	c *RESTClient
}

func (r RESTIndividualClient) Get(ctx context.Context, id string) (*Individual, error) {
	var obj Individual
	err := r.c.Get().Path(fmt.Sprintf("/apis/iam/v1/individuals/%s", id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIndividualClient) Create(ctx context.Context, create *Individual) (*Individual, error) {
	var obj Individual
	err := r.c.Post().Body(create).Path("/apis/iam/v1/individuals").Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIndividualClient) Update(ctx context.Context, update *Individual) (*Individual, error) {
	var obj Individual
	err := r.c.Put().Body(update).Path(fmt.Sprintf("/apis/iam/v1/individuals/%s", update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIndividualClient) List(ctx context.Context, listOptions IndividualListOptions) (*IndividualList, error) {
	var obj IndividualList
	err := r.c.Get().Path("/apis/iam/v1/individuals").WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ IndividualClient = &RESTIndividualClient{}
