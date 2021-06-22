package testing

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite
	ctx    context.Context
	cancel context.CancelFunc
	server *server.Server
}

func (s *Suite) SetupSuite() {

	s.ctx, s.cancel = context.WithCancel(context.Background())
	serverOptions := server.Options{
		TemplateDirectory:    "../../pkg/webapp/templates",
		Address:              "http://localhost:9001",
		MongoDatabase:        "e2e",
		MongoUsername:        "root",
		MongoPassword:        "example",
		KeycloakClientID:     "",
		KeycloakClientSecret: "e6486272-039d-430f-b3c7-47887aa9e206",
		KeycloakBaseURL:      "http://localhost:8080",
		KeycloakRealmName:    "nrc",
		RedisMaxIdleConns:    10,
		RedisNetwork:         "tcp",
		RedisAddress:         "localhost:6379",
		RedisPassword:        "",
		RedisSecretKey:       "some-secret",
	}

	completedConfig, err := serverOptions.Complete(s.ctx)
	if !assert.NoError(s.T(), err) {
		return
	}

	if err := ClearDatabase(s.ctx, completedConfig.MongoDatabase, completedConfig.MongoClient); !assert.NoError(s.T(), err) {
		panic(err)
	}

	s.server = completedConfig.New(s.ctx)
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
