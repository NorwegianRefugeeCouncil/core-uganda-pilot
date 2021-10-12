package attachments

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"path"
)

type RESTAttachmentClient struct {
	c *rest.Client
}

func NewClient(restConfig *rest.Config) *RESTAttachmentClient {
	return &RESTAttachmentClient{
		c: rest.NewClient(restConfig),
	}
}

func (r RESTAttachmentClient) Get(ctx context.Context, id string) (*Attachment, error) {
	var obj Attachment
	err := r.c.Get().Path(path.Join(server.AttachmentsEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTAttachmentClient) Create(ctx context.Context, create *Attachment) (*Attachment, error) {
	var obj Attachment
	err := r.c.Post().Body(create).Path(server.AttachmentsEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTAttachmentClient) Update(ctx context.Context, update *Attachment) (*Attachment, error) {
	var obj Attachment
	err := r.c.Put().Body(update).Path(path.Join(server.AttachmentsEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTAttachmentClient) List(ctx context.Context, listOptions AttachmentListOptions) (*AttachmentList, error) {
	var obj AttachmentList
	err := r.c.Get().Path(server.AttachmentsEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}
