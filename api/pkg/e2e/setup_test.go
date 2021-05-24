package e2e

import (
	"context"
	"github.com/nrc-no/core/api/pkg/client/core"
	"github.com/nrc-no/core/api/pkg/client/rest"
	serveroptions "github.com/nrc-no/core/api/pkg/server/options"
	"github.com/nrc-no/core/api/pkg/store"
	"github.com/stretchr/testify/suite"
	"net"
	"net/url"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	cancel     context.CancelFunc
	client     core.Interface
	restConfig *rest.Config
	baseUrl    *url.URL
}

func (s *Suite) SetupSuite() {
	t := s.T()

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	opts := &serveroptions.Options{
		StorageConfig: serveroptions.MongoOptions{
			StorageConfig: store.Config{
				Transport: store.TransportConfig{
					// Password: "pass12345",
					// Username: "root",
					Database: "test",
					ServerList: []string{
						"mongodb://127.0.0.1:27017",
					},
				},
			},
		},
		BindPort:    8888,
		BindAddress: net.ParseIP("127.0.0.1"),
	}
	if err := opts.Complete(); err != nil {
		panic(err)
	}
	if err := opts.Validate(); err != nil {
		panic(err)
	}

	go func() {
		if err := opts.Run(ctx); err != nil {
			cancel()
		}
	}()

	time.Sleep(5 * time.Second)

	s.restConfig = &rest.Config{
		Host: "http://localhost:8888",
	}
	var err error
	s.client, err = core.NewForConfig(s.restConfig)
	if err != nil {
		t.Fatalf("could not create rest client: %v", err)
	}

}

func (s *Suite) TearDownSuite() {
	s.cancel()
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
