package server

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/cases/cases"
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	"github.com/nrc-no/core-kafka/pkg/intake"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/beneficiaries"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypeschemas"
	"github.com/nrc-no/core-kafka/pkg/parties/relationships"
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	"github.com/nrc-no/core-kafka/pkg/services/vulnerability"
	"github.com/nrc-no/core-kafka/pkg/webapp"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	mongoClient *mongo.Client
}

func NewServer(
	ctx context.Context,
) {

	if err := CreateTopics(kafka.TopicConfig{
		Topic:             "submissions",
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, kafka.TopicConfig{
		Topic:             "beneficiaries",
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, kafka.TopicConfig{
		Topic:             "intake",
		NumPartitions:     1,
		ReplicationFactor: 1,
	}); err != nil {
		panic(err)
	}

	mongoClient, err := mongo.NewClient(options.Client().SetAuth(
		options.Credential{
			Username: "root",
			Password: "example",
		}))
	if err != nil {
		panic(err)
	}

	if err := mongoClient.Connect(ctx); err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	// Attributes
	attributeStore := attributes.NewStore(mongoClient)
	if err := attributes.Init(ctx, attributeStore); err != nil {
		panic(err)
	}
	attributeHandler := attributes.NewHandler(attributeStore)
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
	vulnerabilityStore := vulnerability.NewStore(mongoClient)
	vulnerabilityHandler := vulnerability.NewHandler(vulnerabilityStore)
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
	beneficiariesStore := beneficiaries.NewStore(mongoClient)
	beneficiaryReader, err := NewReader(ctx, "beneficiaries", nil)
	beneficiaryHandler := beneficiaries.NewHandler(beneficiariesStore)
	if err != nil {
		panic(err)
	}
	beneficiariesListener := beneficiaries.NewListener(beneficiariesStore, beneficiaryReader)
	go func() {
		beneficiariesListener.Run(ctx)
	}()
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
	relationshipTypeStore := relationshiptypes.NewStore(mongoClient)
	if err := relationshiptypes.Init(ctx, relationshipTypeStore); err != nil {
		panic(err)
	}
	relationshipTypeHandler := relationshiptypes.NewHandler(relationshipTypeStore)
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
	relationshipStore := relationships.NewStore(mongoClient)
	relationshipHandler := relationships.NewHandler(relationshipStore)
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
	partyStore := parties.NewStore(mongoClient)
	partyHandler := parties.NewHandler(partyStore)
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
	partyTypeStore := partytypes.NewStore(mongoClient)
	if err := partytypes.Init(ctx, partyTypeStore); err != nil {
		panic(err)
	}
	partyTypeHandler := partytypes.NewHandler(partyTypeStore)
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
	partyTypeSchemaStore := partytypeschemas.NewStore(mongoClient)
	partyTypeSchemaHandler := partytypeschemas.NewHandler(partyTypeSchemaStore)
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

	// Intake
	intakeStore := intake.NewStore(mongoClient)
	intakeWriter, err := NewWriter("intake")
	if err != nil {
		panic(err)
	}
	intakeHandler := intake.NewHandler(intakeStore, intakeWriter)
	router.Path("/apis/v1/intake").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		intakeHandler.PostSubmission(w, req)
	})

	// Cases
	caseStore := cases.NewStore(mongoClient)
	caseHandler := cases.NewHandler(caseStore)
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
	caseTypeStore := casetypes.NewStore(mongoClient)
	castTypeHandler := casetypes.NewHandler(caseTypeStore)
	router.Path("/apis/v1/casetypes").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		castTypeHandler.List(w, req)
	})
	router.Path("/apis/v1/casetypes/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		castTypeHandler.Get(w, req)
	})
	router.Path("/apis/v1/casetypes/{id}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		castTypeHandler.Put(w, req)
	})
	router.Path("/apis/v1/casetypes").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		castTypeHandler.Post(w, req)
	})

	// WebApp
	webAppHandler, err := webapp.NewHandler()
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
	router.Path("/cases/casetypes/new").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.NewCaseType(w, req)
	})
	router.Path("/cases/casetypes/{id}").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		webAppHandler.CaseType(w, req)
	})

	http.ListenAndServe(":9000", router)
}

func NewReader(ctx context.Context, topic string, offset *int64) (*kafka.Reader, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic,
	})
	if offset != nil {
		if err := reader.SetOffset(*offset); err != nil {
			return nil, err
		}
	} else {
		if err := reader.SetOffsetAt(ctx, time.Now()); err != nil {
			return nil, err
		}
	}
	return reader, nil
}

func NewWriter(topic string) (*kafka.Writer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Balancer: &kafka.LeastBytes{},
		Topic:    topic,
	}
	return writer, nil
}

func CreateTopics(topics ...kafka.TopicConfig) error {

	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()
	if err := controllerConn.CreateTopics(topics...); err != nil {
		return err
	}

	return nil
}
