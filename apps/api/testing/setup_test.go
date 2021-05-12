package testing

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg/api/defaultscheme"
	"github.com/nrc-no/core/apps/api/pkg/client/informers"
	"github.com/nrc-no/core/apps/api/pkg/client/nrc"
	"github.com/nrc-no/core/apps/api/pkg/client/rest"
	"github.com/nrc-no/core/apps/api/pkg/controlplane"
	"github.com/nrc-no/core/apps/api/pkg/server"
	"github.com/nrc-no/core/apps/api/pkg/server/options"
	serverstorage "github.com/nrc-no/core/apps/api/pkg/server/storage"
	storagebackend "github.com/nrc-no/core/apps/api/pkg/storage/backend"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MainTestSuite struct {
	suite.Suite
	ctx            context.Context
	httpServer     *httptest.Server
	apiServer      *controlplane.Instance
	nrcClient      nrc.Interface
	destroyStorage func()
	restConfig     *rest.Config
	informer       informers.SharedInformerFactory
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

func (s *MainTestSuite) SetupSuite() {

	logrus.SetLevel(logrus.TraceLevel)

	ctx := context.Background()
	s.ctx = ctx

	storageConfig := &storagebackend.Config{
		Transport: storagebackend.TransportConfig{
			ServerList: []string{
				"localhost:30001",
			},
		},
		Prefix: "test",
	}

	config := &controlplane.Config{
		GenericConfig: server.NewConfig(defaultscheme.Codecs),
		ExtraConfig: controlplane.ExtraConfig{
			APIResourceConfigSource: controlplane.DefaultAPIResourceConfigSource(),
		},
	}

	resourceEncoding := serverstorage.NewDefaultResourceEncodingConfig(defaultscheme.Scheme)
	storageFactory := serverstorage.NewDefaultStorageFactory(
		*storageConfig,
		"application/json",
		defaultscheme.Codecs,
		resourceEncoding,
		controlplane.DefaultAPIResourceConfigSource(),
	)
	mongoOptions := options.NewMongoOptions(storageConfig)
	if !assert.NoError(s.T(), mongoOptions.ApplyWithStorageFactoryTo(storageFactory, config.GenericConfig)) {
		return
	}

	config.ExtraConfig.StorageFactory = storageFactory
	config.GenericConfig.LoopbackClientConfig = &rest.Config{
		ContentConfig: rest.ContentConfig{
			NegotiatedSerializer: defaultscheme.Codecs,
		},
		APIPath: "/apis",
	}

	apiServer, err := config.Complete().New(server.NewEmptyDelegate())
	if !assert.NoError(s.T(), err) {
		panic(err)
	}
	s.apiServer = apiServer

	httpServer := httptest.NewServer(&mockHttpServer{
		delegate: func(writer http.ResponseWriter, request *http.Request) {
			apiServer.GenericAPIServer.Handler.ServeHTTP(writer, request)
		},
	})
	s.httpServer = httpServer

	// Create client
	restConfig := &rest.Config{
		ContentConfig: rest.DefaultContentConfig,
		Host:          httpServer.URL,
	}
	s.restConfig = restConfig

	nrcClient, err := nrc.NewForConfig(restConfig)
	if !assert.NoError(s.T(), err) {
		panic(err)
	}
	s.nrcClient = nrcClient

	informer := informers.NewSharedInformerFactory(s.nrcClient, time.Hour*5)
	s.informer = informer

}

type mockHttpServer struct {
	delegate http.HandlerFunc
}

func (m *mockHttpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.delegate(w, req)
}

func (s *MainTestSuite) TearDownSuite() {
	defer s.httpServer.Close()
	if s.destroyStorage != nil {
		defer s.destroyStorage()
	}
}
