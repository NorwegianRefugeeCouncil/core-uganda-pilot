// +build integration

package cms_test

import (
	"context"
	"errors"
	. "github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/testutils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type Suite struct {
	suite.Suite
	*server.GenericServerTestSuite
	server     *Server
	serverOpts *server.GenericServerOptions
	ctx        context.Context
	client     *ClientSet
}

var ctx = context.Background()

func (s *Suite) SetupSuite() {
	args := s.GenericSetupSuite()

	s.ctx = ctx
	s.serverOpts = &args.Options

	srv, err := NewServer(s.ctx, s.serverOpts)
	if !assert.NoError(s.T(), err) {
		s.T().Fatal()
	}

	s.server = srv
	s.client = NewClientSet(testutils.SetXAuthenticatedUserSubject(args.Port))

	go func() {
		if err := http.Serve(args.Listener, srv); err != nil {
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
	err := s.server.ResetDB(s.ctx, s.serverOpts.MongoDatabase)
	if err != nil {
		return
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

func (s *Suite) mockCaseTypes(n int) []*CaseType {
	var caseTypes []*CaseType
	for i := 0; i < n; i++ {
		caseTypes = append(caseTypes, &CaseType{
			Name: "mock",
		})
	}
	return caseTypes
}

func (s *Suite) mockCases(n int) []*Case {
	var cases []*Case
	for i := 0; i < n; i++ {
		cases = append(cases, &Case{
			TeamID:    newUUID(),
			CreatorID: "mock-auth-user",
		})
	}
	return cases
}
