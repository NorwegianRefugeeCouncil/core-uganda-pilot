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

func (r RESTAttachmentClient) Get(ctx context.Context, id string) (*Attachment, error) {
	var obj Attachment
	err := r.c.Get().Path(path.Join(server.AttachmentsEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTAttachmentClient) List(ctx context.Context, listOptions AttachmentListOptions) (*AttachmentList, error) {
	var obj AttachmentList
	err := r.c.Get().Path(server.AttachmentsEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}
