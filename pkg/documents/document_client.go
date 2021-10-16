package documents

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"net/url"
	"path"
	"strconv"
)

type Documents interface {
	Get(ctx context.Context, id string, options GetDocumentOptions) (*Document, error)
	Put(ctx context.Context, document *Document, options PutDocumentOptions) (*PutDocumentResponse, error)
	Delete(ctx context.Context, key string, options DeleteDocumentOptions) error
}

type RESTDocumentClient struct {
	c *rest.Client
}

func NewDocumentsClient(c *rest.Client) *RESTDocumentClient {
	return &RESTDocumentClient{
		c: c,
	}
}

func NewDocumentsClientFromConfig(restConfig *rest.Config) *RESTDocumentClient {
	return NewDocumentsClient(rest.NewClient(restConfig))
}

type GetDocumentOptions struct {
	BucketID string
	Version  string
}

// Get a document
func (r RESTDocumentClient) Get(ctx context.Context, id string, options GetDocumentOptions) (*Document, error) {

	id = normaliseKey(id)

	var obj Document
	body, resp, err := r.c.Get().Path(path.Join(server.DocumentsEndpoint, id)).
		WithParams(url.Values{
			"bucketId": []string{options.BucketID},
			"version":  []string{options.Version},
		}).
		Do(ctx).Raw()
	if err != nil {
		return nil, err
	}

	obj.ID = id

	createdAt, err := parseLastModified(resp.Header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	obj.CreatedAt = createdAt

	obj.ContentType = resp.Header.Get("Content-Type")

	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return nil, err
	}
	obj.ContentLength = int32(contentLength)

	obj.Data = body

	obj.MD5Checksum = resp.Header.Get("ETag")

	obj.SHA512Checksum = resp.Header.Get("x-sha512-checksum")

	metadata, err := getMetadata(resp.Header)
	if err != nil {
		return nil, err
	}
	obj.Metadata = metadata

	if len(resp.Header.Get("x-revision")) != 0 {
		revision, err := strconv.Atoi(resp.Header.Get("x-revision"))
		if err != nil {
			return nil, err
		}
		obj.Revision = revision
	} else {
		obj.Revision = -1
	}

	return &obj, err
}

type PutDocumentResponse struct {
	Key     string
	Version string
	Bucket  string
}

type PutDocumentOptions struct {
}

// Put a document
func (r RESTDocumentClient) Put(ctx context.Context, document *Document, options PutDocumentOptions) (*PutDocumentResponse, error) {

	document.ID = normaliseKey(document.ID)

	_, res, err := r.c.Put().
		Path(path.Join(server.DocumentsEndpoint, document.ID)).
		WithHeader("Content-Type", document.ContentType).
		WithParams(url.Values{
			"bucketId": []string{document.BucketId},
		}).
		Body(document.Data).
		Do(ctx).
		Raw()
	if err != nil {
		return nil, err
	}

	return &PutDocumentResponse{
		Bucket:  res.Header.Get("x-object-bucket"),
		Key:     res.Header.Get("x-object-key"),
		Version: res.Header.Get("x-object-version"),
	}, nil

}

type DeleteDocumentOptions struct {
	BucketID string
	Version  string
}

// Delete a document
func (r RESTDocumentClient) Delete(ctx context.Context, key string, options DeleteDocumentOptions) error {

	key = normaliseKey(key)

	_, _, err := r.c.
		Delete().
		Path(path.Join(server.DocumentsEndpoint, key)).
		WithParams(url.Values{
			"bucketId": []string{options.BucketID},
			"version":  []string{options.Version},
		}).
		Do(ctx).
		Raw()

	if err != nil {
		return err
	}

	return nil

}
