package documents

import (
	"context"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/stretchr/testify/assert"
	"testing"
)

func (s *Suite) TestDeleteDocumentWithNoBucketShouldThrow() {
	err := s.client.Documents().Delete(context.Background(), "some-document", DeleteDocumentOptions{BucketID: ""})
	assert.Error(s.T(), err)
}

func (s *Suite) TestDeleteDocumentInvalidBucketShouldThrow() {
	err := s.client.Documents().Delete(context.Background(), "some-document", DeleteDocumentOptions{BucketID: "non-existing"})
	assert.Error(s.T(), err)
}

func (s *Suite) createBucket(name string) (*Bucket, error) {
	return s.client.Buckets().Create(context.Background(), &Bucket{Name: name}, CreateBucketOptions{})
}

func (s *Suite) getDocumentVersion(bucketId, key, version string) (*Document, error) {
	return s.client.Documents().Get(context.Background(), key, GetDocumentOptions{
		BucketID: bucketId,
		Version:  version,
	})
}

func (s *Suite) getDocument(bucketId, key string) (*Document, error) {
	return s.client.Documents().Get(context.Background(), key, GetDocumentOptions{
		BucketID: bucketId,
	})
}

func (s *Suite) deleteDocumentVersion(bucketId, key, version string) error {
	return s.client.Documents().Delete(context.Background(), key, DeleteDocumentOptions{
		BucketID: bucketId,
		Version:  version,
	})
}

func (s *Suite) putDocument(bucketId, key string) (*PutDocumentResponse, error) {
	return s.client.Documents().Put(context.Background(), &Document{
		ID:          key,
		BucketId:    bucketId,
		ContentType: "application/json",
		Data:        []byte(`{"a":"b"}`),
	}, PutDocumentOptions{})
}

func (s *Suite) TestDeleteDocument() {

	ctx := context.Background()

	bucket, err := s.createBucket("test-document-delete")
	if !assert.NoError(s.T(), err) {
		return
	}

	existingObj := &Document{
		ID:          "/testobj",
		BucketId:    bucket.ID,
		ContentType: "application/json",
		Data:        []byte(`{"a":"b"}`),
	}

	deletedObj := &Document{
		ID:          "/deleted",
		BucketId:    bucket.ID,
		ContentType: "application/json",
		Data:        []byte(`{"a":"b"}`),
	}

	updatedObj := &Document{
		ID:          "/delete-updated",
		BucketId:    bucket.ID,
		ContentType: "application/json",
		Data:        []byte(`{"a":"b"}`),
	}

	for _, document := range []*Document{existingObj, deletedObj, updatedObj} {
		if _, err := s.client.Documents().Put(ctx, document, PutDocumentOptions{}); !assert.NoError(s.T(), err) {
			return
		}
	}

	if err := s.client.Documents().Delete(ctx, deletedObj.ID, DeleteDocumentOptions{BucketID: bucket.ID}); !assert.NoError(s.T(), err) {
		return
	}

	if _, err := s.client.Documents().Put(ctx, updatedObj, PutDocumentOptions{}); !assert.NoError(s.T(), err) {
		return
	}

	type args struct {
		name         string
		key          string
		expectError  bool
		expectReason meta.StatusReason
	}

	tcs := []args{
		{
			name:        "deleteExistingObject",
			key:         "/testobj",
			expectError: false,
		}, {
			name:         "deleteNonExistingObject",
			key:          "/nonExisting",
			expectError:  true,
			expectReason: meta.StatusReasonNotFound,
		}, {
			name:         "deleteAlreadyDeleted",
			key:          "/deleted",
			expectError:  true,
			expectReason: meta.StatusReasonNotFound,
		}, {
			name:        "deleteUpdatedDocument",
			key:         "/delete-updated",
			expectError: false,
		},
	}

	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.client.Documents().Delete(context.Background(), tc.key, DeleteDocumentOptions{
				BucketID: bucket.ID,
			})

			if tc.expectError {
				if !assert.Error(t, err) {
					return
				}

				reason := meta.ReasonForError(err)
				assert.Equal(t, tc.expectReason, reason)
			}

		})
	}

}
