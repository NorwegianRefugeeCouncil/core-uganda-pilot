package e2e

import (
	"context"
	"github.com/nrc-no/core/api/cmd/server/app"
	v1 "github.com/nrc-no/core/api/pkg/generated/clientset/versioned/typed/core/v1"
	"github.com/stretchr/testify/suite"
	apiextensionsclientv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/flowcontrol"
	"net/url"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	cancel     context.CancelFunc
	client     *v1.CoreV1Client
	crdClient  *apiextensionsclientv1.ApiextensionsV1Client
	stopCh     chan struct{}
	stoppedCh  chan struct{}
	restConfig *rest.Config
	baseUrl    *url.URL
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
	s.stopCh = make(chan struct{})
	s.stoppedCh = make(chan struct{})

	go func() {
		if err := options.RunCoreServer(s.stopCh); err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(2 * time.Second)

	baseUrl, _ := url.Parse("http://localhost:8001")
	s.baseUrl = baseUrl

	s.restConfig = &rest.Config{
		Host:        baseUrl.String(),
		RateLimiter: flowcontrol.NewFakeAlwaysRateLimiter(),
	}
	var err error
	s.client, err = v1.NewForConfig(s.restConfig)
	if err != nil {
		t.Fatalf("could not create rest client: %v", err)
	}
	s.crdClient, err = apiextensionsclientv1.NewForConfig(s.restConfig)
	if err != nil {
		t.Fatalf("could not create rest client: %v", err)
	}
}

func (s *Suite) TearDownSuite() {
	s.stopCh <- struct{}{}
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
