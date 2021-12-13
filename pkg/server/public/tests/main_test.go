package tests

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/client"
	"github.com/nrc-no/core/pkg/rest"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/server/public"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type Suite struct {
	suite.Suite
	cli    client.Client
	srv    *public.Server
	cancel context.CancelFunc
	ctx    context.Context
	done   func()
}

func TestMain(m *testing.M) {
	// embedded-postgres runs as a separate process altogether
	// We must setup/teardown here, otherwise the process does
	// not get properly cleaned up
	done, err := testutils.TryGetPostgres()
	defer done()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	exitVal := m.Run()
	done()
	os.Exit(exitVal)
}

func (s *Suite) SetupSuite() {

	s.ctx, s.cancel = context.WithCancel(context.Background())

	dbFactory, err := store.NewFactory("host=localhost port=15432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	db, err := dbFactory.Get()
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	if err := store.Migrate(db); !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	s.srv, err = public.NewServer(s.ctx, public.Options{
		ServerOptions: options.ServerOptions{
			Host: "localhost",
			Port: 0,
			Cors: options.CorsOptions{},
		},
		StoreFactory: dbFactory,
	})

	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	s.srv.Start(s.ctx)
	s.cli = client.NewClientFromConfig(rest.Config{
		Scheme: "http",
		Host:   s.srv.Server.Address(),
	})
}

func (s *Suite) TearDownSuite() {
	if s.cancel != nil {
		s.cancel()
	}
	if s.done != nil {
		s.done()
	}
}

func (s *Suite) TestCanStartServer() {
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
