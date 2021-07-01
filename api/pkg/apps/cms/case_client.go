package cms

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
)

type RESTCaseClient struct {
	c *rest.Client
}

var casesEP = server.Endpoints["cases"]

var _ CaseClient = &RESTCaseClient{}

func (r RESTCaseClient) Get(ctx context.Context, id string) (*Case, error) {
	var obj Case
	err := r.c.Get().Path(casesEP + "/" + id).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseClient) Create(ctx context.Context, create *Case) (*Case, error) {
	var obj Case
	err := r.c.Post().Body(create).Path("/apis/cms/v1/cases").Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseClient) Update(ctx context.Context, update *Case) (*Case, error) {
	var obj Case
	err := r.c.Put().Body(update).Path(casesEP + "/" + update.ID).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseClient) List(ctx context.Context, listOptions CaseListOptions) (*CaseList, error) {
	var obj CaseList
	err := r.c.Get().Path(casesEP).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}
