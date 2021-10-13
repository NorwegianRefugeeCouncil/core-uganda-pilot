package cms

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"path"
)

type RESTCaseClient struct {
	c *rest.Client
}

var _ CaseClient = &RESTCaseClient{}

func (r RESTCaseClient) Get(ctx context.Context, id string) (*Case, error) {
	var obj Case
	err := r.c.Get().Path(path.Join(server.CasesEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseClient) Create(ctx context.Context, create *Case) (*Case, error) {
	var obj Case
	err := r.c.Post().Body(create).Path(server.CasesEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseClient) Update(ctx context.Context, update *Case) (*Case, error) {
	var obj Case
	err := r.c.Put().Body(update).Path(path.Join(server.CasesEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseClient) List(ctx context.Context, listOptions CaseListOptions) (*CaseList, error) {
	var obj CaseList
	err := r.c.Get().Path(server.CasesEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}
