package iam_test

import (
	"context"
	. "github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/generic/server"
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

var ctx = context.Background()

func (s *Suite) SetupSuite() {
	s.GenericServerTestSetup = server.NewGenericServerTestSetup()
	s.server = NewServerOrDie(s.Ctx, s.GenericServerOptions)
	s.client = NewClientSet(testutils.SetXAuthenticatedUserSubject(s.Port))
	s.Serve(s.T(), s.server)
}

// This will run before each test in the suite but must be called manually before subtests
func (s *Suite) SetupTest() {
	err := s.server.ResetDB(ctx, s.GenericServerOptions.MongoDatabase)
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

//func (s *Suite) mockPartyTypes(n int) []*PartyType {
//	var partyTypes []*PartyType
//	for i := 0; i < n; i++ {
//		partyTypes = append(partyTypes, &PartyType{
//			Name:      newUUID(),
//			IsBuiltIn: false,
//		})
//	}
//	return partyTypes
//}

func (s *Suite) mockAttributes(n int) []*Attribute {
	var attributes []*Attribute
	for i := 0; i < n; i++ {
		attributes = append(attributes, &Attribute{
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
			PartyTypeIDs: make([]string, 0),
			Attributes:   make(map[string][]string),
		})
	}
	return parties
}

//func (s *Suite) mockRelationshipTypes(n int) []*RelationshipType {
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

func (s *Suite) mockRelationships(n int) []*Relationship {
	var relationships []*Relationship
	for i := 0; i < n; i++ {
		relationships = append(relationships, &Relationship{
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
			TeamID:       "",
			IndividualID: "",
		})
	}
	return memberships
}

func (s *Suite) mockNationalities(n int) []*Nationality {
	var nationalities []*Nationality
	for i := 0; i < n; i++ {
		nationalities = append(nationalities, &Nationality{
			CountryID: "",
			TeamID:    "",
		})
	}
	return nationalities
}
