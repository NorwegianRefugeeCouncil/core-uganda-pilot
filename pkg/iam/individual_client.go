package iam

import (
	"context"
	"github.com/nrc-no/core/internal/generic/server"
	"github.com/nrc-no/core/internal/rest"
	"path"
)

type RESTIndividualClient struct {
	c *rest.Client
}

func (r RESTIndividualClient) Get(ctx context.Context, id string) (*Individual, error) {
	var obj Individual
	err := r.c.Get().Path(path.Join(server.IndividualsEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIndividualClient) Create(ctx context.Context, create *Individual) (*Individual, error) {
	var obj Individual
	err := r.c.Post().Body(create).Path(server.IndividualsEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIndividualClient) Update(ctx context.Context, update *Individual) (*Individual, error) {
	var obj Individual
	err := r.c.Put().Body(update).Path(path.Join(server.IndividualsEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIndividualClient) List(ctx context.Context, listOptions IndividualListOptions) (*IndividualList, error) {
	var obj IndividualList
	err := r.c.Get().Path(server.IndividualsEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ IndividualClient = &RESTIndividualClient{}
