package documents

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"path"
)

type Buckets interface {
	Get(ctx context.Context, id string, options GetBucketOptions) (*Bucket, error)
	Delete(ctx context.Context, key string) error
	Create(ctx context.Context, obj *Bucket) (*Bucket, error)
}

type RESTBucketClient struct {
	c *rest.Client
}

func NewBucketsClient(c *rest.Client) *RESTBucketClient {
	return &RESTBucketClient{
		c: c,
	}
}

func NewBucketsClientFromConfig(restConfig *rest.Config) *RESTBucketClient {
	return NewBucketsClient(rest.NewClient(restConfig))
}

type GetBucketOptions struct {
}

func (r RESTBucketClient) Get(ctx context.Context, id string, options GetBucketOptions) (*Bucket, error) {
	id = normaliseKey(id)
	var obj Bucket
	err := r.c.Get().Path(path.Join(server.BucketsEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

func (r RESTBucketClient) Delete(ctx context.Context, key string) error {
	key = normaliseKey(key)
	return r.c.Delete().Path(path.Join(server.BucketsEndpoint, key)).Do(ctx).Into(nil)
}

type CreateBucketResponse struct {
}

func (r RESTBucketClient) Create(ctx context.Context, obj *Bucket) (*Bucket, error) {
	var bucket Bucket
	err := r.c.Post().Body(obj).Path(server.BucketsEndpoint).Do(ctx).Into(&bucket)
	return &bucket, err
}
