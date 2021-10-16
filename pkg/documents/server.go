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
	if strings.HasPrefix(req.URL.Path, server.DocumentsEndpoint) {
		if req.Method == "GET" {
			s.getDocument(w, req)
		} else if req.Method == "PUT" {
			s.putDocument(w, req)
		} else if req.Method == "DELETE" {
			s.deleteDocument(w, req)
		}
	} else if strings.HasPrefix(req.URL.Path, server.BucketsEndpoint) {
		if req.Method == "GET" {
			s.getBucket(w, req)
		} else if req.Method == "POST" {
			s.createBucket(w, req)
		} else if req.Method == "DELETE" {
			s.deleteBucket(w, req)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) GetAddress() string {
	return s.listener.Addr().String()
}
