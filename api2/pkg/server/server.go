package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/cases/cases"
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	"github.com/nrc-no/core-kafka/pkg/memberships"
	"github.com/nrc-no/core-kafka/pkg/organizations"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/beneficiaries"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypeschemas"
	"github.com/nrc-no/core-kafka/pkg/parties/relationships"
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	"github.com/nrc-no/core-kafka/pkg/staff"
	"github.com/nrc-no/core-kafka/pkg/staffmock"
	"github.com/nrc-no/core-kafka/pkg/teamorganizations"
	"github.com/nrc-no/core-kafka/pkg/teams"
	"github.com/nrc-no/core-kafka/pkg/webapp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"net"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	MongoClient             *mongo.Client
	AttributeStore          *attributes.Store
	AttributeHandler        *attributes.Handler
	AttributeClient         *attributes.Client
	BeneficiaryStore        *beneficiaries.Store
	BeneficiaryHandler      *beneficiaries.Handler
	BeneficiaryClient       *beneficiaries.Client
	RelationshipTypeStore   *relationshiptypes.Store
	RelationshipTypeHandler *relationshiptypes.Handler
	RelationshipTypeClient  *relationshiptypes.Client
	RelationshipStore       *relationships.Store
	RelationshipHandler     *relationships.Handler
	RelationshipClient      *relationships.Client
	PartyStore              *parties.Store
	PartyHandler            *parties.Handler
	PartyClient             *parties.Client
	PartyTypeStore          *partytypes.Store
	PartyTypeHandler        *partytypes.Handler
	PartyTypeClient         *partytypes.Client
	CaseTypeStore           *casetypes.Store
	CaseTypeHandler         *casetypes.Handler
	CaseTypeClient          *casetypes.Client
	CaseStore               *cases.Store
	CaseHandler             *cases.Handler
	CaseClient              *cases.Client
	WebAppHandler           *webapp.Handler
	HttpServer              *http.Server
	TeamStore               *teams.Store
	TeamHandler             *teams.Handler
	TeamClient              *teams.Client
	MembershipHandler       *memberships.Handler
	MembershipStore         *memberships.Store
	membershipsClient       *memberships.Client
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
}

func (o Options) Complete(ctx context.Context) (CompletedOptions, error) {

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

	completedOptions := CompletedOptions{
		MongoClient:       mongoClient,
		TemplateDirectory: o.TemplateDirectory,
		Address:           o.Address,
		MongoDatabase:     o.MongoDatabase,
	}
	return completedOptions, nil
}

type CompletedOptions struct {
	MongoClient       *mongo.Client
	TemplateDirectory string
	Address           string
	MongoDatabase     string
}

func (c CompletedOptions) New(ctx context.Context) *Server {

	router := mux.NewRouter()

	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			start := time.Now()

			stWriter := &statusWriter{w: writer}
			handler.ServeHTTP(stWriter, request)

			end := time.Now()

			statusCode := stWriter.statusCode
			if stWriter.statusCode == 0 {
				statusCode = 200
			}

			fields := logrus.Fields{
				"method":     request.Method,
				"statusCode": statusCode,
				"path":       request.URL.Path,
				"responseMs": math.Round(float64(end.Sub(start).Nanoseconds())/1000000.0*100.0) / 100.0,
			}

			if stWriter.statusCode < 400 {
				logrus.WithFields(fields).Infof("")
			} else {
				logrus.WithFields(fields).
					WithError(fmt.Errorf("inbound request failed with status code: %d", statusCode)).
					Errorf("")
			}

		})
	})

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

	// Beneficiaries
	beneficiariesStore := beneficiaries.NewStore(c.MongoClient, c.MongoDatabase)
	beneficiaryHandler := beneficiaries.NewHandler(beneficiariesStore)
	beneficiaryClient := beneficiaries.NewClient(c.Address)
	if err := beneficiaries.SeedDatabase(ctx, beneficiariesStore); err != nil {
		panic(err)
	}
	router.Path("/apis/v1/beneficiaries").Methods("GET").HandlerFunc(beneficiaryHandler.List)
	router.Path("/apis/v1/beneficiaries/{id}").Methods("GET").HandlerFunc(beneficiaryHandler.Get)
	router.Path("/apis/v1/beneficiaries/{id}").Methods("PUT").HandlerFunc(beneficiaryHandler.Update)
	router.Path("/apis/v1/beneficiaries").Methods("POST").HandlerFunc(beneficiaryHandler.Create)

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

	// Teams
	teamStore := teams.NewStore(partyStore)
	if err := teams.Init(ctx, teamStore, partyTypeStore); err != nil {
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

	// WebApp
	webAppOptions := webapp.Options{
		TemplateDirectory: c.TemplateDirectory,
	}
	webAppHandler, err := webapp.NewHandler(webAppOptions,
		attributeClient,
		beneficiaryClient,
		relationshipTypeClient,
		relationshipClient,
		partyClient,
		partyTypeClient,
		caseTypeClient,
		caseClient,
		teamClient,
		membershipClient,
	)
	if err != nil {
		panic(err)
	}

	router.Path("/beneficiaries").HandlerFunc(webAppHandler.Beneficiaries)
	router.Path("/beneficiaries/{id}").HandlerFunc(webAppHandler.Beneficiary)
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

	// Seed database for development
	if err := beneficiaries.SeedDatabase(ctx, beneficiariesStore); err != nil {
		panic(err)
	}

	httpServer := &http.Server{
		Addr:    ":9000",
		Handler: router,
	}

	srv := &Server{
		MongoClient:             c.MongoClient,
		AttributeStore:          attributeStore,
		AttributeHandler:        attributeHandler,
		AttributeClient:         attributeClient,
		BeneficiaryStore:        beneficiariesStore,
		BeneficiaryHandler:      beneficiaryHandler,
		BeneficiaryClient:       beneficiaryClient,
		RelationshipTypeStore:   relationshipTypeStore,
		RelationshipTypeHandler: relationshipTypeHandler,
		RelationshipTypeClient:  relationshipTypeClient,
		RelationshipStore:       relationshipStore,
		RelationshipHandler:     relationshipHandler,
		RelationshipClient:      relationshipClient,
		PartyStore:              partyStore,
		PartyHandler:            partyHandler,
		PartyClient:             partyClient,
		PartyTypeStore:          partyTypeStore,
		PartyTypeHandler:        partyTypeHandler,
		PartyTypeClient:         partyTypeClient,
		CaseTypeStore:           caseTypeStore,
		CaseTypeHandler:         caseTypeHandler,
		CaseTypeClient:          caseTypeClient,
		CaseStore:               caseStore,
		CaseHandler:             caseHandler,
		CaseClient:              caseClient,
		TeamStore:               teamStore,
		TeamHandler:             teamHandler,
		TeamClient:              teamClient,
		MembershipStore:         membershipStore,
		MembershipHandler:       membershipHandler,
		membershipsClient:       membershipClient,
		WebAppHandler:           webAppHandler,
		HttpServer:              httpServer,
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
