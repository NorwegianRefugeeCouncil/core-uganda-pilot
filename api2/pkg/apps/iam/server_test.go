package iam

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type Suite struct {
	suite.Suite
	server *Server
	ctx    context.Context
	client *ClientSet
}

func (s *Suite) SetupSuite() {

	opts := NewServerOptions().
		WithMongoHosts([]string{"localhost:27017"}).
		WithMongoDatabase("iam_test").
		WithMongoUsername("root").
		WithMongoPassword("example").
		WithListenAddress(":9001")

	ctx := context.Background()
	srv, err := NewServer(ctx, opts)
	if !assert.NoError(s.T(), err) {
		s.T().Fatal()
	}

	if err := srv.mongoClient.Database(opts.MongoDatabase).Drop(ctx); !assert.NoError(s.T(), err) {
		s.T().Fatal()
	}

	if err := srv.Init(ctx); !assert.NoError(s.T(), err) {
		s.T().Fatal()
	}

	s.ctx = ctx
	s.server = srv
	s.client = NewClientSet(&RESTConfig{
		Scheme: "http",
		Host:   "localhost:9001",
	})

	go func() {
		if err := http.ListenAndServe(opts.ListenAddress, srv); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
		} else {
			s.T().Fatal(err)
		}
	}()

}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
