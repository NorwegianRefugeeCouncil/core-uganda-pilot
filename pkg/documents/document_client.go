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

	// BucketID of the document (required)
	BucketID string

	// Version of the document (optional)
	Version string
}

// Get a document
func (r RESTDocumentClient) Get(ctx context.Context, id string, options GetDocumentOptions) (*Document, error) {

	id = removeLeadingTrailingSlashes(id)

	var obj Document
	body, resp, err := r.c.Get().Path(path.Join(server.DocumentsEndpoint, id)).
		WithParams(url.Values{
			paramBucketID: []string{options.BucketID},
			paramVersion:  []string{options.Version},
		}).
		Do(ctx).Raw()
	if err != nil {
		return nil, err
	}

	obj.ID = id

	createdAt, err := parseHTTPLastModified(resp.Header.Get(headerLastModified))
	if err != nil {
		return nil, err
	}
	obj.CreatedAt = createdAt

	obj.ContentType = resp.Header.Get(headerContentType)

	contentLength, err := strconv.Atoi(resp.Header.Get(headerContentLength))
	if err != nil {
		return nil, err
	}
	obj.ContentLength = int32(contentLength)

	obj.Data = body

	obj.MD5Checksum = resp.Header.Get(headerETag)

	obj.SHA512Checksum = resp.Header.Get(headerSha512Checksum)

	metadata, err := getDocumentMetadataFromHTTPHeader(resp.Header)
	if err != nil {
		return nil, err
	}
	obj.Metadata = metadata

	if len(resp.Header.Get(headerObjectVersion)) != 0 {
		revision, err := strconv.Atoi(resp.Header.Get(headerObjectVersion))
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

	// Key of the document
	Key string

	// Version of the document
	Version string

	// BucketID of the document
	BucketID string
}

type PutDocumentOptions struct {
}

// Put a document
func (r RESTDocumentClient) Put(ctx context.Context, document *Document, options PutDocumentOptions) (*PutDocumentResponse, error) {

	document.ID = removeLeadingTrailingSlashes(document.ID)

	_, res, err := r.c.Put().
		Path(path.Join(server.DocumentsEndpoint, document.ID)).
		WithHeader(headerContentType, document.ContentType).
		WithParams(url.Values{
			paramBucketID: []string{document.BucketId},
		}).
		Body(document.Data).
		Do(ctx).
		Raw()
	if err != nil {
		return nil, err
	}

	return &PutDocumentResponse{
		BucketID: res.Header.Get(headerBucketID),
		Key:      res.Header.Get(headerObjectKey),
		Version:  res.Header.Get(headerObjectVersion),
	}, nil

}

type DeleteDocumentOptions struct {

	// BucketID of the document (required)
	BucketID string

	// Version of the document (optional)
	Version string
}

// Delete a document
func (r RESTDocumentClient) Delete(ctx context.Context, key string, options DeleteDocumentOptions) error {

	key = removeLeadingTrailingSlashes(key)

	_, _, err := r.c.
		Delete().
		Path(path.Join(server.DocumentsEndpoint, key)).
		WithParams(url.Values{
			paramBucketID: []string{options.BucketID},
			paramVersion:  []string{options.Version},
		}).
		Do(ctx).
		Raw()

	if err != nil {
		return err
	}

	return nil

}
