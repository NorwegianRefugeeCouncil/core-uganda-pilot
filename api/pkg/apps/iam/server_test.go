package iam

import (
	"context"
	"errors"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"testing"
)

type Suite struct {
	suite.Suite
	server     *Server
	serverOpts *server.GenericServerOptions
	ctx        context.Context
	client     *ClientSet
}

func (s *Suite) SetupSuite() {

	ctx := context.Background()

	mongoClient, err := mongo.NewClient(options.Client().SetAuth(options.Credential{Username: "root", Password: "example"}).SetHosts([]string{"localhost:27017"}))
	if !assert.NoError(s.T(), err) {
		s.T().Fatal(err)
		return
	}
	if err := mongoClient.Connect(ctx); !assert.NoError(s.T(), err) {
		s.T().Fatal(err)
		return
	}

	opts := &server.GenericServerOptions{
		MongoClient:   mongoClient,
		MongoDatabase: "e2e",
		Environment:   "Development",
	}
	s.serverOpts = opts

	srv, err := NewServer(ctx, opts)
	if !assert.NoError(s.T(), err) {
		s.T().Fatal()
	}

	s.ctx = ctx
	s.server = srv
	s.client = NewClientSet(&rest.RESTConfig{
		Scheme: "http",
		Host:   "localhost:9001",
		Headers: map[string][]string{
			"X-Authenticated-User-Subject": {"mock-auth-user"},
		},
	})

	s.ResetDB()

	go func() {
		if err := http.ListenAndServe(":9001", srv); err != nil {
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
