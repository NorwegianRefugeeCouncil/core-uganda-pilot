package iam_test

import (
	"context"
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/generic/server"
	. "github.com/nrc-no/core/pkg/iam"
	"github.com/nrc-no/core/pkg/testutils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	*server.GenericServerTestSetup
	suite.Suite
	server *Server
	client *ClientSet
}

func (s *Suite) SetupSuite() {
	s.GenericServerTestSetup = server.NewGenericServerTestSetup(context.Background())
	s.server = NewServerOrDie(s.Ctx, s.GenericServerOptions)
	s.client = NewClientSet(testutils.SetXAuthenticatedUserSubject(s.Port))
	s.Serve(s.T(), s.server)
}

// This will run before each test in the suite but must be called manually before subtests
func (s *Suite) SetupTest() {
	err := s.server.ResetDB(s.Ctx, s.GenericServerOptions.MongoDatabase)
	if err != nil {
		return
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}

func (s *Suite) TearDownSuite() {
	s.SetupTest()
}

//
// Helpers
//

func newUUID() string {
	return uuid.NewV4().String()
}

func mockPartyAttributeDefinition() *PartyAttributeDefinition {
	return &PartyAttributeDefinition{
		CountryID:    newUUID(),
		PartyTypeIDs: []string{newUUID()},
		FormControl:  form.Control{Name: "mock"},
	}
}

func mockPartyAttributeDefinitions(n int) []*PartyAttributeDefinition {
	var attributes []*PartyAttributeDefinition
	for i := 0; i < n; i++ {
		attributes = append(attributes, mockPartyAttributeDefinition())
	}
	return attributes
}

func mockParty() *Party {
	return &Party{
		PartyTypeIDs: make([]string, 0),
		Attributes:   make(map[string][]string),
	}
}

func mockParties(n int) []*Party {
	var parties []*Party
	for i := 0; i < n; i++ {
		parties = append(parties, mockParty())
	}
	return parties
}

// func (s *Suite) mockRelationshipTypes(n int) []*RelationshipType {
//	var relationshipTypes []*RelationshipType
//	for i := 0; i < n; i++ {
//		relationshipTypes = append(relationshipTypes, &RelationshipType{
//			IsDirectional:   false,
//			Name:            newUUID(),
//			FirstPartyRole:  "",
//			SecondPartyRole: "",
//			Rules:           nil,
//		})
//	}
//	return relationshipTypes
//}

func mockRelationship() *Relationship {
	return &Relationship{}
}

func mockRelationships(n int) []*Relationship {
	var relationships []*Relationship
	for i := 0; i < n; i++ {
		relationships = append(relationships, mockRelationship())
	}
	return relationships
}

func mockIndividual() *Individual {
	individual := NewIndividual(newUUID())
	individual.Attributes.Add(FullNameAttribute.ID, "mock")
	individual.Attributes.Add(DisplayNameAttribute.ID, "mock")
	return individual
}

func mockIndividuals(n int) []*Individual {
	var individuals []*Individual
	for i := 0; i < n; i++ {
		individuals = append(individuals, mockIndividual())
	}
	return individuals
}

func mockMembership() *Membership {
	return &Membership{}
}

func mockNationality() *Nationality {
	return &Nationality{}
}
func mockNationalities(n int) []*Nationality {
	var nationalities []*Nationality
	for i := 0; i < n; i++ {
		nationalities = append(nationalities, mockNationality())
	}
	return nationalities
}
