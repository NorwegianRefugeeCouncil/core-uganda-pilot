package testing

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/client"
	"github.com/nrc-no/core/apps/api/pkg/client/nrc"
	"github.com/nrc-no/core/apps/api/pkg/client/rest"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/formdefinitions"
	"github.com/nrc-no/core/apps/api/pkg/server"
	"github.com/nrc-no/core/apps/api/pkg/storage/eventstoredb"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http/httptest"
	"testing"
)

type MainTestSuite struct {
	suite.Suite
	ctx                context.Context
	httpServer         *httptest.Server
	apiServer          *server.Server
	nrcClient          *nrc.NrcCoreClient
	mongoClient        *mongo.Client
	eventStoreDBClient *client.Client
	store              *eventstoredb.Store
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

	// Create eventdb client
	eventStoreDBClient, err := client.NewClient(&client.Configuration{
		Address:    "localhost:2113",
		DisableTLS: true,
	})
	if err != nil {
		s.T().Errorf("failed to create eventstoredb client: %v", err)
		return
	}
	if err := eventStoreDBClient.Connect(); err != nil {
		s.T().Errorf("failed to connect to evenstore: %v", err)
		return
	}
	s.eventStoreDBClient = eventStoreDBClient

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:30001"))
	if err != nil {
		s.T().Errorf("could not connect to mongo: %v", err)
		return
	}
	s.mongoClient = mongoClient

	// Create storage
	store := eventstoredb.NewStore(eventStoreDBClient, mongoClient)
	s.store = store

	// Install FormDefinitions api
	formdefinitions.Install(apiServer.Container, store)

}

func (s *MainTestSuite) TearDownSuite() {
	defer s.httpServer.Close()
	defer s.mongoClient.Disconnect(s.ctx)
	defer s.eventStoreDBClient.Close()
}
