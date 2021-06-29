package cms

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/rest"
)

type RestCommentClient struct {
	c *rest.Client
}

var _ CaseTypeClient = &RESTCaseTypeClient{}

func (r RestCommentClient) Get(ctx context.Context, id string) (*Comment, error) {
	var obj Comment
	err := r.c.Get().Path(fmt.Sprintf("/apis/cms/v1/comments/%s", id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RestCommentClient) Create(ctx context.Context, create *Comment) (*Comment, error) {
	var obj Comment
	err := r.c.Post().Body(create).Path("/apis/cms/v1/comments").Do(ctx).Into(&obj)
	return &obj, err
}

func (r RestCommentClient) Update(ctx context.Context, update *Comment) (*Comment, error) {
	var obj Comment
	err := r.c.Put().Body(update).Path(fmt.Sprintf("/apis/cms/v1/comments/%s", update.ID)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RestCommentClient) List(ctx context.Context, listOptions CommentListOptions) (*CommentList, error) {
	var obj CommentList
	err := r.c.Get().Path("/apis/cms/v1/comments").WithParams(listOptions).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RestCommentClient) Delete(ctx context.Context, id string) error {
	return r.c.Delete().Path(fmt.Sprintf("/apis/cms/v1/comments/%s", id)).Do(ctx).Error()
}
