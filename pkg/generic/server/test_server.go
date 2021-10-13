package server

import (
	"context"
	"errors"
	"github.com/ory/hydra-client-go/client"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"net/http"
	"os"
	"testing"
)

type GenericServerTestSetup struct {
	*GenericServerOptions
	Ctx      context.Context
	Listener *net.TCPListener
	Port     string
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

var (
	mongoUsername = GetEnvOrDefault("MONGO_USERNAME", "root")
	mongoPassword = GetEnvOrDefault("MONGO_PASSWORD", "example")
	mongoHost     = GetEnvOrDefault("MONGO_HOST", "localhost:27017")
	mongoDatabase = GetEnvOrDefault("MONGO_DATABASE", "e2e")
)

type GenericServerTestSuiteArgs struct {
	Listener      *net.TCPListener
	MongoClientFn func(ctx context.Context) (*mongo.Client, error)
	Port          string
	Options       GenericServerOptions
}

func NewGenericServerTestSetup() *GenericServerTestSetup { // Using a random port
	ip := net.ParseIP("127.0.0.1")
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP: ip,
	})

	if err != nil || listener == nil {
		panic(err)
	}
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		panic(err)
	}

	var mongoClientFn = func(ctx context.Context) (*mongo.Client, error) {
		mongoClient, err := mongo.NewClient(options.Client().SetAuth(options.Credential{Username: mongoUsername, Password: mongoPassword}).SetHosts([]string{mongoHost}))
		if err != nil {
			panic(err)
		}
		if err := mongoClient.Connect(ctx); err != nil {
			logrus.WithError(err).Errorf("failed to connect to mongo")
			return nil, err
		}
		return mongoClient, nil
	}

	hydraClient := client.NewHTTPClient(nil)

	return &GenericServerTestSetup{
		GenericServerOptions: &GenericServerOptions{
			MongoClientFn:    mongoClientFn,
			MongoDatabase:    mongoDatabase,
			Environment:      "Development",
			HydraAdminClient: hydraClient,
		},
		Ctx:      context.Background(),
		Listener: listener,
		Port:     port,
	}
}

func (s *GenericServerTestSetup) Serve(t *testing.T, handler http.Handler) {
	go func() {
		if err := http.Serve(s.Listener, handler); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
		} else {
			t.Fatal(err)
		}
	}()
}
