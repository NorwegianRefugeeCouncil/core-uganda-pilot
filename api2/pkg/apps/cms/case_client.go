package cms

import (
	"context"
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/rest"
)

type RESTCaseClient struct {
	c *rest.Client
}

var _ CaseClient = &RESTCaseClient{}

func (r RESTCaseClient) Get(ctx context.Context, id string) (*Case, error) {
	var obj Case
	err := r.c.Get().Path(fmt.Sprintf("/apis/cms/v1/cases/%s", id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseClient) Create(ctx context.Context, create *Case) (*Case, error) {
	var obj Case
	err := r.c.Post().Body(create).Path("/apis/cms/v1/cases").Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseClient) Update(ctx context.Context, update *Case) (*Case, error) {
	var obj Case
	err := r.c.Put().Body(update).Path(fmt.Sprintf("/apis/cms/v1/cases/%s", update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseClient) List(ctx context.Context, listOptions CaseListOptions) (*CaseList, error) {
	var obj CaseList
	err := r.c.Get().Path("/apis/cms/v1/cases").WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}
