package tests

import (
	"testing"

	coreDBModel "github.com/nrc-no/core/pkg/server/core-db/models"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	DBFactory store.Factory
}

func (s *Suite) SetupSuite() {
	dbFactory := testutils.SetupPostgres(s.T())

	db, err := dbFactory.Get()

	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	if err := coreDBModel.Migrate(db); !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	s.DBFactory = dbFactory
}

func (s *Suite) BeforeTest(suiteName, testName string) {
	db, err := s.DBFactory.Get()

	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	if err := db.Exec("TRUNCATE TABLE entity_definition CASCADE;").Error; !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	if err := db.Exec("TRUNCATE TABLE entity_attribute CASCADE;").Error; !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	if err := db.Exec("TRUNCATE TABLE entity_relationship CASCADE;").Error; !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
