package testing

import (
	"context"
	corev1 "github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/client/nrc"
	"github.com/nrc-no/core/apps/api/pkg/client/rest"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/formdefinitions"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/serializer/json"
	"github.com/nrc-no/core/apps/api/pkg/server"
	"github.com/nrc-no/core/apps/api/pkg/storage/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

type MainTestSuite struct {
	suite.Suite
	ctx            context.Context
	httpServer     *httptest.Server
	apiServer      *server.Server
	nrcClient      *nrc.NrcCoreClient
	destroyStorage func()
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

func (s *MainTestSuite) SetupSuite() {
	ctx := context.Background()
	s.ctx = ctx

	// Create API server
	apiServer := server.NewServer()
	s.apiServer = apiServer

	// Create HTTP server
	httpServer := httptest.NewServer(apiServer)
	s.httpServer = httpServer

	// Create client
	nrcClient, err := nrc.NewForConfig(&rest.Config{
		ContentConfig: rest.DefaultContentConfig,
		Host:          httpServer.URL,
	})
	if err != nil {
		s.T().Errorf("unable to create rest client: %v", err)
		return
	}
	s.nrcClient = nrcClient

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		s.T().Errorf("unable to register scheme: %v", err)
		return
	}

	serializer := json.NewSerializer(json.DefaultMetaFactory, scheme, scheme)

	storageBackend, destroyStorage, err := backend.Create(backend.Config{
		Codec:           corev1.Codecs,
		EncodeVersioner: corev1.SchemeGroupVersion,
		Prefix:          "core_nrc_no/formdefinitions",
		Transport: backend.TransportConfig{
			ServerList: []string{"localhost:30001"},
		},
	}, func() runtime.Object { return &corev1.FormDefinition{} })
	if !assert.NoError(s.T(), err) {
		return
	}
	s.destroyStorage = destroyStorage

	// Install FormDefinitions api
	formdefinitions.Install(
		apiServer.Container,
		storageBackend,
		corev1.SchemeGroupVersion.WithKind("FormDefinition"),
		corev1.SchemeGroupVersion.WithResource("formdefinitions"),
		scheme,
		scheme,
		serializer,
		scheme,
	)

}

func (s *MainTestSuite) TearDownSuite() {
	defer s.httpServer.Close()
	if s.destroyStorage != nil {
		defer s.destroyStorage()
	}
}
