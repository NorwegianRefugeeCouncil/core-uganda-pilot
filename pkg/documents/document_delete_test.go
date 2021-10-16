package documents

import (
	"context"
	"github.com/nrc-no/core/pkg/api/meta"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestDeleteDocumentWithNoBucketShouldThrow() {
	err := s.client.Documents().Delete(context.Background(), "some-document", DeleteDocumentOptions{BucketID: ""})
	assert.Error(s.T(), err)
}

func (s *Suite) TestDeleteDocumentInvalidBucketShouldThrow() {
	err := s.client.Documents().Delete(context.Background(), "some-document", DeleteDocumentOptions{BucketID: "non-existing"})
	assert.Error(s.T(), err)
}

func (s *Suite) TestDeleteUnexistingDocShouldReturn404() {
	bucketName := uuid.NewV4().String()
	bucket := s.createVersionedBucketOrDie(bucketName)
	s.assertDocumentNotFound(bucket.ID, "not-existing")
}

func (s *Suite) TestDeleteVersionedDoc() {
	bucketName := uuid.NewV4().String()
	bucket := s.createVersionedBucketOrDie(bucketName)
	v1 := s.putDocumentOrDie(bucket.ID, "somedocument")
	v2 := s.putDocumentOrDie(bucket.ID, "somedocument")
	v3 := s.putDocumentOrDie(bucket.ID, "somedocument")
	s.deleteDocumentVersionOrDie(bucket.ID, "somedocument", v2.Version)
	s.getDocumentVersionOrDie(bucket.ID, "somedocument", v1.Version)
	s.getDocumentVersionOrDie(bucket.ID, "somedocument", v3.Version)
	_, err := s.getDocumentVersion(bucket.ID, "somedocument", v2.Version)
	assert.Equal(s.T(), meta.StatusReasonNotFound, meta.ReasonForError(err))
}

func (s *Suite) TestDeleteDeletedVersionShouldReturn404() {
	bucketName := uuid.NewV4().String()
	bucket := s.createVersionedBucketOrDie(bucketName)
	v1 := s.putDocumentOrDie(bucket.ID, "somedocument")
	v2 := s.putDocumentOrDie(bucket.ID, "somedocument")
	s.deleteDocumentVersionOrDie(bucket.ID, "somedocument", v2.Version)
	s.getDocumentVersionOrDie(bucket.ID, "somedocument", v1.Version)
	s.assertDocumentVersionNotFound(bucket.ID, "somedocument", v2.Version)
	s.deleteDocumentVersionOrDie(bucket.ID, "somedocument", v1.Version)
	s.assertDocumentVersionNotFound(bucket.ID, "somedocument", v1.Version)
}

func (s *Suite) TestDeleteWithoutVersionShouldDeleteLastVersion() {
	bucketName := uuid.NewV4().String()
	bucket := s.createVersionedBucketOrDie(bucketName)
	v1 := s.putDocumentOrDie(bucket.ID, "somedocument")
	v2 := s.putDocumentOrDie(bucket.ID, "somedocument")
	s.deleteDocumentOrDie(bucket.ID, "somedocument")
	s.getDocumentVersionOrDie(bucket.ID, "somedocument", v1.Version)
	s.assertDocumentVersionNotFound(bucket.ID, "somedocument", v2.Version)
}

func (s *Suite) TestDeleteNonVersionedDoc() {
	bucketName := uuid.NewV4().String()
	bucket := s.createUnversionedBucketOrDie(bucketName)
	s.putDocumentOrDie(bucket.ID, "somedocument")
	s.deleteDocumentOrDie(bucket.ID, "somedocument")
	s.assertDocumentNotFound(bucket.ID, "somedocument")
}

func (s *Suite) TestDeleteDocVersionInUnversionedBucketShouldThrow() {
	bucketName := uuid.NewV4().String()
	bucket := s.createUnversionedBucketOrDie(bucketName)
	d := s.putDocumentOrDie(bucket.ID, "somedocument")
	err := s.deleteDocumentVersion(bucket.ID, "somedocument", d.Version)
	assert.Equal(s.T(), meta.StatusReasonBadRequest, meta.ReasonForError(err))
}
