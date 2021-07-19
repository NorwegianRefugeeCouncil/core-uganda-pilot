package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"os"
)

type GenericServerTestSuite struct {
	serverOpts *GenericServerOptions
	ctx        context.Context
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

func (s *GenericServerTestSuite) GenericSetupSuite() GenericServerTestSuiteArgs {
	// Using a random port
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

	opts := GenericServerOptions{
		MongoClientFn: mongoClientFn,
		MongoDatabase: mongoDatabase,
		Environment:   "Development",
	}
	return GenericServerTestSuiteArgs{
		Listener:      listener,
		MongoClientFn: mongoClientFn,
		Port:          port,
		Options:       opts,
	}
}
