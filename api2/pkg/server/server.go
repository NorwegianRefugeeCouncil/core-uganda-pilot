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
	router.Path("/apis/v1/attributes").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		attributeHandler.List(w, req)
	})
	router.Path("/apis/v1/attributes/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		attributeHandler.Get(w, req)
	})
	router.Path("/apis/v1/attributes/{id}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		attributeHandler.Update(w, req)
	})
	router.Path("/apis/v1/attributes").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		attributeHandler.Post(w, req)
	})

	// Vulnerabilities
	vulnerabilityStore := vulnerability.NewStore(c.MongoClient, c.MongoDatabase)
	vulnerabilityHandler := vulnerability.NewHandler(vulnerabilityStore)
	vulnerabilityClient := vulnerability.NewClient(c.Address)
	router.Path("/apis/v1/vulnerabilities").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vulnerabilityHandler.ListVulnerabilities(w, req)
	})
	router.Path("/apis/v1/vulnerabilities/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vulnerabilityHandler.GetVulnerability(w, req)
	})
	router.Path("/apis/v1/vulnerabilities/{id}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vulnerabilityHandler.UpdateVulnerability(w, req)
	})
	router.Path("/apis/v1/vulnerabilities").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vulnerabilityHandler.PostVulnerability(w, req)
	})

	// Beneficiaries
	beneficiariesStore := beneficiaries.NewStore(c.MongoClient, c.MongoDatabase)
	beneficiaryHandler := beneficiaries.NewHandler(beneficiariesStore)
	beneficiaryClient := beneficiaries.NewClient(c.Address)
	if err := beneficiaries.SeedDatabase(ctx, beneficiariesStore); err != nil {
		panic(err)
	}
	router.Path("/apis/v1/beneficiaries").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		beneficiaryHandler.List(w, req)
	})
	router.Path("/apis/v1/beneficiaries/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		beneficiaryHandler.Get(w, req)
	})
	router.Path("/apis/v1/beneficiaries/{id}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		beneficiaryHandler.Update(w, req)
	})
	router.Path("/apis/v1/beneficiaries").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		beneficiaryHandler.Create(w, req)
	})

	// RelationshipTypes
	relationshipTypeStore := relationshiptypes.NewStore(c.MongoClient, c.MongoDatabase)
	if err := relationshiptypes.Init(ctx, relationshipTypeStore); err != nil {
		panic(err)
	}
	relationshipTypeHandler := relationshiptypes.NewHandler(relationshipTypeStore)
	relationshipTypeClient := relationshiptypes.NewClient(c.Address)
	router.Path("/apis/v1/relationshiptypes").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		relationshipTypeHandler.List(w, req)
	})
	router.Path("/apis/v1/relationshiptypes/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		relationshipTypeHandler.Get(w, req)
	})
	router.Path("/apis/v1/relationshiptypes/{id}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		relationshipTypeHandler.Put(w, req)
	})
	router.Path("/apis/v1/relationshiptypes").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		relationshipTypeHandler.Post(w, req)
	})

	// Relationships
	relationshipStore := relationships.NewStore(c.MongoClient, c.MongoDatabase)
	relationshipHandler := relationships.NewHandler(relationshipStore)
	relationshipClient := relationships.NewClient(c.Address)
	router.Path("/apis/v1/relationships").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		relationshipHandler.List(w, req)
	})
	router.Path("/apis/v1/relationships/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		relationshipHandler.Get(w, req)
	})
	router.Path("/apis/v1/relationships/{id}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		relationshipHandler.Put(w, req)
	})
	router.Path("/apis/v1/relationships").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		relationshipHandler.Post(w, req)
	})

	// Parties
	partyStore := parties.NewStore(c.MongoClient, c.MongoDatabase)
	if err := parties.Init(ctx, partyStore); err != nil {
		panic(err)
	}
	partyHandler := parties.NewHandler(partyStore)
	partyClient := parties.NewClient(c.Address)
	router.Path("/apis/v1/parties").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		partyHandler.List(w, req)
	})
	router.Path("/apis/v1/parties/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		partyHandler.Get(w, req)
	})
	router.Path("/apis/v1/parties/{id}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		partyHandler.Put(w, req)
	})
	router.Path("/apis/v1/parties").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		partyHandler.Post(w, req)
	})

	// Party Types
	partyTypeStore := partytypes.NewStore(c.MongoClient, c.MongoDatabase)
	if err := partytypes.Init(ctx, partyTypeStore); err != nil {
		panic(err)
	}
	partyTypeHandler := partytypes.NewHandler(partyTypeStore)
	partyTypeClient := partytypes.NewClient(c.Address)
	router.Path("/apis/v1/partytypes").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		partyTypeHandler.List(w, req)
	})
	router.Path("/apis/v1/partytypes/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		partyTypeHandler.Get(w, req)
	})
	router.Path("/apis/v1/partytypes/{id}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		partyTypeHandler.Put(w, req)
	})
	router.Path("/apis/v1/partytypes").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		partyTypeHandler.Post(w, req)
	})

	// PartyTypeSchemas
	partyTypeSchemaStore := partytypeschemas.NewStore(c.MongoClient, c.MongoDatabase)
	partyTypeSchemaHandler := partytypeschemas.NewHandler(partyTypeSchemaStore)
	// TOOD: partyTypeSchemaClient := partytypeschemas.NewClient(serverOptions.Address)
	router.Path("/apis/v1/partytypeschemas").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		partyTypeSchemaHandler.List(w, req)
	})
	router.Path("/apis/v1/partytypeschemas/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		partyTypeSchemaHandler.Get(w, req)
	})
	router.Path("/apis/v1/partytypeschemas/{id}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		partyTypeSchemaHandler.Put(w, req)
	})
	router.Path("/apis/v1/partytypeschemas").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		partyTypeSchemaHandler.Post(w, req)
	})

	// Cases
	caseStore := cases.NewStore(c.MongoClient, c.MongoDatabase)
	if err := cases.Init(ctx, caseStore); err != nil {
		panic(err)
	}
	caseHandler := cases.NewHandler(caseStore)
	caseClient := cases.NewClient(c.Address)
	router.Path("/apis/v1/cases").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		caseHandler.List(w, req)
	})
	router.Path("/apis/v1/cases/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		caseHandler.Get(w, req)
	})
	router.Path("/apis/v1/cases/{id}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		caseHandler.Put(w, req)
	})
	router.Path("/apis/v1/cases").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		caseHandler.Post(w, req)
	})
	// CaseTypes
	caseTypeStore := casetypes.NewStore(c.MongoClient, c.MongoDatabase)
	if err := casetypes.Init(ctx, caseTypeStore); err != nil {
		panic(err)
	}
	caseTypeHandler := casetypes.NewHandler(caseTypeStore)
	caseTypeClient := casetypes.NewClient(c.Address)
	router.Path("/apis/v1/casetypes").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		caseTypeHandler.List(w, req)
	})
	router.Path("/apis/v1/casetypes/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		caseTypeHandler.Get(w, req)
	})
	router.Path("/apis/v1/casetypes/{id}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		caseTypeHandler.Put(w, req)
	})
	router.Path("/apis/v1/casetypes").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		caseTypeHandler.Post(w, req)
	})

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

	router.Path("/vulnerabilities/{id}").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.Vulnerability(w, req)
	})
	router.Path("/vulnerabilities").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.Vulnerabilities(w, req)
	})
	router.Path("/beneficiaries").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.Beneficiaries(w, req)
	})
	router.Path("/beneficiaries/{id}").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.Beneficiary(w, req)
	})
	router.Path("/settings").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.Settings(w, req)
	})
	router.Path("/settings/attributes").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.Attributes(w, req)
	})
	router.Path("/settings/attributes/new").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.NewAttribute(w, req)
	})
	router.Path("/settings/attributes/{id}").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.Attribute(w, req)
	})
	router.Path("/settings/relationshiptypes").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.RelationshipTypes(w, req)
	})
	router.Path("/settings/relationshiptypes/new").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.NewRelationshipType(w, req)
	})
	router.Path("/settings/relationshiptypes/{id}").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.RelationshipType(w, req)
	})
	router.Path("/settings/partytypes").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.PartyTypes(w, req)
	})
	router.Path("/settings/partytypes/{id}").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.PartyType(w, req)
	})
	router.Path("/settings/countries").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.CountrySettings(w, req)
	})

	router.Path("/cases").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.Cases(w, req)
	})
	router.Path("/cases/new").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.NewCase(w, req)
	})
	router.Path("/cases/{id}").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.Case(w, req)
	})
	router.Path("/settings/casetypes").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.CaseTypes(w, req)
	})
	router.Path("/settings/casetypes/new").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.NewCaseType(w, req)
	})
	router.Path("/settings/casetypes/{id}").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.CaseType(w, req)
	})

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
