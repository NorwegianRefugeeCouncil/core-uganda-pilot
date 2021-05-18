package e2e

import (
	"context"
	"github.com/nrc-no/core/api/cmd/server/app"
	v1 "github.com/nrc-no/core/api/pkg/generated/clientset/versioned/typed/core/v1"
	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/flowcontrol"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	cancel context.CancelFunc
	client *v1.CoreV1Client
	stopCh chan struct{}
}

func (s *Suite) SetupSuite() {
	t := s.T()

	options := app.NewCoreServerOptions()
	options.RecommendedOptions.Etcd.StorageConfig.Transport.ServerList = []string{
		"localhost:2379",
	}
	if err := options.Complete(); err != nil {
		t.Fatal(err)
	}
	if err := options.Validate([]string{}); err != nil {
		t.Fatal(err)
	}
	stopCh := make(chan struct{})
	s.stopCh = stopCh

	go func() {
		if err := options.RunCoreServer(stopCh); err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(2 * time.Second)

	restConfig := &rest.Config{
		Host:        "http://localhost:8001",
		RateLimiter: flowcontrol.NewFakeAlwaysRateLimiter(),
	}
	client, err := v1.NewForConfig(restConfig)
	if err != nil {
		t.Fatalf("could not create rest client: %v", err)
	}
	s.client = client
}

func (s *Suite) TearDownSuite() {
	s.stopCh <- struct{}{}
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
