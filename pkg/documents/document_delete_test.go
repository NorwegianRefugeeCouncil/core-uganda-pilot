package documents

import (
	"context"
	"github.com/nrc-no/core/pkg/validation"
	"github.com/stretchr/testify/assert"
	"net/http"
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

func (s *Suite) TestDeleteDocument() {

	ctx := context.Background()

	bucket, err := s.client.Buckets().Create(ctx, &Bucket{Name: "test-document-delete"}, CreateBucketOptions{})
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
		name            string
		key             string
		expectError     bool
		expectErrorCode int
	}

	tcs := []args{
		{
			name:        "deleteExistingObject",
			key:         "/testobj",
			expectError: false,
		}, {
			name:            "deleteNonExistingObject",
			key:             "/nonExisting",
			expectError:     true,
			expectErrorCode: http.StatusNotFound,
		}, {
			name:            "deleteAlreadyDeleted",
			key:             "/deleted",
			expectError:     true,
			expectErrorCode: http.StatusNotFound,
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
				status := validation.AsStatus(err)
				assert.Equal(t, tc.expectErrorCode, status.Code)
			}

		})
	}

}
