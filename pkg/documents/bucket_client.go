package documents

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"path"
)

type Buckets interface {
	Get(ctx context.Context, id string, options GetBucketOptions) (*Bucket, error)
	Delete(ctx context.Context, key string, options DeleteBucketOptions) error
	Create(ctx context.Context, obj *Bucket, options CreateBucketOptions) (*Bucket, error)
}

type RESTBucketClient struct {
	c *rest.Client
}

// NewBucketsClient returns a RESTBucketClient from a rest.Client
func NewBucketsClient(c *rest.Client) *RESTBucketClient {
	return &RESTBucketClient{
		c: c,
	}
}

// NewBucketsClientFromConfig returns a RESTBucketClient from a rest.Config
func NewBucketsClientFromConfig(restConfig *rest.Config) *RESTBucketClient {
	return NewBucketsClient(rest.NewClient(restConfig))
}

type GetBucketOptions struct {
}

// Get a Bucket
func (r RESTBucketClient) Get(ctx context.Context, id string, options GetBucketOptions) (*Bucket, error) {
	id = removeLeadingTrailingSlashes(id)
	var obj Bucket
	err := r.c.Get().Path(path.Join(server.BucketsEndpoint, id)).Do(ctx).Into(&obj)
	return &obj, err
}

type DeleteBucketOptions struct {
}

// Delete a bucket
func (r RESTBucketClient) Delete(ctx context.Context, key string, options DeleteBucketOptions) error {
	key = removeLeadingTrailingSlashes(key)
	return r.c.Delete().Path(path.Join(server.BucketsEndpoint, key)).Do(ctx).Into(nil)
}

type CreateBucketOptions struct {
}

// Create a bucket
func (r RESTBucketClient) Create(ctx context.Context, obj *Bucket, options CreateBucketOptions) (*Bucket, error) {
	var bucket Bucket
	err := r.c.Post().Body(obj).Path(server.BucketsEndpoint).Do(ctx).Into(&bucket)
	return &bucket, err
}
