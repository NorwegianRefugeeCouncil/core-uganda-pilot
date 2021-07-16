// +build integration

package cms_test

import (
	"context"
	"errors"
	. "github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
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

var (
	ctx = context.Background()

	mongoUsername = GetEnvOrDefault("MONGO_USERNAME", "root")
	mongoPassword = GetEnvOrDefault("MONGO_PASSWORD", "example")
	mongoHost     = GetEnvOrDefault("MONGO_HOST", "localhost:27017")
	mongoDatabase = GetEnvOrDefault("MONGO_DATABASE", "e2e")
)

func (s *Suite) SetupSuite() {
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

	var mongoClientFn = func(ctx context.Context) (*mongo.Client, error) {
		mongoClient, err := mongo.NewClient(options.Client().SetAuth(options.Credential{Username: mongoUsername, Password: mongoPassword}).SetHosts([]string{mongoHost}))
		if !assert.NoError(s.T(), err) {
			s.T().Fatal(err)
		}
		if err := mongoClient.Connect(ctx); err != nil {
			logrus.WithError(err).Errorf("failed to connect to mongo")
			return nil, err
		}
		return mongoClient, nil

	}

	opts := &server.GenericServerOptions{
		MongoClientFn: mongoClientFn,
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

// This will run before each test in the suite but must be called manually before subtests
func (s *Suite) SetupTest() {
	err := s.server.ResetDB(ctx, mongoDatabase)
	if err != nil {
		return
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}

//
// Helpers
//
func contains(s []string, item string) bool {
	for _, a := range s {
		if a == item {
			return true
		}
	}
	return false
}

func newUUID() string {
	return uuid.NewV4().String()
}

func (s *Suite) mockCaseTypes(n int) []*CaseType {
	var caseTypes []*CaseType
	for i := 0; i < n; i++ {
		caseTypes = append(caseTypes, &CaseType{
			ID:   newUUID(),
			Name: "mock",
		})
	}
	return caseTypes
}

func (s *Suite) mockCases(n int) []*Case {
	var cases []*Case
	for i := 0; i < n; i++ {
		cases = append(cases, &Case{
			ID:        newUUID(),
			TeamID:    newUUID(),
			CreatorID: "mock-auth-user",
		})
	}
	return cases
}
