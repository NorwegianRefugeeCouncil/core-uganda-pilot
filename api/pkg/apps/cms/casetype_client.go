package cms

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/rest"
)

type RESTCaseTypeClient struct {
	c *rest.Client
}

var _ CaseTypeClient = &RESTCaseTypeClient{}

func (r RESTCaseTypeClient) Get(ctx context.Context, id string) (*CaseType, error) {
	var obj CaseType
	err := r.c.Get().Path(fmt.Sprintf("/apis/cms/v1/casetypes/%s", id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseTypeClient) Create(ctx context.Context, create *CaseType) (*CaseType, error) {
	var obj CaseType
	err := r.c.Post().Body(create).Path("/apis/cms/v1/casetypes").Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseTypeClient) Update(ctx context.Context, update *CaseType) (*CaseType, error) {
	var obj CaseType
	err := r.c.Put().Body(update).Path(fmt.Sprintf("/apis/cms/v1/casetypes/%s", update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseTypeClient) List(ctx context.Context, listOptions CaseTypeListOptions) (*CaseTypeList, error) {
	var obj CaseTypeList
	err := r.c.Get().Path("/apis/cms/v1/casetypes").WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}
