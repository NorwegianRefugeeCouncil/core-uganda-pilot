package cms

import (
	"context"
	"github.com/nrc-no/core/internal/generic/server"
	"github.com/nrc-no/core/internal/rest"
	"path"
)

type RESTCaseTypeClient struct {
	c *rest.Client
}

var _ CaseTypeClient = &RESTCaseTypeClient{}

func (r RESTCaseTypeClient) Get(ctx context.Context, id string) (*CaseType, error) {
	var obj CaseType
	err := r.c.Get().Path(path.Join(server.CaseTypesEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseTypeClient) Create(ctx context.Context, create *CaseType) (*CaseType, error) {
	var obj CaseType
	err := r.c.Post().Body(create).Path(server.CaseTypesEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseTypeClient) Update(ctx context.Context, update *CaseType) (*CaseType, error) {
	var obj CaseType
	err := r.c.Put().Body(update).Path(path.Join(server.CaseTypesEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseTypeClient) List(ctx context.Context, listOptions CaseTypeListOptions) (*CaseTypeList, error) {
	var obj CaseTypeList
	err := r.c.Get().Path(server.CaseTypesEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}
