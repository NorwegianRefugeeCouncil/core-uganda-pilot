package documents

import (
	"context"
	"github.com/nrc-no/core/pkg/api/meta"
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

func (s *Suite) TestGetObjectVersionWithOtherVersionsDeleted() {

	bucket, err := s.createVersionedBucket("TestGetObjectVersionWithOtherVersionsDeleted")
	if !assert.NoError(s.T(), err) {
		return
	}

	bucketId := bucket.ID
	key := "/testObject"

	firstVersion, err := s.putDocument(bucketId, key)
	if !assert.NoError(s.T(), err) {
		return
	}

	_, err = s.putDocument(bucketId, key)
	if !assert.NoError(s.T(), err) {
		return
	}

	err = s.deleteDocumentVersion(bucketId, key, firstVersion.Version)
	if !assert.NoError(s.T(), err) {
		return
	}

	_, err = s.getDocumentVersion(bucketId, key, firstVersion.Version)
	assert.Error(s.T(), err)

	_, err = s.getDocument(bucketId, key)
	assert.NoError(s.T(), err)
}

func (s *Suite) TestGetObjectWithPreviousVersionsDeleted() {
	bucket, err := s.createVersionedBucket("TestGetObjectWithPreviousVersionsDeleted")
	if !assert.NoError(s.T(), err) {
		return
	}

	bucketId := bucket.ID
	key := "/testObject"

	firstVersion, err := s.putDocument(bucketId, key)
	if !assert.NoError(s.T(), err) {
		return
	}

	secondVersion, err := s.putDocument(bucketId, key)
	if !assert.NoError(s.T(), err) {
		return
	}

	err = s.deleteDocumentVersion(bucketId, key, firstVersion.Version)
	if !assert.NoError(s.T(), err) {
		return
	}

	_, err = s.getDocumentVersion(bucketId, key, firstVersion.Version)
	assert.Error(s.T(), err)

	_, err = s.getDocumentVersion(bucketId, key, secondVersion.Version)
	assert.NoError(s.T(), err)
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
		name               string
		path               string
		expectError        bool
		expectErrorReason  meta.StatusReason
		expectEtag         string
		expectLastModified string
		expectBody         []byte
	}

	tcs := []args{
		{
			name:               "getExistingObject",
			path:               "/testobj",
			expectError:        false,
			expectEtag:         getMD5Checksum(existingObj.Data),
			expectLastModified: formatHTTPLastModified(s.timeTeller.TellTime()),
			expectBody:         existingObj.Data,
		}, {
			name:              "getNonExistingObject",
			path:              "/nonExisting",
			expectError:       true,
			expectErrorReason: meta.StatusReasonNotFound,
		}, {
			name:              "getDeletedObject",
			path:              deletedObj.ID,
			expectError:       true,
			expectErrorReason: meta.StatusReasonNotFound,
		}, {
			name:              "getUpdatedThenDeletedObject",
			path:              updatedDeletedObj.ID,
			expectError:       true,
			expectErrorReason: meta.StatusReasonNotFound,
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
				reason := meta.ReasonForError(err)
				assert.Equal(t, tc.expectErrorReason, reason)
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
