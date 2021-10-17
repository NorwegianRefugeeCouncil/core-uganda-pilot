package documents

import (
	"context"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/storage"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	server       *Server
	mongoFn      func() (*mongo.Client, error)
	mongoCli     *mongo.Client
	databaseName string
	timeTeller   utils.TimeTeller
	uidGenerator utils.UIDGenerator
	client       Interface
	done         chan struct{}
}

func (s *Suite) SetupSuite() {

	s.databaseName = "test"

	var err error

	s.mongoCli, err = mongo.Connect(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}

	s.mongoFn = func() (*mongo.Client, error) {
		return s.mongoCli, nil
	}

	s.timeTeller = utils.NewMockTimeTeller(time.Now())
	s.uidGenerator = utils.NewUIDGenerator()

	dbFactory := storage.NewFactory(s.mongoFn)

	s.server = NewServer(dbFactory, s.databaseName, s.timeTeller, s.uidGenerator)

	s.done = make(chan struct{}, 1)

	if err := s.server.Start(s.done); err != nil {
		s.T().Fatal(err)
	}

	s.client = s.server.NewClient()

	if err := ClearCollections(context.Background(), s.mongoCli, s.databaseName); err != nil {
		s.T().Fatal(err)
	}

}

func (s *Suite) TearDownSuite() {
	s.done <- struct{}{}
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}

func (s *Suite) createVersionedBucket(name string) (*Bucket, error) {
	return s.client.Buckets().Create(context.Background(), &Bucket{Name: name, Versioning: VersioningEnabled}, CreateBucketOptions{})
}

func (s *Suite) createVersionedBucketOrDie(name string) *Bucket {
	bucket, err := s.createVersionedBucket(name)
	if !assert.NoError(s.T(), err) {
		s.T().Fatal(err)
	}
	return bucket
}

func (s *Suite) createUnversionedBucket(name string) (*Bucket, error) {
	return s.client.Buckets().Create(context.Background(), &Bucket{Name: name, Versioning: VersioningDisabled}, CreateBucketOptions{})
}

func (s *Suite) createUnversionedBucketOrDie(name string) *Bucket {
	bucket, err := s.client.Buckets().Create(context.Background(), &Bucket{Name: name, Versioning: VersioningDisabled}, CreateBucketOptions{})
	if !assert.NoError(s.T(), err) {
		s.T().Fatal(err)
	}
	return bucket
}

func (s *Suite) getDocumentVersion(bucketId, key, version string) (*Document, error) {
	return s.client.Documents().Get(context.Background(), key, GetDocumentOptions{
		BucketID: bucketId,
		Version:  version,
	})
}
func (s *Suite) getDocumentVersionOrDie(bucketId, key, version string) *Document {
	doc, err := s.getDocumentVersion(bucketId, key, version)
	if !assert.NoError(s.T(), err) {
		s.T().Fatal(err)
	}
	return doc
}

func (s *Suite) getDocument(bucketId, key string) (*Document, error) {
	return s.client.Documents().Get(context.Background(), key, GetDocumentOptions{
		BucketID: bucketId,
	})
}

func (s *Suite) deleteDocument(bucketId, key string) error {
	return s.client.Documents().Delete(context.Background(), key, DeleteDocumentOptions{
		BucketID: bucketId,
	})
}

func (s *Suite) deleteDocumentOrDie(bucketId, key string) {
	err := s.deleteDocument(bucketId, key)
	if !assert.NoError(s.T(), err) {
		s.T().Fatal(err)
	}
	return
}

func (s *Suite) deleteDocumentVersion(bucketId, key, version string) error {
	return s.client.Documents().Delete(context.Background(), key, DeleteDocumentOptions{
		BucketID: bucketId,
		Version:  version,
	})
}

func (s *Suite) deleteDocumentVersionOrDie(bucketId, key, version string) {
	err := s.deleteDocumentVersion(bucketId, key, version)
	if !assert.NoError(s.T(), err) {
		s.T().Fatal(err)
	}
	return
}

func (s *Suite) putDocument(bucketId, key string) (*PutDocumentResponse, error) {
	return s.client.Documents().Put(context.Background(), &Document{
		ID:          key,
		BucketId:    bucketId,
		ContentType: "application/json",
		Data:        []byte(`{"a":"b"}`),
	}, PutDocumentOptions{})
}

func (s *Suite) putDocumentOrDie(bucketId, key string) *PutDocumentResponse {
	r, err := s.putDocument(bucketId, key)
	if !assert.NoError(s.T(), err) {
		s.T().Fatal(err)
	}
	return r
}

func (s *Suite) assertDocumentVersionNotFound(bucketId, key, version string) bool {
	_, err := s.getDocumentVersion(bucketId, key, version)
	return assert.Equal(s.T(), meta.StatusReasonNotFound, meta.ReasonForError(err))
}

func (s *Suite) assertDocumentNotFound(bucketId, key string) bool {
	_, err := s.getDocument(bucketId, key)
	return assert.Equal(s.T(), meta.StatusReasonNotFound, meta.ReasonForError(err))
}
