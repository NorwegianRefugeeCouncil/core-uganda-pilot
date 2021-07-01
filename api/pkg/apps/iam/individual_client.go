package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
)

type RESTIndividualClient struct {
	c *rest.Client
}

var individualsEP = server.Endpoints["individuals"]

func (r RESTIndividualClient) Get(ctx context.Context, id string) (*Individual, error) {
	var obj Individual
	err := r.c.Get().Path(individualsEP + "/" + id).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIndividualClient) Create(ctx context.Context, create *Individual) (*Individual, error) {
	var obj Individual
	err := r.c.Post().Body(create).Path(individualsEP).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIndividualClient) Update(ctx context.Context, update *Individual) (*Individual, error) {
	var obj Individual
	err := r.c.Put().Body(update).Path(individualsEP + "/" + update.ID).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIndividualClient) List(ctx context.Context, listOptions IndividualListOptions) (*IndividualList, error) {
	var obj IndividualList
	err := r.c.Get().Path(individualsEP).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ IndividualClient = &RESTIndividualClient{}
