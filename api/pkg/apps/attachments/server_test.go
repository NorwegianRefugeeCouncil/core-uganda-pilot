package attachments_test

import (
	"context"
	. "github.com/nrc-no/core/pkg/apps/attachments"
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
	client *RESTAttachmentClient
}

var ctx = context.Background()

func (s *Suite) SetupSuite() {
	s.GenericServerTestSetup = server.NewGenericServerTestSetup()
	s.server = NewServerOrDie(s.Ctx, s.GenericServerOptions)
	s.client = NewClient(testutils.SetXAuthenticatedUserSubject(s.Port))
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

func (s *Suite) aMockAttachment(attachedToId string) *Attachment {
	var newUUIDForAttachment = newUUID()
	return &Attachment{
		ID:           newUUIDForAttachment,
		AttachedToID: attachedToId,
		Body:         "{\"data\":\"" + newUUIDForAttachment + "\"}",
	}
}

func (s *Suite) mockAttachments(n int, attachedToId string) []*Attachment {
	var attachments []*Attachment
	for i := 0; i < n; i++ {
		attachments = append(attachments, s.aMockAttachment(attachedToId))
	}
	return attachments
}
