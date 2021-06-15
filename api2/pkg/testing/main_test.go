package testing

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
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
		TemplateDirectory: os.Getenv("TEMPLATE_DIRECTORY"),
		MongoUsername:     "root",
		MongoPassword:     "example",
		MongoDatabase:     "e2e",
		Address:           "http://localhost:9001",
	}

	completedConfig, err := serverOptions.Complete(s.ctx)
	if !assert.NoError(s.T(), err) {
		return
	}

	if err := ClearDatabase(s.ctx, completedConfig.MongoClient); !assert.NoError(s.T(), err) {
		panic(err)
	}

	s.server = completedConfig.New(s.ctx)
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
