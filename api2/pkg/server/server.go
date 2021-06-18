package server

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/cases/cases"
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/beneficiaries"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypeschemas"
	"github.com/nrc-no/core-kafka/pkg/parties/relationships"
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	"github.com/nrc-no/core-kafka/pkg/services/vulnerability"
	"github.com/nrc-no/core-kafka/pkg/webapp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"net/http"
	"strings"
)

type Server struct {
	MongoClient             *mongo.Client
	AttributeStore          *attributes.Store
	AttributeHandler        *attributes.Handler
	AttributeClient         *attributes.Client
	VulnerabilityStore      *vulnerability.Store
	VulnerabilityHandler    *vulnerability.Handler
	VulnerabilityClient     *vulnerability.Client
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
}

type Options struct {
	TemplateDirectory string
	Address           string
	MongoDatabase     string
	MongoUsername     string
	MongoPassword     string
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

	// Setup address
	if len(o.Address) == 0 {
		o.Address = "http://localhost:9000"
	}

	// Setup template directory
	if len(o.TemplateDirectory) == 0 {
		o.TemplateDirectory = "pkg/webapp/templates"
	}

	if len(o.MongoDatabase) == 0 {
		o.MongoDatabase = "core"
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

	// Vulnerabilities
	vulnerabilityStore := vulnerability.NewStore(c.MongoClient, c.MongoDatabase)
	vulnerabilityHandler := vulnerability.NewHandler(vulnerabilityStore)
	vulnerabilityClient := vulnerability.NewClient(c.Address)
	router.Path("/apis/v1/vulnerabilities").Methods("GET").HandlerFunc(vulnerabilityHandler.ListVulnerabilities)
	router.Path("/apis/v1/vulnerabilities/{id}").Methods("GET").HandlerFunc(vulnerabilityHandler.GetVulnerability)
	router.Path("/apis/v1/vulnerabilities/{id}").Methods("PUT").HandlerFunc(vulnerabilityHandler.UpdateVulnerability)
	router.Path("/apis/v1/vulnerabilities").Methods("POST").HandlerFunc(vulnerabilityHandler.PostVulnerability)

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
	caseHandler := cases.NewHandler(caseStore)
	caseClient := cases.NewClient(c.Address)
	router.Path("/apis/v1/cases").Methods("GET").HandlerFunc(caseHandler.List)
	router.Path("/apis/v1/cases/{id}").Methods("GET").HandlerFunc(caseHandler.Get)
	router.Path("/apis/v1/cases/{id}").Methods("PUT").HandlerFunc(caseHandler.Put)
	router.Path("/apis/v1/cases").Methods("POST").HandlerFunc(caseHandler.Post)
	// CaseTypes
	caseTypeStore := casetypes.NewStore(c.MongoClient, c.MongoDatabase)
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
	webAppHandler, err := webapp.NewHandler(webAppOptions,
		attributeClient,
		vulnerabilityClient,
		beneficiaryClient,
		relationshipTypeClient,
		relationshipClient,
		partyClient,
		partyTypeClient,
		caseTypeClient,
		caseClient,
	)
	if err != nil {
		panic(err)
	}

	router.Path("/vulnerabilities/{id}").HandlerFunc(webAppHandler.Vulnerability)
	router.Path("/vulnerabilities").HandlerFunc(webAppHandler.Vulnerabilities)
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
		VulnerabilityStore:      vulnerabilityStore,
		VulnerabilityHandler:    vulnerabilityHandler,
		VulnerabilityClient:     vulnerabilityClient,
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
