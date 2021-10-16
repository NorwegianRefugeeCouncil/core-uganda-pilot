package documents

import (
	"context"
	"github.com/nrc-no/core/pkg/validation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func (s *Suite) TestGetDocumentWithNoBucketShouldThrow() {
	_, err := s.client.Documents().Get(context.Background(), "some-document", GetDocumentOptions{BucketID: ""})
	assert.Error(s.T(), err)
}

func (s *Suite) TestGetocumentInvalidBucketShouldThrow() {
	_, err := s.client.Documents().Get(context.Background(), "some-document", GetDocumentOptions{BucketID: "non-existing"})
	assert.Error(s.T(), err)
}

func (s *Suite) TestGetDocument() {

	ctx := context.Background()

	bucket, err := s.client.Buckets().Create(ctx, &Bucket{Name: "test-document-get"}, CreateBucketOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}

	existingObj := &Document{
		ID:          "testobj",
		ContentType: "application/json",
		BucketId:    bucket.ID,
		Data:        []byte(`{"a":"b"}`),
	}

	deletedObj := &Document{
		ID:          "/get-deleted",
		ContentType: "application/json",
		BucketId:    bucket.ID,
		Data:        []byte(`{"a":"b"}`),
	}

	updatedObj := &Document{
		ID:          "/get-updated",
		ContentType: "application/json",
		BucketId:    bucket.ID,
		Data:        []byte(`{"a":"b"}`),
	}

	updatedDeletedObj := &Document{
		ID:          "/get-updated-deleted",
		ContentType: "application/json",
		BucketId:    bucket.ID,
		Data:        []byte(`{"a":"b"}`),
	}

	for _, document := range []*Document{existingObj, deletedObj, updatedObj, updatedDeletedObj} {
		if _, err := s.client.Documents().Put(ctx, document, PutDocumentOptions{}); !assert.NoError(s.T(), err) {
			return
		}
	}

	for _, document := range []*Document{updatedObj, updatedDeletedObj} {
		if _, err := s.client.Documents().Put(ctx, document, PutDocumentOptions{}); !assert.NoError(s.T(), err) {
			return
		}
	}

	for _, document := range []*Document{deletedObj, updatedDeletedObj} {
		if err := s.client.Documents().Delete(ctx, document.ID, DeleteDocumentOptions{
			BucketID: bucket.ID,
		}); !assert.NoError(s.T(), err) {
			return
		}
	}

	type args struct {
		name                  string
		path                  string
		expectError           bool
		expectErrorStatusCode int
		expectEtag            string
		expectLastModified    string
		expectBody            []byte
	}

	tcs := []args{
		{
			name:               "getExistingObject",
			path:               "/testobj",
			expectError:        false,
			expectEtag:         getMD5Checksum(existingObj.Data),
			expectLastModified: getLastModified(s.timeTeller.TellTime()),
			expectBody:         existingObj.Data,
		}, {
			name:                  "getNonExistingObject",
			path:                  "/nonExisting",
			expectError:           true,
			expectErrorStatusCode: 404,
		}, {
			name:                  "getDeletedObject",
			path:                  deletedObj.ID,
			expectError:           true,
			expectErrorStatusCode: 404,
		}, {
			name:                  "getUpdatedThenDeletedObject",
			path:                  updatedDeletedObj.ID,
			expectError:           true,
			expectErrorStatusCode: 404,
		},
	}

	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			doc, err := s.client.Documents().Get(context.Background(), tc.path, GetDocumentOptions{
				BucketID: bucket.ID,
			})
			if tc.expectError {
				if !assert.Error(t, err) {
					return
				}
				status := validation.AsStatus(err)
				assert.Equal(t, tc.expectErrorStatusCode, status.Code)

			}
			if !tc.expectError {
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, tc.expectEtag, doc.MD5Checksum)
				assert.Equal(t, tc.expectBody, doc.Data)
			}
		})
	}

}
