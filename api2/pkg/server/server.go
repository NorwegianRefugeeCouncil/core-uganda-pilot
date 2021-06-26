package server

import (
	"context"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/apps/cms"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"github.com/nrc-no/core-kafka/pkg/apps/seed"
	"github.com/nrc-no/core-kafka/pkg/auth"
	"github.com/nrc-no/core-kafka/pkg/middleware"
	"github.com/nrc-no/core-kafka/pkg/rest"
	"github.com/nrc-no/core-kafka/pkg/sessionmanager"
	"github.com/nrc-no/core-kafka/pkg/webapp"
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
	WebAppHandler     *webapp.Handler
	HttpServer        *http.Server
	SessionManager    sessionmanager.Store
	HydraPublicClient *client.OryHydra
	HydraAdminClient  *client.OryHydra
	CredentialsClient *auth.CredentialsClient
}

type Options struct {
	TemplateDirectory    string
	Address              string
	MongoDatabase        string
	MongoUsername        string
	MongoPassword        string
	KeycloakClientID     string
	KeycloakClientSecret string
	KeycloakBaseURL      string
	KeycloakRealmName    string
	RedisMaxIdleConns    int
	RedisNetwork         string
	RedisAddress         string
	RedisPassword        string
	RedisSecretKey       string
	HydraAdminURL        string
	HydraPublicURL       string
}

func NewOptions() *Options {
	return &Options{
		TemplateDirectory:    "pkg/webapp/templates",
		Address:              "http://localhost:9000",
		MongoDatabase:        "core",
		MongoUsername:        "",
		MongoPassword:        "",
		KeycloakClientID:     "",
		KeycloakClientSecret: "",
		KeycloakBaseURL:      "",
		KeycloakRealmName:    "",
		RedisMaxIdleConns:    10,
		RedisNetwork:         "tcp",
		RedisAddress:         "localhost:6379",
		RedisPassword:        "",
		RedisSecretKey:       "some-secret",
		HydraAdminURL:        "http://localhost:4445",
		HydraPublicURL:       "http://localhost:4444",
	}
}

func (o *Options) Flags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Address, "address", o.Address, "Address")
	fs.StringVar(&o.MongoDatabase, "mongo-database", o.MongoDatabase, "Mongo database name")
	fs.StringVar(&o.MongoUsername, "mongo-username", o.MongoUsername, "Mongo username")
	fs.StringVar(&o.MongoPassword, "mongo-password", o.MongoPassword, "Mongo password")
	fs.StringVar(&o.KeycloakBaseURL, "keycloak-base-url", o.KeycloakBaseURL, "Keycloak base URL")
	fs.StringVar(&o.KeycloakRealmName, "keycloak-realm-name", o.KeycloakRealmName, "Keycloak realm name")
	fs.StringVar(&o.KeycloakClientID, "keycloak-client-id", o.KeycloakClientID, "Keycloak client id")
	fs.StringVar(&o.KeycloakClientSecret, "keycloak-client-secret", o.KeycloakClientSecret, "Keycloak client secret")
	fs.IntVar(&o.RedisMaxIdleConns, "redis-max-idle-conns", o.RedisMaxIdleConns, "Redis maximum number of idle connections")
	fs.StringVar(&o.RedisNetwork, "redis-network", o.RedisNetwork, "Redis network")
	fs.StringVar(&o.RedisPassword, "redis-password", o.RedisPassword, "Redis password")
	fs.StringVar(&o.RedisSecretKey, "redis-secret-key", o.RedisSecretKey, "Redis secret key")
	fs.StringVar(&o.HydraAdminURL, "hydra-admin-url", o.HydraAdminURL, "Hydra Admin URL")
	fs.StringVar(&o.HydraPublicURL, "hydra-public-url", o.HydraPublicURL, "Hydra Public URL")
}

type CompletedOptions struct {
	*Options
	MongoClient       *mongo.Client
	SessionManager    sessionmanager.Store
	HydraAdminClient  *client.OryHydra
	HydraPublicClient *client.OryHydra
	CredentialsClient *auth.CredentialsClient
}

func (o *Options) Complete(ctx context.Context) (CompletedOptions, error) {

	pool := &redis.Pool{
		MaxIdle: o.RedisMaxIdleConns,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(o.RedisNetwork, o.RedisAddress)
		},
	}

	sm := sessionmanager.New(pool, sessionmanager.Options{})

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

	credentialsClient := auth.NewCredentialsClient(o.MongoDatabase, mongoClient)

	completedOptions := CompletedOptions{
		Options:           o,
		MongoClient:       mongoClient,
		HydraAdminClient:  hydraAdminClient,
		HydraPublicClient: hydraPublicCLient,
		CredentialsClient: credentialsClient,
		SessionManager:    sm,
	}
	return completedOptions, nil
}

func (c CompletedOptions) New(ctx context.Context) *Server {

	router := mux.NewRouter()
	router.Use(middleware.UseLogging())
	router.Use(c.SessionManager.LoadAndSave)

	addressUrl, err := url.Parse(c.Address)
	if err != nil {
		panic(err)
	}

	iamServer, err := iam.NewServer(
		ctx,
		iam.NewServerOptions().
			WithMongoDatabase(c.MongoDatabase).
			WithMongoUsername(c.MongoUsername).
			WithMongoPassword(c.MongoPassword).
			WithMongoHosts([]string{"localhost:27017"}))
	if err != nil {
		panic(err)
	}
	if err := iamServer.Init(ctx); err != nil {
		panic(err)
	}
	router.PathPrefix("/apis/iam").Handler(iamServer)
	iamClient := iam.NewClientSet(&rest.RESTConfig{
		Scheme: addressUrl.Scheme,
		Host:   addressUrl.Host,
	})

	cmsServer, err := cms.NewServer(ctx, cms.NewServerOptions().
		WithMongoDatabase(c.MongoDatabase).
		WithMongoUsername(c.MongoUsername).
		WithMongoPassword(c.MongoPassword).
		WithMongoHosts([]string{"localhost:27017"}))
	if err != nil {
		panic(err)
	}
	router.PathPrefix("/apis/cms").Handler(cmsServer)
	cmsClient := cms.NewClientSet(&rest.RESTConfig{
		Scheme: addressUrl.Scheme,
		Host:   addressUrl.Host,
	})

	// WebApp
	webAppOptions := webapp.Options{
		TemplateDirectory: c.TemplateDirectory,
	}
	webAppHandler, err := webapp.NewHandler(
		webAppOptions,
		c.HydraAdminClient,
		c.HydraPublicClient,
		c.SessionManager,
		c.CredentialsClient,
		iamClient,
		cmsClient,
	)
	if err != nil {
		panic(err)
	}

	router.Path("/").HandlerFunc(webAppHandler.Individuals)
	router.Path("/individuals").HandlerFunc(webAppHandler.Individuals)
	router.Path("/individuals/{id}").HandlerFunc(webAppHandler.Individual)
	router.Path("/individuals/{id}/credentials").HandlerFunc(webAppHandler.IndividualCredentials)
	router.Path("/teams").HandlerFunc(webAppHandler.Teams)
	router.Path("/teams/{id}").HandlerFunc(webAppHandler.Team)
	router.Path("/cases").HandlerFunc(webAppHandler.Cases)
	router.Path("/cases/new").HandlerFunc(webAppHandler.NewCase)
	router.Path("/cases/{id}").HandlerFunc(webAppHandler.Case)
	router.Path("/settings").HandlerFunc(webAppHandler.Settings)
	router.Path("/settings/attributes").HandlerFunc(webAppHandler.Attributes)
	router.Path("/settings/attributes/new").HandlerFunc(webAppHandler.NewAttribute)
	router.Path("/settings/attributes/{id}").HandlerFunc(webAppHandler.Attribute)
	router.Path("/settings/relationshiptypes").HandlerFunc(webAppHandler.RelationshipTypes)
	router.Path("/settings/relationshiptypes/new").HandlerFunc(webAppHandler.NewRelationshipType)
	router.Path("/settings/relationshiptypes/{id}").HandlerFunc(webAppHandler.RelationshipType)
	router.Path("/settings/partytypes").HandlerFunc(webAppHandler.PartyTypes)
	router.Path("/settings/partytypes/{id}").HandlerFunc(webAppHandler.PartyType)
	router.Path("/settings/casetypes").HandlerFunc(webAppHandler.CaseTypes)
	router.Path("/settings/casetypes/new").HandlerFunc(webAppHandler.NewCaseType)
	router.Path("/settings/casetypes/{id}").HandlerFunc(webAppHandler.CaseType)
	router.Path("/settings/authclients").HandlerFunc(webAppHandler.AuthClients)
	router.Path("/settings/authclients/{id}").HandlerFunc(webAppHandler.AuthClient)
	router.Path("/settings/authclients/{id}/newsecret").HandlerFunc(webAppHandler.AuthClientNewSecret)
	router.Path("/settings/authclients/{id}/delete").HandlerFunc(webAppHandler.DeleteAuthClient)

	httpServer := &http.Server{
		Addr:    ":9000",
		Handler: router,
	}

	srv := &Server{
		MongoClient:       c.MongoClient,
		WebAppHandler:     webAppHandler,
		SessionManager:    c.SessionManager,
		HydraPublicClient: c.HydraPublicClient,
		HydraAdminClient:  c.HydraAdminClient,
		CredentialsClient: c.CredentialsClient,
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

	if err := seed.Seed(ctx, c.MongoDatabase, c.MongoClient); err != nil {
		panic(err)
	}

	return srv

}
