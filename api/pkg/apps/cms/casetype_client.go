package cms

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
)

type RESTCaseTypeClient struct {
	c *rest.Client
}

var caseTypesEP = server.Endpoints["casetypes"]

var _ CaseTypeClient = &RESTCaseTypeClient{}

func (r RESTCaseTypeClient) Get(ctx context.Context, id string) (*CaseType, error) {
	var obj CaseType
	err := r.c.Get().Path(fmt.Sprintf(caseTypesEP + id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseTypeClient) Create(ctx context.Context, create *CaseType) (*CaseType, error) {
	var obj CaseType
	err := r.c.Post().Body(create).Path(caseTypesEP).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseTypeClient) Update(ctx context.Context, update *CaseType) (*CaseType, error) {
	var obj CaseType
	err := r.c.Put().Body(update).Path(fmt.Sprintf(caseTypesEP + update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTCaseTypeClient) List(ctx context.Context, listOptions CaseTypeListOptions) (*CaseTypeList, error) {
	var obj CaseTypeList
	err := r.c.Get().Path(caseTypesEP).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}
