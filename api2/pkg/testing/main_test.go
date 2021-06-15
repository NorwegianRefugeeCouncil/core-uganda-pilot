package testing

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/server"
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
	s.server = server.NewServer(s.ctx, server.ServerOptions{
		TemplateDirectory: os.Getenv("TEMPLATE_DIRECTORY"),
	})
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
