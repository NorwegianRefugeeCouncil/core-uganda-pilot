package cms

import (
	"context"
	"github.com/nrc-no/core/internal/generic/server"
	"github.com/nrc-no/core/internal/rest"
	"path"
)

type RestCommentClient struct {
	c *rest.Client
}

var _ CaseTypeClient = &RESTCaseTypeClient{}

func (r RestCommentClient) Get(ctx context.Context, id string) (*Comment, error) {
	var obj Comment
	err := r.c.Get().Path(path.Join(server.CommentsEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RestCommentClient) Create(ctx context.Context, create *Comment) (*Comment, error) {
	var obj Comment
	err := r.c.Post().Body(create).Path(server.CommentsEndpoint).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RestCommentClient) Update(ctx context.Context, update *Comment) (*Comment, error) {
	var obj Comment
	err := r.c.Put().Body(update).Path(path.Join(server.CommentsEndpoint, update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RestCommentClient) List(ctx context.Context, listOptions CommentListOptions) (*CommentList, error) {
	var obj CommentList
	err := r.c.Get().Path(server.CommentsEndpoint).WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RestCommentClient) Delete(ctx context.Context, id string) error {
	return r.c.Delete().Path(path.Join(server.CommentsEndpoint, id)).Do(ctx).Error()
}
