package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"path"
)

type RESTIdentificationDocumentTypeClient struct {
	c *rest.Client
}

func (r RESTIdentificationDocumentTypeClient) Get(ctx context.Context, id string) (*IdentificationDocumentType, error) {
	var obj IdentificationDocumentType
	err := r.c.Get().Path(path.Join(server.IdentificationDocumentTypesEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIdentificationDocumentTypeClient) Create(ctx context.Context, create *IdentificationDocumentType) (*IdentificationDocumentType, error) {
	var obj IdentificationDocumentType
	err := r.c.Post().Body(create).Path(server.IdentificationDocumentTypesEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIdentificationDocumentTypeClient) Update(ctx context.Context, update *IdentificationDocumentType) (*IdentificationDocumentType, error) {
	var obj IdentificationDocumentType
	err := r.c.Put().Body(update).Path(path.Join(server.IdentificationDocumentTypesEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIdentificationDocumentTypeClient) List(ctx context.Context, listOptions IdentificationDocumentTypeListOptions) (*IdentificationDocumentTypeList, error) {
	var obj IdentificationDocumentTypeList
	err := r.c.Get().Path(server.IdentificationDocumentTypesEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

var _ IdentificationDocumentTypeClient = &RESTIdentificationDocumentTypeClient{}
