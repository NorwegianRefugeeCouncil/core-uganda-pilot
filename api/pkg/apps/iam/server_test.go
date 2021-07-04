// +build integration

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
	"net"
	"net/http"
	"os"
	"testing"
)

type Suite struct {
	suite.Suite
	server     *Server
	serverOpts *server.GenericServerOptions
	ctx        context.Context
	client     *ClientSet
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func (s *Suite) SetupSuite() {

	ctx := context.Background()

	mongoUsername := GetEnvOrDefault("MONGO_USERNAME", "root")
	mongoPassword := GetEnvOrDefault("MONGO_PASSWORD", "example")
	mongoHost := GetEnvOrDefault("MONGO_HOST", "localhost:27017")
	mongoDatabase := GetEnvOrDefault("MONGO_DATABASE", "e2e")

	// Using a random port
	ip := net.ParseIP("127.0.0.1")
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP: ip,
	})
	if !assert.NoError(s.T(), err) {
		s.T().Fatal(err)
		return
	}
	s.T().Logf("Listening at: %s", listener.Addr().String())
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if !assert.NoError(s.T(), err) {
		s.T().Fatal(err)
		return
	}

	mongoClient, err := mongo.NewClient(options.Client().SetAuth(options.Credential{Username: mongoUsername, Password: mongoPassword}).SetHosts([]string{mongoHost}))
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
		MongoDatabase: mongoDatabase,
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
		Host:   "localhost:" + port,
		Headers: map[string][]string{
			"X-Authenticated-User-Subject": {"mock-auth-user"},
		},
	})

	s.ResetDB()

	go func() {
		if err := http.Serve(listener, srv); err != nil {
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
