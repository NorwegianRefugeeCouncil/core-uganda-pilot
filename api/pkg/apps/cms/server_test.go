package cms_test

import (
	. "github.com/nrc-no/core/pkg/apps/cms"
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

func (s *Suite) SetupSuite() {
	s.GenericServerTestSetup = server.NewGenericServerTestSetup()
	s.server = NewServerOrDie(s.Ctx, s.GenericServerOptions)
	s.client = NewClientSet(testutils.SetXAuthenticatedUserSubject(s.Port))
	s.Serve(s.T(), s.server)
}

// This will run before each test in the suite but must be called manually before subtests
func (s *Suite) SetupTest() {
	err := s.server.ResetDB(s.Ctx, s.GenericServerTestSetup.GenericServerOptions.MongoDatabase)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *Suite) TearDownSuite() {
	s.SetupTest()
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

func uuidSlice(n int) []string {
	s := []string{}
	for i := 0; i < n; i++ {
		s = append(s, newUUID())
	}
	return s
}

func (s *Suite) mockCaseTypes(n int) []*CaseType {
	mockTemplate := &CaseTemplate{
		FormElements: []FormElement{{
			Type: "textarea",
			Attributes: FormElementAttribute{
				Label: "mock",
				ID:    "mock",
			},
		}},
	}
	var caseTypes []*CaseType
	for i := 0; i < n; i++ {
		caseTypes = append(caseTypes, &CaseType{
			Name:         "mock",
			PartyTypeID:  newUUID(),
			TeamID:       newUUID(),
			CaseTemplate: mockTemplate,
		})
	}
	return caseTypes
}

func aMockCase() *Case {
	return &Case{
		TeamID:  newUUID(),
		PartyID: newUUID(),
		Form: &CaseTemplate{FormElements: []FormElement{{
			Type: "textarea",
			Attributes: FormElementAttribute{
				Label: "mock",
				ID:    "mock",
			},
			Validation: FormElementValidation{},
		}}},
	}
}

func mockCases(n int) []*Case {
	var cases []*Case
	for i := 0; i < n; i++ {
		cases = append(cases, aMockCase())
	}
	return cases
}
