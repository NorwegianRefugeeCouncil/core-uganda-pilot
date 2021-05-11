package testing

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg/api/defaultscheme"
	"github.com/nrc-no/core/apps/api/pkg/client/nrc"
	"github.com/nrc-no/core/apps/api/pkg/client/rest"
	"github.com/nrc-no/core/apps/api/pkg/controlplane"
	"github.com/nrc-no/core/apps/api/pkg/server"
	"github.com/nrc-no/core/apps/api/pkg/server/options"
	serverstorage "github.com/nrc-no/core/apps/api/pkg/server/storage"
	storagebackend "github.com/nrc-no/core/apps/api/pkg/storage/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

type MainTestSuite struct {
	suite.Suite
	ctx            context.Context
	httpServer     *httptest.Server
	apiServer      *controlplane.Instance
	nrcClient      *nrc.NrcCoreClient
	destroyStorage func()
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

func (s *MainTestSuite) SetupSuite() {
	ctx := context.Background()
	s.ctx = ctx

	storageConfig := &storagebackend.Config{
		Transport: storagebackend.TransportConfig{
			ServerList: []string{
				"localhost:30001",
			},
		},
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

	httpServer := httptest.NewServer(apiServer.GenericAPIServer.Handler.GoRestfulContainer.ServeMux)
	s.httpServer = httpServer

	// Create client
	nrcClient, err := nrc.NewForConfig(&rest.Config{
		ContentConfig: rest.DefaultContentConfig,
		Host:          httpServer.URL,
	})
	if !assert.NoError(s.T(), err) {
		panic(err)
	}
	s.nrcClient = nrcClient

}

func (s *MainTestSuite) TearDownSuite() {
	defer s.httpServer.Close()
	if s.destroyStorage != nil {
		defer s.destroyStorage()
	}
}
