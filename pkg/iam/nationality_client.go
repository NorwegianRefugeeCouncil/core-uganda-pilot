package iam

import (
	"context"
	"github.com/nrc-no/core/internal/generic/server"
	"github.com/nrc-no/core/internal/rest"
	"path"
)

type RESTNationalityClient struct {
	c *rest.Client
}

func (r RESTNationalityClient) Get(ctx context.Context, id string) (*Nationality, error) {
	var obj Nationality
	err := r.c.Get().Path(path.Join(server.NationalitiesEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTNationalityClient) Create(ctx context.Context, create *Nationality) (*Nationality, error) {
	var obj Nationality
	err := r.c.Post().Body(create).Path(server.NationalitiesEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTNationalityClient) Update(ctx context.Context, update *Nationality) (*Nationality, error) {
	var obj Nationality
	err := r.c.Put().Body(update).Path(path.Join(server.NationalitiesEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTNationalityClient) List(ctx context.Context, listOptions NationalityListOptions) (*NationalityList, error) {
	var obj NationalityList
	err := r.c.Get().Path(server.NationalitiesEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ NationalityClient = &RESTNationalityClient{}
