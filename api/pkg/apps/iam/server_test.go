// +build integration

package iam_test

import (
	"context"
	"errors"
	. "github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"net/http"
	"os"
	"testing"
)

type Suite struct {
	suite.Suite
	server     *Server
	serverOpts *server.GenericServerOptions
	ctx        context.Context
	client     *ClientSet
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

var (
	ctx = context.Background()

	mongoUsername = GetEnvOrDefault("MONGO_USERNAME", "root")
	mongoPassword = GetEnvOrDefault("MONGO_PASSWORD", "example")
	mongoHost     = GetEnvOrDefault("MONGO_HOST", "localhost:27017")
	mongoDatabase = GetEnvOrDefault("MONGO_DATABASE", "e2e")
)

func (s *Suite) SetupSuite() {
	// Using a random port
	ip := net.ParseIP("127.0.0.1")
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP: ip,
	})
	if !assert.NoError(s.T(), err) {
		s.T().Fatal(err)
		return
	}
	s.T().Logf("Listening at: %s", listener.Addr().String())
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if !assert.NoError(s.T(), err) {
		s.T().Fatal(err)
		return
	}

	mongoClient, err := mongo.NewClient(options.Client().SetAuth(options.Credential{Username: mongoUsername, Password: mongoPassword}).SetHosts([]string{mongoHost}))
	if !assert.NoError(s.T(), err) {
		s.T().Fatal(err)
		return
	}
	if err := mongoClient.Connect(ctx); !assert.NoError(s.T(), err) {
		s.T().Fatal(err)
		return
	}

	opts := &server.GenericServerOptions{
		MongoClientFn: func(ctx context.Context) (*mongo.Client, error) {
			return mongoClient, nil
		},
		MongoDatabase: mongoDatabase,
		Environment:   "Development",
	}
	s.serverOpts = opts

	srv, err := NewServer(ctx, opts)
	if !assert.NoError(s.T(), err) {
		s.T().Fatal()
	}

	s.ctx = ctx
	s.server = srv
	s.client = NewClientSet(&rest.RESTConfig{
		Scheme: "http",
		Host:   "localhost:" + port,
		Headers: map[string][]string{
			"X-Authenticated-User-Subject": {"mock-auth-user"},
		},
	})

	go func() {
		if err := http.Serve(listener, srv); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
		} else {
			s.T().Fatal(err)
		}
	}()

}

// This will run before each test in the suite but must be called manually before subtests
func (s *Suite) SetupTest() {
	err := s.server.ResetDB(ctx, mongoDatabase)
	if err != nil {
		return
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}

//
// Helpers
//
func contains(s []string, item string) bool {
	for _, a := range s {
		if a == item {
			return true
		}
	}
	return false
}

func newUUID() string {
	return uuid.NewV4().String()
}

func (s *Suite) mockPartyTypes(n int) []*PartyType {
	var partyTypes []*PartyType
	for i := 0; i < n; i++ {
		partyTypes = append(partyTypes, &PartyType{
			ID:        newUUID(),
			Name:      newUUID(),
			IsBuiltIn: false,
		})
	}
	return partyTypes
}

func (s *Suite) mockAttributes(n int) []*Attribute {
	var attributes []*Attribute
	for i := 0; i < n; i++ {
		attributes = append(attributes, &Attribute{
			ID:                           newUUID(),
			Name:                         newUUID(),
			PartyTypeIDs:                 make([]string, 0),
			IsPersonallyIdentifiableInfo: false,
			Translations:                 make([]AttributeTranslation, 0),
		})
	}
	return attributes
}

func (s *Suite) mockParties(n int) []*Party {
	var parties []*Party
	for i := 0; i < n; i++ {
		parties = append(parties, &Party{
			ID:           newUUID(),
			PartyTypeIDs: make([]string, 0),
			Attributes:   make(map[string][]string),
		})
	}
	return parties
}

func (s *Suite) mockRelationshipTypes(n int) []*RelationshipType {
	var relationshipTypes []*RelationshipType
	for i := 0; i < n; i++ {
		relationshipTypes = append(relationshipTypes, &RelationshipType{
			ID:              newUUID(),
			IsDirectional:   false,
			Name:            newUUID(),
			FirstPartyRole:  "",
			SecondPartyRole: "",
			Rules:           nil,
		})
	}
	return relationshipTypes
}

func (s *Suite) mockRelationships(n int) []*Relationship {
	var relationships []*Relationship
	for i := 0; i < n; i++ {
		relationships = append(relationships, &Relationship{
			ID:                 newUUID(),
			RelationshipTypeID: "",
			FirstPartyID:       "",
			SecondPartyID:      "",
		})
	}
	return relationships
}

func (s *Suite) mockIndividuals(n int) []*Individual {
	var individuals []*Individual
	for i := 0; i < n; i++ {
		individual := *NewIndividual(newUUID())
		individual.Attributes.Add(FirstNameAttribute.ID, "mock")
		individual.Attributes.Add(LastNameAttribute.ID, "mock")
		individuals = append(individuals, &individual)
	}
	return individuals
}

func (s *Suite) mockMemberships(n int) []*Membership {
	var memberships []*Membership
	for i := 0; i < n; i++ {
		memberships = append(memberships, &Membership{
			ID:           newUUID(),
			TeamID:       "",
			IndividualID: "",
		})
	}
	return memberships
}
