package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"path"
)

type RESTIdentificationDocumentClient struct {
	c *rest.Client
}

func (r RESTIdentificationDocumentClient) Get(ctx context.Context, id string) (*IdentificationDocument, error) {
	var obj IdentificationDocument
	err := r.c.Get().Path(path.Join(server.IdentificationDocumentsEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIdentificationDocumentClient) Create(ctx context.Context, create *IdentificationDocument) (*IdentificationDocument, error) {
	var obj IdentificationDocument
	err := r.c.Post().Body(create).Path(server.IdentificationDocumentsEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIdentificationDocumentClient) Update(ctx context.Context, update *IdentificationDocument) (*IdentificationDocument, error) {
	var obj IdentificationDocument
	err := r.c.Put().Body(update).Path(path.Join(server.IdentificationDocumentsEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIdentificationDocumentClient) List(ctx context.Context, listOptions IdentificationDocumentListOptions) (*IdentificationDocumentList, error) {
	var obj IdentificationDocumentList
	err := r.c.Get().Path(server.IdentificationDocumentsEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTIdentificationDocumentClient) Delete(ctx context.Context, id string) error {
	return r.c.Delete().Path(server.IdentificationDocumentsEndpoint + "/" + id).Do(ctx).Error()
}

var _ IdentificationDocumentClient = &RESTIdentificationDocumentClient{}
