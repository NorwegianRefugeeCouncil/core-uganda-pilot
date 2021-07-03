package server

import (
	"context"
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/seeder"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/middleware"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/public"
	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/url"
	"strings"
)

type Options struct {

	// Seed
	ClearDB bool

	// Serve
	Environment   string
	ListenAddress string
	BaseURL       string

	// Mongo
	MongoDatabase string
	MongoUsername string
	MongoPassword string
	MongoHosts    []string

	// Redis
	RedisMaxIdleConnections int
	RedisNetwork            string
	RedisAddress            string
	RedisPassword           string
	RedisSecretKey          string

	// Hydra
	HydraAdminURL  string
	HydraPublicURL string

	// Web App
	WebAppTemplateDirectory string
	WebAppBasePath          string
	WebAppClientID          string
	WebAppClientSecret      string
	WebAppClientName        string
	WebAppIAMScheme         string
	WebAppIAMHost           string
	WebAppCMSScheme         string
	WebAppCMSHost           string

	// CMS
	CMSBasePath string

	// IAM
	IAMBasePath string

	// Login
	LoginBasePath          string
	LoginClientName        string
	LoginClientID          string
	LoginClientSecret      string
	LoginTemplateDirectory string
	LoginIAMHost           string
	LoginIAMScheme         string
}

func NewOptions() *Options {
	defaultEnv := "Production"
	defaultRedisAddress := "localhost:6379"
	defaultMongoHosts := []string{"http://localhost:27017"}
	defaultHydraAdminURL := "http://localhost:4445"
	defaultHydraPublicURL := "http://localhost:4444"
	defaultHost := "localhost"
	defaultScheme := "http"
	defaultPort := "9000"
	defaultUrl := url.URL{
		Scheme: defaultScheme,
		Host:   defaultHost + ":" + defaultPort,
	}
	return &Options{
		ClearDB:                 false,
		Environment:             defaultEnv,
		ListenAddress:           ":" + defaultUrl.Port(),
		BaseURL:                 defaultUrl.String(),
		MongoDatabase:           "core",
		MongoUsername:           "",
		MongoPassword:           "",
		MongoHosts:              defaultMongoHosts,
		RedisMaxIdleConnections: 10,
		RedisNetwork:            "tcp",
		RedisAddress:            defaultRedisAddress,
		RedisPassword:           "",
		RedisSecretKey:          "some-secret",
		HydraAdminURL:           defaultHydraAdminURL,
		HydraPublicURL:          defaultHydraPublicURL,
		WebAppTemplateDirectory: "pkg/apps/webapp/templates",
		WebAppBasePath:          "",
		WebAppClientID:          "webapp",
		WebAppClientSecret:      "",
		WebAppClientName:        "webapp",
		WebAppIAMScheme:         defaultUrl.Scheme,
		WebAppIAMHost:           defaultUrl.Host,
		WebAppCMSScheme:         defaultUrl.Scheme,
		WebAppCMSHost:           defaultUrl.Host,
		CMSBasePath:             "/apis/cms",
		IAMBasePath:             "/apis/iam",
		LoginBasePath:           "/auth",
		LoginClientName:         "login",
		LoginClientID:           "login",
		LoginClientSecret:       "",
		LoginTemplateDirectory:  "",
		LoginIAMScheme:          defaultUrl.Scheme,
		LoginIAMHost:            defaultUrl.Host,
	}
}

func (o *Options) Flags(fs *pflag.FlagSet) {

	// Seed
	fs.BoolVar(&o.ClearDB, "fresh", o.ClearDB, "Clear user-created DB entries")

	// Serve
	fs.StringVar(&o.Environment, "environment", o.Environment, "Environment (Production / Development)")
	fs.StringVar(&o.ListenAddress, "listen-address", o.ListenAddress, "Listen Address")
	fs.StringVar(&o.BaseURL, "base-url", o.BaseURL, "Base URL")

	// Mongo
	fs.StringVar(&o.MongoDatabase, "mongo-database", o.MongoDatabase, "Mongo database name")
	fs.StringVar(&o.MongoUsername, "mongo-username", o.MongoUsername, "Mongo username")
	fs.StringVar(&o.MongoPassword, "mongo-password", o.MongoPassword, "Mongo password")
	fs.StringSliceVar(&o.MongoHosts, "mongo-hosts", o.MongoHosts, "Mongo hosts")

	// Redis
	fs.IntVar(&o.RedisMaxIdleConnections, "redis-max-idle-conns", o.RedisMaxIdleConnections, "Redis maximum number of idle connections")
	fs.StringVar(&o.RedisAddress, "redis-address", o.RedisAddress, "Redis address")
	fs.StringVar(&o.RedisNetwork, "redis-network", o.RedisNetwork, "Redis network")
	fs.StringVar(&o.RedisPassword, "redis-password", o.RedisPassword, "Redis password")
	fs.StringVar(&o.RedisSecretKey, "redis-secret-key", o.RedisSecretKey, "Redis secret key")

	// Hydra
	fs.StringVar(&o.HydraAdminURL, "hydra-admin-url", o.HydraAdminURL, "Hydra Admin URL")
	fs.StringVar(&o.HydraPublicURL, "hydra-public-url", o.HydraPublicURL, "Hydra Public URL")

	// Login
	fs.StringVar(&o.LoginTemplateDirectory, "login-template-directory", o.LoginTemplateDirectory, "Template directory for login module")
	fs.StringVar(&o.LoginBasePath, "login-base-path", o.LoginBasePath, "Base path for the login module")
	fs.StringVar(&o.LoginClientName, "login-client-name", o.LoginClientName, "Login OAuth client name")
	fs.StringVar(&o.LoginClientID, "login-client-id", o.LoginClientID, "Login OAuth client ID")
	fs.StringVar(&o.LoginClientSecret, "login-client-secret", o.LoginClientSecret, "Login OAuth client secret")
	fs.StringVar(&o.LoginIAMHost, "login-iam-host", o.LoginIAMHost, "Login IAM Host")
	fs.StringVar(&o.LoginIAMScheme, "login-iam-scheme", o.LoginIAMScheme, "Login IAM Scheme")

	// IAM
	fs.StringVar(&o.IAMBasePath, "iam-base-path", o.IAMBasePath, "Base path for the IAM module")

	// CMS
	fs.StringVar(&o.CMSBasePath, "cms-base-path", o.CMSBasePath, "Base path for the CMS module")

	// Web App
	fs.StringVar(&o.WebAppBasePath, "web-base-path", o.WebAppBasePath, "Base path for the Web module")
	fs.StringVar(&o.WebAppTemplateDirectory, "web-template-directory", o.WebAppTemplateDirectory, "Directory for the web app templates")
	fs.StringVar(&o.WebAppClientID, "web-client-id", o.WebAppClientID, "Web app OAuth2 client ID")
	fs.StringVar(&o.WebAppClientSecret, "web-client-secret", o.WebAppClientSecret, "Web app OAuth2 client secret")
	fs.StringVar(&o.WebAppClientName, "web-client-name", o.WebAppClientName, "Web app OAuth2 client name")
	fs.StringVar(&o.WebAppIAMScheme, "web-iam-scheme", o.WebAppIAMScheme, "Web app IAM scheme")
	fs.StringVar(&o.WebAppIAMHost, "web-iam-host", o.WebAppIAMHost, "Web app IAM host")
	fs.StringVar(&o.WebAppCMSScheme, "web-cms-scheme", o.WebAppCMSScheme, "Web app CMS scheme")
	fs.StringVar(&o.WebAppCMSHost, "web-cms-host", o.WebAppCMSHost, "Web app CMS host")
}

type CompletedOptions struct {
	*Options
	MongoClient        *mongo.Client
	HydraAdminClient   *client.OryHydra
	HydraPublicClient  *client.OryHydra
	RedisPool          *redis.Pool
	OAuthTokenEndpoint string
	OIDCProvider       *oidc.Provider
}

func (o *Options) Complete(ctx context.Context) (CompletedOptions, error) {

	issuerUrl := o.HydraPublicURL
	if !strings.HasSuffix(issuerUrl, "/") {
		issuerUrl = issuerUrl + "/"
	}

	mongoClient, err := MongoClient(o.MongoHosts, o.MongoUsername, o.MongoPassword)
	if err != nil {
		return CompletedOptions{}, err
	}
	if err := mongoClient.Connect(ctx); err != nil {
		return CompletedOptions{}, err
	}

	hydraAdminClient, err := HydraAdminClient(o.HydraAdminURL)
	if err != nil {
		return CompletedOptions{}, err
	}

	hydraPublicCLient, err := HydraAdminClient(o.HydraPublicURL)
	if err != nil {
		return CompletedOptions{}, err
	}

	openIdConf, err := hydraPublicCLient.Public.DiscoverOpenIDConfiguration(&public.DiscoverOpenIDConfigurationParams{
		Context: ctx,
	})
	if err != nil {
		panic(err)
	}

	oidcProvider, err := oidc.NewProvider(ctx, issuerUrl)
	if err != nil {
		panic(err)
	}

	pool := &redis.Pool{
		MaxIdle: o.RedisMaxIdleConnections,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(o.RedisNetwork, o.RedisAddress)
		},
	}

	completedOptions := CompletedOptions{
		Options:            o,
		MongoClient:        mongoClient,
		HydraAdminClient:   hydraAdminClient,
		HydraPublicClient:  hydraPublicCLient,
		RedisPool:          pool,
		OAuthTokenEndpoint: *openIdConf.Payload.TokenEndpoint,
		OIDCProvider:       oidcProvider,
	}
	return completedOptions, nil
}

func HydraAdminClient(adminURL string) (*client.OryHydra, error) {
	hydraAdminURL, err := url.Parse(adminURL)
	if err != nil {
		return nil, err
	}
	hydraAdminClient := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Schemes: []string{
			hydraAdminURL.Scheme,
		},
		Host:     hydraAdminURL.Host,
		BasePath: hydraAdminURL.Path,
	})
	return hydraAdminClient, nil
}

func MongoClient(hosts []string, username, password string) (*mongo.Client, error) {
	// Setup mongo client
	mongoClient, err := mongo.NewClient(options.Client().
		SetHosts(hosts).
		SetAuth(
			options.Credential{
				Username: username,
				Password: password,
			}))
	if err != nil {
		return nil, err
	}
	return mongoClient, nil
}

func (c CompletedOptions) Generic() *server.GenericServerOptions {
	return &server.GenericServerOptions{
		MongoClient:       c.MongoClient,
		MongoDatabase:     c.MongoDatabase,
		Environment:       c.Environment,
		HydraAdminClient:  c.HydraAdminClient,
		HydraPublicClient: c.HydraPublicClient,
		RedisPool:         c.RedisPool,
	}
}

func (c CompletedOptions) New(ctx context.Context) *Server {

	// Prep db
	if c.ClearDB {
		if err := seeder.Clear(ctx, c.MongoClient, c.MongoDatabase); err != nil {
			panic(err)
		}
	}

	genericServerOptions := c.Generic()

	// Create IAM Server
	iamServer, err := c.CreateIAMServer(ctx, genericServerOptions)
	if err != nil {
		panic(err)
	}

	loginServer, err := c.CreateLoginServer(ctx, genericServerOptions)
	if err != nil {
		panic(err)
	}

	// Create CMS Server
	cmsServer, err := cms.NewServer(ctx, genericServerOptions)
	if err != nil {
		panic(err)
	}

	webAppServer, err := c.CreateWebAppServer(ctx, genericServerOptions)
	if err != nil {
		panic(err)
	}

	srv := &Server{
		MongoClient:       c.MongoClient,
		HydraPublicClient: c.HydraPublicClient,
		HydraAdminClient:  c.HydraAdminClient,
		WebAppServer:      webAppServer,
		IAMServer:         iamServer,
		LoginServer:       loginServer,
		CMSServer:         cmsServer,
	}

	router := c.CreateRouter(srv)
	srv.Router = router

	go func() {
		c.StartServer(srv)
	}()

	if err := seeder.Seed(ctx, c.MongoClient, c.MongoDatabase); err != nil {
		panic(err)
	}

	return srv

}

func (c CompletedOptions) CreateRouter(srv *Server) *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.UseLogging())
	router.PathPrefix("/apis/cms").Handler(srv.CMSServer)
	router.PathPrefix("/apis/iam").Handler(srv.IAMServer)
	router.PathPrefix("/auth").Handler(srv.LoginServer)
	router.PathPrefix("/apis/login").Handler(srv.LoginServer)
	router.PathPrefix("/").Handler(srv.WebAppServer)
	return router
}

func (c CompletedOptions) StartServer(server *Server) {
	if err := http.ListenAndServe(c.ListenAddress, server.Router); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		panic(err)
	}
}
