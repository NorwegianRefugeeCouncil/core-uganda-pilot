package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/seeder"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/middleware"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/public"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Options struct {

	// Seed
	ClearDB bool

	// Serve
	Environment   string
	ListenAddress string
	BaseURL       string
	TLSDisable    bool
	TLSCertPath   string
	TLSKeyPath    string

	// Mongo
	MongoDatabase     string
	MongoUsername     string
	MongoUsernameFile string
	MongoPassword     string
	MongoPasswordFile string
	MongoHosts        []string

	// Redis
	RedisMaxIdleConnections int
	RedisNetwork            string
	RedisAddress            string
	RedisPassword           string
	RedisPasswordFile       string
	RedisSecretKey          string

	// Hydra
	HydraAdminURL  string
	HydraPublicURL string

	// Web App
	WebAppTemplateDirectory string
	WebAppBasePath          string
	WebAppClientID          string
	WebAppClientSecret      string
	WebAppClientSecretFile  string
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
	LoginClientSecretFile  string
	LoginTemplateDirectory string
	LoginIAMHost           string
	LoginIAMScheme         string
	SeedDB                 bool
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
		LoginTemplateDirectory:  "pkg/apps/login/templates",
		LoginIAMScheme:          defaultUrl.Scheme,
		LoginIAMHost:            defaultUrl.Host,
	}
}

func (o *Options) Flags(fs *pflag.FlagSet) {

	// Seed
	fs.BoolVar(&o.ClearDB, "fresh", o.ClearDB, "Clear user-created DB entries")
	fs.BoolVar(&o.SeedDB, "seed", o.SeedDB, "Seed database with mock data")

	// Serve
	fs.StringVar(&o.Environment, "environment", o.Environment, "Environment (Production / Development)")
	fs.StringVar(&o.ListenAddress, "listen-address", o.ListenAddress, "Listen Address")
	fs.StringVar(&o.BaseURL, "base-url", o.BaseURL, "Base URL")
	fs.StringVar(&o.TLSKeyPath, "tls-key-path", o.BaseURL, "TLS Key Path")
	fs.StringVar(&o.TLSCertPath, "tls-cert-path", o.BaseURL, "TLS Cert Path")
	fs.BoolVar(&o.TLSDisable, "tls-disable", o.TLSDisable, "Disable TLS")

	// Mongo
	fs.StringVar(&o.MongoDatabase, "mongo-database", o.MongoDatabase, "Mongo database name")
	fs.StringVar(&o.MongoUsername, "mongo-username", o.MongoUsername, "Mongo username")
	fs.StringVar(&o.MongoUsernameFile, "mongo-username-file", o.MongoUsernameFile, "Mongo username file")
	fs.StringVar(&o.MongoPassword, "mongo-password", o.MongoPassword, "Mongo password")
	fs.StringVar(&o.MongoPasswordFile, "mongo-password-file", o.MongoPasswordFile, "Mongo password file")
	fs.StringSliceVar(&o.MongoHosts, "mongo-hosts", o.MongoHosts, "Mongo hosts")

	// Redis
	fs.IntVar(&o.RedisMaxIdleConnections, "redis-max-idle-conns", o.RedisMaxIdleConnections, "Redis maximum number of idle connections")
	fs.StringVar(&o.RedisAddress, "redis-address", o.RedisAddress, "Redis address")
	fs.StringVar(&o.RedisNetwork, "redis-network", o.RedisNetwork, "Redis network")
	fs.StringVar(&o.RedisPassword, "redis-password", o.RedisPassword, "Redis password file")
	fs.StringVar(&o.RedisPasswordFile, "redis-password-file", o.RedisPasswordFile, "Redis password")
	fs.StringVar(&o.RedisSecretKey, "redis-secret-key", o.RedisSecretKey, "Redis secret key")

	// Hydra
	fs.StringVar(&o.HydraAdminURL, "hydra-admin-url", o.HydraAdminURL, "Hydra Admin URL")
	fs.StringVar(&o.HydraPublicURL, "hydra-public-url", o.HydraPublicURL, "Hydra Public URL")

	// Login
	fs.StringVar(&o.LoginTemplateDirectory, "login-templates-directory", o.LoginTemplateDirectory, "Template directory for login module")
	fs.StringVar(&o.LoginBasePath, "login-base-path", o.LoginBasePath, "Base path for the login module")
	fs.StringVar(&o.LoginClientName, "login-client-name", o.LoginClientName, "Login OAuth client name")
	fs.StringVar(&o.LoginClientID, "login-client-id", o.LoginClientID, "Login OAuth client ID")
	fs.StringVar(&o.LoginClientSecret, "login-client-secret", o.LoginClientSecret, "Login OAuth client secret")
	fs.StringVar(&o.LoginClientSecretFile, "login-client-secret-file", o.LoginClientSecretFile, "Login OAuth client secret file")
	fs.StringVar(&o.LoginIAMHost, "login-iam-host", o.LoginIAMHost, "Login IAM Host")
	fs.StringVar(&o.LoginIAMScheme, "login-iam-scheme", o.LoginIAMScheme, "Login IAM Scheme")

	// IAM
	fs.StringVar(&o.IAMBasePath, "iam-base-path", o.IAMBasePath, "Base path for the IAM module")

	// CMS
	fs.StringVar(&o.CMSBasePath, "cms-base-path", o.CMSBasePath, "Base path for the CMS module")

	// Web App
	fs.StringVar(&o.WebAppBasePath, "web-base-path", o.WebAppBasePath, "Base path for the Web module")
	fs.StringVar(&o.WebAppTemplateDirectory, "web-templates-directory", o.WebAppTemplateDirectory, "Directory for the web app templates")
	fs.StringVar(&o.WebAppClientID, "web-client-id", o.WebAppClientID, "Web app OAuth2 client ID")
	fs.StringVar(&o.WebAppClientSecret, "web-client-secret", o.WebAppClientSecret, "Web app OAuth2 client secret")
	fs.StringVar(&o.WebAppClientSecretFile, "web-client-secret-file", o.WebAppClientSecretFile, "Web app OAuth2 client secret file")
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
	HydraTLSClient     *http.Client
}

func readFile(path string) (string, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(fileBytes), nil
}

func (o *Options) Complete(ctx context.Context) (CompletedOptions, error) {

	issuerUrl := o.HydraPublicURL
	if !strings.HasSuffix(issuerUrl, "/") {
		issuerUrl = issuerUrl + "/"
	}

	var err error

	if len(o.MongoUsername) == 0 && len(o.MongoUsernameFile) > 0 {
		o.MongoUsername, err = readFile(o.MongoUsernameFile)
		if err != nil {
			return CompletedOptions{}, err
		}
	}

	if len(o.MongoPassword) == 0 && len(o.MongoPasswordFile) > 0 {
		o.MongoPassword, err = readFile(o.MongoPasswordFile)
		if err != nil {
			return CompletedOptions{}, err
		}
	}

	if len(o.RedisPassword) == 0 && len(o.RedisPasswordFile) > 0 {
		o.RedisPassword, err = readFile(o.RedisPasswordFile)
		if err != nil {
			return CompletedOptions{}, err
		}
	}

	if len(o.LoginClientSecret) == 0 && len(o.LoginClientSecretFile) > 0 {
		o.LoginClientSecret, err = readFile(o.LoginClientSecretFile)
		if err != nil {
			return CompletedOptions{}, err
		}
	}

	if len(o.WebAppClientSecret) == 0 && len(o.WebAppClientSecretFile) > 0 {
		o.WebAppClientSecret, err = readFile(o.WebAppClientSecretFile)
		if err != nil {
			return CompletedOptions{}, err
		}
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

	hydraPublicClient, err := HydraAdminClient(o.HydraPublicURL)
	if err != nil {
		return CompletedOptions{}, err
	}

	hydraHttpClient := http.DefaultClient
	if !o.TLSDisable {
		hydraHttpClient, err = tlsClient(o.TLSCertPath)
		if err != nil {
			return CompletedOptions{}, err
		}
	}

	openIdConf, err := hydraPublicClient.Public.DiscoverOpenIDConfiguration(&public.DiscoverOpenIDConfigurationParams{
		Context:    ctx,
		HTTPClient: hydraHttpClient,
	})
	if err != nil {
		panic(err)
	}

	oidcCtx := oidc.ClientContext(ctx, hydraHttpClient)
	oidcProvider, err := oidc.NewProvider(oidcCtx, issuerUrl)
	if err != nil {
		panic(err)
	}

	pool := &redis.Pool{
		MaxIdle: o.RedisMaxIdleConnections,
		Dial: func() (redis.Conn, error) {
			var redisOptions []redis.DialOption
			if len(o.RedisPassword) > 0 {
				redisOptions = append(redisOptions, redis.DialPassword(o.RedisPassword))
			}
			return redis.Dial(o.RedisNetwork, o.RedisAddress, redisOptions...)
		},
	}

	completedOptions := CompletedOptions{
		Options:            o,
		MongoClient:        mongoClient,
		HydraAdminClient:   hydraAdminClient,
		HydraPublicClient:  hydraPublicClient,
		HydraTLSClient:     hydraHttpClient,
		RedisPool:          pool,
		OAuthTokenEndpoint: *openIdConf.Payload.TokenEndpoint,
		OIDCProvider:       oidcProvider,
	}
	return completedOptions, nil
}

func tlsClient(tlsCertPath string) (*http.Client, error) {
	certFile, err := ioutil.ReadFile(tlsCertPath)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(certFile); !ok {
		err := fmt.Errorf("failed to append cert to CertPool")
		logrus.WithError(err).Errorf("")
		return nil, err
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}
	return httpClient, nil
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
		HydraHTTPClient:   c.HydraTLSClient,
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
		logrus.WithError(err).Errorf("failed to create IAM server")
		panic(err)
	}

	loginServer, err := c.CreateLoginServer(ctx, genericServerOptions)
	if err != nil {
		logrus.WithError(err).Errorf("failed to create login server")
		panic(err)
	}

	// Create CMS Server
	cmsServer, err := cms.NewServer(ctx, genericServerOptions)
	if err != nil {
		logrus.WithError(err).Errorf("failed to create CMS server")
		panic(err)
	}

	webAppServer, err := c.CreateWebAppServer(ctx, genericServerOptions)
	if err != nil {
		logrus.WithError(err).Errorf("failed to create WebApp server")
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

	if c.SeedDB {
		if err := seeder.Seed(ctx, c.MongoClient, c.MongoDatabase); err != nil {
			panic(err)
		}
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
	if c.TLSDisable {
		if err := http.ListenAndServe(c.ListenAddress, server.Router); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			panic(err)
		}
	} else {
		if err := http.ListenAndServeTLS(c.ListenAddress, c.TLSCertPath, c.TLSKeyPath, server.Router); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			panic(err)
		}
	}
}
