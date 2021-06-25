package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/auth"
	"github.com/nrc-no/core-kafka/pkg/cases/cases"
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	"github.com/nrc-no/core-kafka/pkg/individuals"
	"github.com/nrc-no/core-kafka/pkg/keycloak"
	"github.com/nrc-no/core-kafka/pkg/memberships"
	"github.com/nrc-no/core-kafka/pkg/middleware"
	"github.com/nrc-no/core-kafka/pkg/organizations"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypeschemas"
	"github.com/nrc-no/core-kafka/pkg/parties/relationships"
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	"github.com/nrc-no/core-kafka/pkg/relationshipparties"
	"github.com/nrc-no/core-kafka/pkg/sessionmanager"
	"github.com/nrc-no/core-kafka/pkg/staff"
	"github.com/nrc-no/core-kafka/pkg/staffmock"
	"github.com/nrc-no/core-kafka/pkg/teamorganizations"
	"github.com/nrc-no/core-kafka/pkg/teams"
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
	MongoClient               *mongo.Client
	AttributeStore            *attributes.Store
	AttributeHandler          *attributes.Handler
	AttributeClient           *attributes.Client
	IndividualStore           *individuals.Store
	IndividualHandler         *individuals.Handler
	IndividualClient          *individuals.Client
	RelationshipTypeStore     *relationshiptypes.Store
	RelationshipTypeHandler   *relationshiptypes.Handler
	RelationshipTypeClient    *relationshiptypes.Client
	RelationshipStore         *relationships.Store
	RelationshipHandler       *relationships.Handler
	RelationshipClient        *relationships.Client
	PartyStore                *parties.Store
	PartyHandler              *parties.Handler
	PartyClient               *parties.Client
	PartyTypeStore            *partytypes.Store
	PartyTypeHandler          *partytypes.Handler
	PartyTypeClient           *partytypes.Client
	RelationshipPartiesClient *relationshipparties.Client
	CaseTypeStore             *casetypes.Store
	CaseTypeHandler           *casetypes.Handler
	CaseTypeClient            *casetypes.Client
	CaseStore                 *cases.Store
	CaseHandler               *cases.Handler
	CaseClient                *cases.Client
	WebAppHandler             *webapp.Handler
	HttpServer                *http.Server
	TeamStore                 *teams.Store
	TeamHandler               *teams.Handler
	TeamClient                *teams.Client
	MembershipHandler         *memberships.Handler
	MembershipStore           *memberships.Store
	membershipsClient         *memberships.Client
	KeycloakClient            *keycloak.Client
	SessionManager            sessionmanager.Store
	HydraPublicClient         *client.OryHydra
	HydraAdminClient          *client.OryHydra
	CredentialsClient         *auth.CredentialsClient
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
	KeycloakClient       *keycloak.Client
	KeycloakClientID     string
	KeycloakClientSecret string
	KeycloakBaseURL      string
	KeycloakRealmName    string
	MongoClient          *mongo.Client
	TemplateDirectory    string
	Address              string
	MongoDatabase        string
	SessionManager       sessionmanager.Store
	HydraAdminClient     *client.OryHydra
	HydraPublicClient    *client.OryHydra
	CredentialsClient    *auth.CredentialsClient
}

func (o Options) Complete(ctx context.Context) (CompletedOptions, error) {

	pool := &redis.Pool{
		MaxIdle: o.RedisMaxIdleConns,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(o.RedisNetwork, o.RedisAddress)
		},
	}

	sm := sessionmanager.New(pool, sessionmanager.Options{})

	// Setup mongo client
	mongoClient, err := mongo.NewClient(options.Client().SetAuth(
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

	keycloakClient, err := keycloak.NewClient(o.KeycloakBaseURL, o.KeycloakRealmName, o.KeycloakClientID, o.KeycloakClientSecret)
	if err != nil {
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
		MongoClient:          mongoClient,
		TemplateDirectory:    o.TemplateDirectory,
		Address:              o.Address,
		MongoDatabase:        o.MongoDatabase,
		KeycloakClient:       keycloakClient,
		KeycloakClientSecret: o.KeycloakClientSecret,
		KeycloakClientID:     o.KeycloakClientID,
		KeycloakBaseURL:      o.KeycloakBaseURL,
		KeycloakRealmName:    o.KeycloakRealmName,
		HydraAdminClient:     hydraAdminClient,
		HydraPublicClient:    hydraPublicCLient,
		CredentialsClient:    credentialsClient,
		SessionManager:       sm,
	}
	return completedOptions, nil
}

func (c CompletedOptions) New(ctx context.Context) *Server {

	router := mux.NewRouter()

	router.Use(middleware.UseLogging())
	router.Use(c.SessionManager.LoadAndSave)

	// Attributes
	attributeStore := attributes.NewStore(c.MongoClient, c.MongoDatabase)
	if err := attributes.Init(ctx, attributeStore); err != nil {
		panic(err)
	}
	attributeHandler := attributes.NewHandler(attributeStore)
	attributeClient := attributes.NewClient(c.Address)
	router.Path("/apis/v1/attributes").Methods("GET").HandlerFunc(attributeHandler.List)
	router.Path("/apis/v1/attributes/{id}").Methods("GET").HandlerFunc(attributeHandler.Get)
	router.Path("/apis/v1/attributes/{id}").Methods("PUT").HandlerFunc(attributeHandler.Update)
	router.Path("/apis/v1/attributes").Methods("POST").HandlerFunc(attributeHandler.Post)

	// RelationshipTypes
	relationshipTypeStore := relationshiptypes.NewStore(c.MongoClient, c.MongoDatabase)
	if err := relationshiptypes.Init(ctx, relationshipTypeStore); err != nil {
		panic(err)
	}
	relationshipTypeHandler := relationshiptypes.NewHandler(relationshipTypeStore)
	relationshipTypeClient := relationshiptypes.NewClient(c.Address)
	router.Path("/apis/v1/relationshiptypes").Methods("GET").HandlerFunc(relationshipTypeHandler.List)
	router.Path("/apis/v1/relationshiptypes/{id}").Methods("GET").HandlerFunc(relationshipTypeHandler.Get)
	router.Path("/apis/v1/relationshiptypes/{id}").Methods("PUT").HandlerFunc(relationshipTypeHandler.Put)
	router.Path("/apis/v1/relationshiptypes").Methods("POST").HandlerFunc(relationshipTypeHandler.Post)

	// Relationships
	relationshipStore := relationships.NewStore(c.MongoClient, c.MongoDatabase)
	relationshipHandler := relationships.NewHandler(relationshipStore)
	relationshipClient := relationships.NewClient(c.Address)
	router.Path("/apis/v1/relationships").Methods("GET").HandlerFunc(relationshipHandler.List)
	router.Path("/apis/v1/relationships/{id}").Methods("GET").HandlerFunc(relationshipHandler.Get)
	router.Path("/apis/v1/relationships/{id}").Methods("PUT").HandlerFunc(relationshipHandler.Put)
	router.Path("/apis/v1/relationships").Methods("POST").HandlerFunc(relationshipHandler.Post)
	router.Path("/apis/v1/relationships/{id}").Methods("DELETE").HandlerFunc(relationshipHandler.Delete)

	// Parties
	partyStore := parties.NewStore(c.MongoClient, c.MongoDatabase)
	if err := parties.Init(ctx, partyStore); err != nil {
		panic(err)
	}
	partyHandler := parties.NewHandler(partyStore)
	partyClient := parties.NewClient(c.Address)
	router.Path("/apis/v1/parties").Methods("GET").HandlerFunc(partyHandler.List)
	router.Path("/apis/v1/parties/{id}").Methods("GET").HandlerFunc(partyHandler.Get)
	router.Path("/apis/v1/parties/{id}").Methods("PUT").HandlerFunc(partyHandler.Put)
	router.Path("/apis/v1/parties").Methods("POST").HandlerFunc(partyHandler.Post)

	// Party Types
	partyTypeStore := partytypes.NewStore(c.MongoClient, c.MongoDatabase)
	if err := partytypes.Init(ctx, partyTypeStore); err != nil {
		panic(err)
	}
	partyTypeHandler := partytypes.NewHandler(partyTypeStore)
	partyTypeClient := partytypes.NewClient(c.Address)
	router.Path("/apis/v1/partytypes").Methods("GET").HandlerFunc(partyTypeHandler.List)
	router.Path("/apis/v1/partytypes/{id}").Methods("GET").HandlerFunc(partyTypeHandler.Get)
	router.Path("/apis/v1/partytypes/{id}").Methods("PUT").HandlerFunc(partyTypeHandler.Put)
	router.Path("/apis/v1/partytypes").Methods("POST").HandlerFunc(partyTypeHandler.Post)

	// PartyTypeSchemas
	partyTypeSchemaStore := partytypeschemas.NewStore(c.MongoClient, c.MongoDatabase)
	partyTypeSchemaHandler := partytypeschemas.NewHandler(partyTypeSchemaStore)
	// TOOD: partyTypeSchemaClient := partytypeschemas.NewClient(serverOptions.Address)
	router.Path("/apis/v1/partytypeschemas").Methods("GET").HandlerFunc(partyTypeSchemaHandler.List)
	router.Path("/apis/v1/partytypeschemas/{id}").Methods("GET").HandlerFunc(partyTypeSchemaHandler.Get)
	router.Path("/apis/v1/partytypeschemas/{id}").Methods("PUT").HandlerFunc(partyTypeSchemaHandler.Put)
	router.Path("/apis/v1/partytypeschemas").Methods("POST").HandlerFunc(partyTypeSchemaHandler.Post)

	// Relationship <> Parties
	relationshipPartiesStore := relationshipparties.NewStore(partyStore)
	relationshipPartiesHandler := relationshipparties.NewHandler(relationshipPartiesStore)
	relationshipPartiesClient := relationshipparties.NewClient(c.Address)
	router.Path("/apis/v1/relationshipparties/picker").Methods("GET").HandlerFunc(relationshipPartiesHandler.PickParty)

	// individuals
	individualsStore := individuals.NewStore(c.MongoClient, c.MongoDatabase)
	if err := individuals.Init(ctx, attributeStore, partyStore); err != nil {
		panic(err)
	}
	individualHandler := individuals.NewHandler(individualsStore)
	individualClient := individuals.NewClient(c.Address)
	if err := individuals.SeedDatabase(ctx, individualsStore); err != nil {
		panic(err)
	}
	router.Path("/apis/v1/individuals").Methods("GET").HandlerFunc(individualHandler.List)
	router.Path("/apis/v1/individuals/{id}").Methods("GET").HandlerFunc(individualHandler.Get)
	router.Path("/apis/v1/individuals/{id}").Methods("PUT").HandlerFunc(individualHandler.Update)
	router.Path("/apis/v1/individuals").Methods("POST").HandlerFunc(individualHandler.Create)

	// Teams
	teamStore := teams.NewStore(partyStore)
	if err := teams.Init(ctx, teamStore, partyTypeStore, attributeStore); err != nil {
		panic(err)
	}
	teamHandler := teams.NewHandler(teamStore)
	teamClient := teams.NewClient(c.Address)
	router.Path("/apis/v1/teams").Methods("GET").HandlerFunc(teamHandler.List)
	router.Path("/apis/v1/teams/{id}").Methods("GET").HandlerFunc(teamHandler.Get)
	router.Path("/apis/v1/teams/{id}").Methods("PUT").HandlerFunc(teamHandler.Put)
	router.Path("/apis/v1/teams").Methods("POST").HandlerFunc(teamHandler.Post)

	// Staff
	staffStore := staff.NewStore(relationshipStore)
	if err := staff.Init(ctx, relationshipTypeStore); err != nil {
		panic(err)
	}

	// Memberships
	membershipStore := memberships.NewStore(relationshipStore)
	membershipHandler := memberships.NewHandler(membershipStore)
	membershipClient := memberships.NewClient(c.Address)
	if err := memberships.Init(ctx, relationshipTypeStore); err != nil {
		panic(err)
	}
	router.Path("/apis/v1/memberships").Methods("GET").HandlerFunc(membershipHandler.List)
	router.Path("/apis/v1/memberships/{id}").Methods("GET").HandlerFunc(membershipHandler.Get)
	router.Path("/apis/v1/memberships").Methods("POST").HandlerFunc(membershipHandler.Post)

	// Organizations
	if err := organizations.Init(ctx, partyTypeStore, attributeStore, partyStore); err != nil {
		panic(err)
	}

	// TeamOrganization
	if err := teamorganizations.Init(ctx, relationshipTypeStore, relationshipStore); err != nil {
		panic(err)
	}

	// Mock staff
	if err := staffmock.Init(ctx, partyStore, staffStore, membershipStore); err != nil {
		panic(err)
	}

	// Auth
	authHandler, err := auth.NewHandler(
		ctx,
		fmt.Sprintf("%s/auth/realms/%s", c.KeycloakBaseURL, c.KeycloakRealmName),
		c.KeycloakClientID,
		c.KeycloakClientSecret,
		"http://localhost:9000/auth/callback",
		c.SessionManager,
		c.KeycloakClient)
	if err != nil {
		panic(err)
	}
	if err := auth.Init(ctx, c.KeycloakClient); err != nil {
		panic(err)
	}
	// router.Use(authHandler.Authenticate())
	router.Path("/auth/login").Methods("GET").HandlerFunc(authHandler.Login)
	router.Path("/auth/logout").Methods("GET").HandlerFunc(authHandler.Logout)
	router.Path("/auth/callback").Methods("GET").HandlerFunc(authHandler.Callback)

	// Cases
	caseStore := cases.NewStore(c.MongoClient, c.MongoDatabase)
	if err := cases.Init(ctx, caseStore); err != nil {
		panic(err)
	}
	caseHandler := cases.NewHandler(caseStore)
	caseClient := cases.NewClient(c.Address)
	router.Path("/apis/v1/cases").Methods("GET").HandlerFunc(caseHandler.List)
	router.Path("/apis/v1/cases/{id}").Methods("GET").HandlerFunc(caseHandler.Get)
	router.Path("/apis/v1/cases/{id}").Methods("PUT").HandlerFunc(caseHandler.Put)
	router.Path("/apis/v1/cases").Methods("POST").HandlerFunc(caseHandler.Post)

	// CaseTypes
	caseTypeStore := casetypes.NewStore(c.MongoClient, c.MongoDatabase)
	if err := casetypes.Init(ctx, caseTypeStore); err != nil {
		panic(err)
	}
	caseTypeHandler := casetypes.NewHandler(caseTypeStore)
	caseTypeClient := casetypes.NewClient(c.Address)
	router.Path("/apis/v1/casetypes").Methods("GET").HandlerFunc(caseTypeHandler.List)
	router.Path("/apis/v1/casetypes/{id}").Methods("GET").HandlerFunc(caseTypeHandler.Get)
	router.Path("/apis/v1/casetypes/{id}").Methods("PUT").HandlerFunc(caseTypeHandler.Put)
	router.Path("/apis/v1/casetypes").Methods("POST").HandlerFunc(caseTypeHandler.Post)

	// WebApp
	webAppOptions := webapp.Options{
		TemplateDirectory: c.TemplateDirectory,
	}
	webAppHandler, err := webapp.NewHandler(
		webAppOptions,
		attributeClient,
		individualClient,
		relationshipTypeClient,
		relationshipClient,
		partyClient,
		partyTypeClient,
		caseTypeClient,
		caseClient,
		relationshipPartiesClient,
		teamClient,
		membershipClient,
		c.HydraAdminClient,
		c.HydraPublicClient,
		c.SessionManager,
		c.CredentialsClient,
		partyStore,
	)
	if err != nil {
		panic(err)
	}

	router.Path("/").HandlerFunc(webAppHandler.Individuals)
	router.Path("/individuals").HandlerFunc(webAppHandler.Individuals)
	router.Path("/individuals/{id}").HandlerFunc(webAppHandler.Individual)
	router.Path("/individuals/{id}/credentials").HandlerFunc(webAppHandler.IndividualCredentials)
	router.Path("/settings").HandlerFunc(webAppHandler.Settings)
	router.Path("/settings/attributes").HandlerFunc(webAppHandler.Attributes)
	router.Path("/settings/attributes/new").HandlerFunc(webAppHandler.NewAttribute)
	router.Path("/settings/attributes/{id}").HandlerFunc(webAppHandler.Attribute)
	router.Path("/settings/relationshiptypes").HandlerFunc(webAppHandler.RelationshipTypes)
	router.Path("/settings/relationshiptypes/new").HandlerFunc(webAppHandler.NewRelationshipType)
	router.Path("/settings/relationshiptypes/{id}").HandlerFunc(webAppHandler.RelationshipType)
	router.Path("/settings/partytypes").HandlerFunc(webAppHandler.PartyTypes)
	router.Path("/settings/partytypes/{id}").HandlerFunc(webAppHandler.PartyType)
	router.Path("/settings/countries").HandlerFunc(webAppHandler.CountrySettings)
	router.Path("/teams").HandlerFunc(webAppHandler.Teams)
	router.Path("/teams/{id}").HandlerFunc(webAppHandler.Team)
	router.Path("/cases").HandlerFunc(webAppHandler.Cases)
	router.Path("/cases/new").HandlerFunc(webAppHandler.NewCase)
	router.Path("/cases/{id}").HandlerFunc(webAppHandler.Case)
	router.Path("/settings/casetypes").HandlerFunc(webAppHandler.CaseTypes)
	router.Path("/settings/casetypes/new").HandlerFunc(webAppHandler.NewCaseType)
	router.Path("/settings/casetypes/{id}").HandlerFunc(webAppHandler.CaseType)
	router.Path("/settings/authclients").HandlerFunc(webAppHandler.AuthClients)
	router.Path("/settings/authclients/{id}").HandlerFunc(webAppHandler.AuthClient)
	router.Path("/settings/authclients/{id}/newsecret").HandlerFunc(webAppHandler.AuthClientNewSecret)
	router.Path("/settings/authclients/{id}/delete").HandlerFunc(webAppHandler.DeleteAuthClient)
	router.Path("/login").Methods("GET").HandlerFunc(webAppHandler.Login)
	router.Path("/login").Methods("POST").HandlerFunc(webAppHandler.PostLogin)

	// Seed database for development
	if err := individuals.SeedDatabase(ctx, individualsStore); err != nil {
		panic(err)
	}

	httpServer := &http.Server{
		Addr:    ":9000",
		Handler: router,
	}

	srv := &Server{
		MongoClient:               c.MongoClient,
		AttributeStore:            attributeStore,
		AttributeHandler:          attributeHandler,
		AttributeClient:           attributeClient,
		IndividualStore:           individualsStore,
		IndividualHandler:         individualHandler,
		IndividualClient:          individualClient,
		RelationshipTypeStore:     relationshipTypeStore,
		RelationshipTypeHandler:   relationshipTypeHandler,
		RelationshipTypeClient:    relationshipTypeClient,
		RelationshipStore:         relationshipStore,
		RelationshipHandler:       relationshipHandler,
		RelationshipClient:        relationshipClient,
		PartyStore:                partyStore,
		PartyHandler:              partyHandler,
		PartyClient:               partyClient,
		PartyTypeStore:            partyTypeStore,
		PartyTypeHandler:          partyTypeHandler,
		PartyTypeClient:           partyTypeClient,
		RelationshipPartiesClient: relationshipPartiesClient,
		CaseTypeStore:             caseTypeStore,
		CaseTypeHandler:           caseTypeHandler,
		CaseTypeClient:            caseTypeClient,
		CaseStore:                 caseStore,
		CaseHandler:               caseHandler,
		CaseClient:                caseClient,
		TeamStore:                 teamStore,
		TeamHandler:               teamHandler,
		TeamClient:                teamClient,
		MembershipStore:           membershipStore,
		MembershipHandler:         membershipHandler,
		membershipsClient:         membershipClient,
		WebAppHandler:             webAppHandler,
		KeycloakClient:            c.KeycloakClient,
		SessionManager:            c.SessionManager,
		HydraPublicClient:         c.HydraPublicClient,
		HydraAdminClient:          c.HydraAdminClient,
		CredentialsClient:         c.CredentialsClient,
		HttpServer:                httpServer,
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

	return srv

}
