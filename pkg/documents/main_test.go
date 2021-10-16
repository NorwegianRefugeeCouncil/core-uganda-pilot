package documents

import (
	"context"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	server         *Server
	mongoFn        func() (*mongo.Client, error)
	mongoCli       *mongo.Client
	databaseName   string
	collectionName string
	timeTeller     utils.TimeTeller
	uidGenerator   utils.UIDGenerator
	client         Interface
	done           chan struct{}
}

func (s *Suite) SetupSuite() {

	s.databaseName = "test"
	s.collectionName = "documents"

	var err error

	s.mongoCli, err = mongo.Connect(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}

	s.mongoFn = func() (*mongo.Client, error) {
		return s.mongoCli, nil
	}

	s.timeTeller = utils.NewMockTimeTeller(time.Now())
	s.uidGenerator = utils.NewUIDGenerator()

	s.server = NewServer(s.mongoFn, s.databaseName, s.collectionName, s.timeTeller, s.uidGenerator)

	s.done = make(chan struct{}, 1)

	if err := s.server.Start(s.done); err != nil {
		s.T().Fatal(err)
	}

	s.client = s.server.NewClient()

	_, err = s.mongoCli.Database(s.databaseName).Collection(s.collectionName).DeleteMany(context.Background(), bson.M{})
	if err != nil {
		s.T().Fatal(err)
	}

}

func (s *Suite) TearDownSuite() {
	s.done <- struct{}{}
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
