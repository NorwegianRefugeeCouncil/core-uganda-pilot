package iam

import (
	"context"
	"github.com/nrc-no/core/internal/generic/server"
	"github.com/nrc-no/core/internal/rest"
	"path"
)

type RESTCountryClient struct {
	c *rest.Client
}

func (r RESTCountryClient) Get(ctx context.Context, id string) (*Country, error) {
	var obj Country
	err := r.c.Get().Path(path.Join(server.CountriesEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCountryClient) Create(ctx context.Context, create *Country) (*Country, error) {
	var obj Country
	err := r.c.Post().Body(create).Path(server.CountriesEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCountryClient) Update(ctx context.Context, update *Country) (*Country, error) {
	var obj Country
	err := r.c.Put().Body(update).Path(path.Join(server.CountriesEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCountryClient) List(ctx context.Context, listOptions CountryListOptions) (*CountryList, error) {
	var obj CountryList
	err := r.c.Get().Path(server.CountriesEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ CountryClient = &RESTCountryClient{}
