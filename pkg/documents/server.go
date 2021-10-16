package documents

import (
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"net"
	"net/http"
	"strings"
)

const (
	collDocuments = "documents"
	collBuckets   = "document_buckets"

	headerContentType    = "Content-Type"
	headerContentLength  = "Content-Length"
	headerETag           = "ETag"
	headerLastModified   = "Last-Modified"
	headerBucketID       = "x-bucket-id"
	headerObjectVersion  = "x-object-version"
	headerObjectKey      = "x-object-key"
	headerSha512Checksum = "x-sha512-checksum"
	headerTags           = "x-tags"

	paramVersion  = "version"
	paramBucketID = "bucketId"

	keyID             = "id"
	keyBucketID       = "bucketId"
	keyIsLastRevision = "isLastRevision"
	keyRevision       = "revision"
	keyIsDeleted      = "isDeleted"
	keyDeletedAt      = "deletedAt"

	mimeTypeApplicationJson = "application/json"
	mimeTypeTextPlain       = "text/plain"
	mimeTypeTextHtml        = "text/html"
)

type Server struct {
	mongoFn        func() (*mongo.Client, error)
	databaseName   string
	timeTeller     utils.TimeTeller
	uidGenerator   utils.UIDGenerator
	listener       net.Listener
	getDocument    http.HandlerFunc
	putDocument    http.HandlerFunc
	deleteDocument http.HandlerFunc
	createBucket   http.HandlerFunc
	deleteBucket   http.HandlerFunc
	getBucket      http.HandlerFunc
}

func (s *Server) NewClient() Interface {
	return NewFromConfig(&rest.Config{
		Scheme:     "http",
		Host:       s.listener.Addr().String(),
		HTTPClient: http.DefaultClient,
	})
}

func NewServer(
	mongoFn func() (*mongo.Client, error),
	databaseName string,
	timeTeller utils.TimeTeller,
	uidGenerator utils.UIDGenerator,
) *Server {
	s := &Server{
		mongoFn:        mongoFn,
		databaseName:   databaseName,
		timeTeller:     timeTeller,
		uidGenerator:   uidGenerator,
		getDocument:    Get(databaseName, mongoFn),
		putDocument:    Put(timeTeller, mongoFn, databaseName),
		deleteDocument: Delete(mongoFn, databaseName, timeTeller),
		getBucket:      GetBucket(mongoFn, databaseName),
		createBucket:   CreateBucket(mongoFn, databaseName, uidGenerator),
		deleteBucket:   DeleteBucket(mongoFn, databaseName),
	}
	return s
}

func (s *Server) Start(done chan struct{}) error {

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return err
	}
	s.listener = listener

	go func() {
		<-done
		listener.Close()
	}()

	go func() {
		if err := http.Serve(listener, s); err != nil {

		}
	}()

	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := s.getHandler(req)
	handler.ServeHTTP(w, req)
}

func (s *Server) getHandler(req *http.Request) http.Handler {
	if strings.HasPrefix(req.URL.Path, server.DocumentsEndpoint) {
		switch req.Method {
		case http.MethodGet:
			return s.getDocument
		case http.MethodPut:
			return s.putDocument
		case http.MethodDelete:
			return s.deleteDocument
		}
	} else if strings.HasPrefix(req.URL.Path, server.BucketsEndpoint) {
		switch req.Method {
		case http.MethodGet:
			return s.getBucket
		case http.MethodPost:
			return s.createBucket
		case http.MethodDelete:
			return s.deleteBucket
		}
	}
	return http.NotFoundHandler()
}

func (s *Server) GetAddress() string {
	return s.listener.Addr().String()
}
