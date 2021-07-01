package iam

import (
	"context"
	"errors"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"testing"
)

type Suite struct {
	suite.Suite
	server                  *Server
	serverOpts              *server.GenericServerOptions
	ctx                     context.Context
	client                  *ClientSet
	authenticatedHttpClient *http.Client
}

func (s *Suite) SetupSuite() {

	opts := server.NewServerOptions().
		WithMongoHosts([]string{"localhost:27017"}).
		WithMongoDatabase("iam_test").
		WithMongoUsername("root").
		WithMongoPassword("example").
		WithListenAddress(":9001")

	s.serverOpts = opts

	ctx := context.Background()
	srv, err := NewServer(ctx, opts)
	if !assert.NoError(s.T(), err) {
		s.T().Fatal()
	}

	clientCredsConfig := clientcredentials.Config{
		ClientID:     "webapp",
		ClientSecret: "somesupersecret",
		TokenURL:     "http://localhost:4444/oauth2/token",
	}
	httpClient := clientCredsConfig.Client(ctx)

	s.ctx = ctx
	s.server = srv
	s.client = NewClientSet(&rest.RESTConfig{
		Scheme:     "http",
		Host:       "localhost:9001",
		HTTPClient: httpClient,
	})

	s.ResetDB()

	s.authenticatedHttpClient = httpClient
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

func (s *Suite) ResetDB() {
	if err := s.server.mongoClient.Database(s.serverOpts.MongoDatabase).Drop(s.ctx); !assert.NoError(s.T(), err) {
		s.T().Fatal()
	}
	if err := s.server.Init(s.ctx); !assert.NoError(s.T(), err) {
		s.T().Fatal()
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
