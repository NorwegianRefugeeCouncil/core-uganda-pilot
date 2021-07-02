package server

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/apps/login"
	"github.com/nrc-no/core/pkg/apps/seeder"
	webapp2 "github.com/nrc-no/core/pkg/apps/webapp"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/middleware"
	"github.com/ory/hydra-client-go/client"
	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type Server struct {
	MongoClient       *mongo.Client
	WebAppHandler     *webapp2.Server
	HttpServer        *http.Server
	HydraPublicClient *client.OryHydra
	HydraAdminClient  *client.OryHydra
}

type Options struct {
	ClearDB           bool
	Environment       string
	TemplateDirectory string
	Address           string
	MongoDatabase     string
	MongoUsername     string
	MongoPassword     string
	RedisMaxIdleConns int
	RedisNetwork      string
	RedisAddress      string
	RedisPassword     string
	RedisSecretKey    string
	HydraAdminURL     string
	HydraPublicURL    string
}

func NewOptions() *Options {
	return &Options{
		ClearDB:           false,
		Environment:       "Production",
		TemplateDirectory: "pkg/apps/webapp/templates",
		Address:           "http://localhost:9000",
		MongoDatabase:     "core",
		MongoUsername:     "",
		MongoPassword:     "",
		RedisMaxIdleConns: 10,
		RedisNetwork:      "tcp",
		RedisAddress:      "localhost:6379",
		RedisPassword:     "",
		RedisSecretKey:    "some-secret",
		HydraAdminURL:     "http://localhost:4445",
		HydraPublicURL:    "http://localhost:4444",
	}
}

func (o *Options) Flags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Address, "address", o.Address, "Address")
	fs.StringVar(&o.MongoDatabase, "mongo-database", o.MongoDatabase, "Mongo database name")
	fs.StringVar(&o.MongoUsername, "mongo-username", o.MongoUsername, "Mongo username")
	fs.StringVar(&o.MongoPassword, "mongo-password", o.MongoPassword, "Mongo password")
	fs.IntVar(&o.RedisMaxIdleConns, "redis-max-idle-conns", o.RedisMaxIdleConns, "Redis maximum number of idle connections")
	fs.StringVar(&o.RedisNetwork, "redis-network", o.RedisNetwork, "Redis network")
	fs.StringVar(&o.RedisPassword, "redis-password", o.RedisPassword, "Redis password")
	fs.StringVar(&o.RedisSecretKey, "redis-secret-key", o.RedisSecretKey, "Redis secret key")
	fs.StringVar(&o.HydraAdminURL, "hydra-admin-url", o.HydraAdminURL, "Hydra Admin URL")
	fs.StringVar(&o.HydraPublicURL, "hydra-public-url", o.HydraPublicURL, "Hydra Public URL")
	fs.StringVar(&o.Environment, "environment", o.Environment, "Environment (Production / Development)")
	fs.BoolVar(&o.ClearDB, "fresh", o.ClearDB, "Clear user-created DB entries")
}

type CompletedOptions struct {
	*Options
	MongoClient       *mongo.Client
	HydraAdminClient  *client.OryHydra
	HydraPublicClient *client.OryHydra
}

func (o *Options) Complete(ctx context.Context) (CompletedOptions, error) {

	// Setup mongo client
	mongoClient, err := mongo.NewClient(options.Client().
		SetAuth(
			options.Credential{
				Username: o.MongoUsername,
				Password: o.MongoPassword,
			}))
	if err != nil {
		return CompletedOptions{}, err
	}

	if err := mongoClient.Connect(ctx); err != nil {
		return CompletedOptions{}, err
	}

	hydraAdminURL, err := url.Parse(o.HydraAdminURL)
	if err != nil {
		return CompletedOptions{}, err
	}
	hydraAdminClient := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Schemes: []string{
			hydraAdminURL.Scheme,
		},
		Host:     hydraAdminURL.Host,
		BasePath: hydraAdminURL.Path,
	})

	hydraPublicURL, err := url.Parse(o.HydraPublicURL)
	if err != nil {
		return CompletedOptions{}, err
	}
	hydraPublicCLient := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Schemes: []string{
			hydraPublicURL.Scheme,
		},
		Host:     hydraPublicURL.Host,
		BasePath: hydraPublicURL.Path,
	})

	completedOptions := CompletedOptions{
		Options:           o,
		MongoClient:       mongoClient,
		HydraAdminClient:  hydraAdminClient,
		HydraPublicClient: hydraPublicCLient,
	}
	return completedOptions, nil
}

func (c CompletedOptions) New(ctx context.Context) *Server {

	// Create aggregated router
	router := mux.NewRouter()

	// Add logging middleware
	router.Use(middleware.UseLogging())

	// Prep db
	if c.ClearDB {
		if err := seeder.Clear(ctx, c.MongoClient, c.MongoDatabase); err != nil {
			panic(err)
		}
	}

	genericServerOptions := server.NewServerOptions().
		WithEnvironment(c.Environment).
		WithMongoDatabase(c.MongoDatabase).
		WithMongoPassword(c.MongoPassword).
		WithMongoUsername(c.MongoUsername).
		WithMongoHosts([]string{"localhost:27017"})

	// Create IAM Server
	iamServer, err := iam.NewServer(ctx, genericServerOptions)
	if err != nil {
		panic(err)
	}
	if err := iamServer.Init(ctx); err != nil {
		panic(err)
	}
	router.PathPrefix("/apis/iam").Handler(iamServer)

	loginOptions := login.NewServerOptions()
	loginOptions.MongoUsername = c.MongoUsername
	loginOptions.MongoPassword = c.MongoPassword
	loginOptions.MongoPassword = c.MongoPassword
	loginOptions.MongoDatabase = c.MongoDatabase
	loginOptions.MongoHosts = []string{"localhost:27017"}
	loginOptions.BCryptCost = 15
	loginOptions.HydraAdminURL = "http://localhost:4445"
	loginServer, err := login.NewServer(ctx, loginOptions)
	if err != nil {
		panic(err)
	}
	router.PathPrefix("/auth").Handler(loginServer)
	router.PathPrefix("/apis/login").Handler(loginServer)

	// Create CMS Server
	cmsServer, err := cms.NewServer(ctx, genericServerOptions)
	if err != nil {
		panic(err)
	}
	router.PathPrefix("/apis/cms").Handler(cmsServer)

	// Create Webapp
	// WebApp
	webAppOptions := webapp2.ServerOptions{
		TemplateDirectory:       c.TemplateDirectory,
		RedisMaxIdleConnections: c.RedisMaxIdleConns,
		RedisAddress:            c.RedisAddress,
		RedisNetwork:            c.RedisNetwork,
		Environment:             c.Environment,
	}
	webappServer, err := webapp2.NewServer(
		webAppOptions,
		c.HydraAdminClient,
		c.HydraPublicClient,
	)
	if err != nil {
		panic(err)
	}
	router.PathPrefix("/").Handler(webappServer)

	httpServer := &http.Server{
		Addr:    ":9000",
		Handler: router,
	}

	srv := &Server{
		MongoClient:       c.MongoClient,
		WebAppHandler:     webappServer,
		HydraPublicClient: c.HydraPublicClient,
		HydraAdminClient:  c.HydraAdminClient,
		HttpServer:        httpServer,
	}

	go func() {
		listenAddress := c.Address
		listenAddress = strings.Replace(listenAddress, "https://", "", -1)
		listenAddress = strings.Replace(listenAddress, "http://", "", -1)
		_, port, err := net.SplitHostPort(listenAddress)
		if err != nil {
			panic(err)
		}
		if err := http.ListenAndServe(":"+port, router); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			panic(err)
		}
	}()

	if err := webappServer.Init(ctx); err != nil {
		panic(err)
	}

	if err := seeder.Seed(ctx, c.MongoClient, c.MongoDatabase); err != nil {
		panic(err)
	}

	return srv

}
