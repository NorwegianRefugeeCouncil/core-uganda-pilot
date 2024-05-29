package engine

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nrc-no/core/pkg/server/data/api"
	"github.com/nrc-no/core/pkg/server/data/test"
	"github.com/nrc-no/core/pkg/server/data/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	Bench *Bench
}

func (s *Suite) SetupSuite() {
	s.Bench = NewTestBench("file::memory:?cache=shared")
}

func (s *Suite) TearDownSuite() {
	if err := s.Bench.TearDown(); err != nil {
		panic(err)
	}
}

func (s *Suite) SetupTest() {
	if err := s.Bench.Reset(); err != nil {
		s.T().Fatal(err)
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestEngineCreateTable() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type testCase struct {
		// name of the test case
		name string
		// description of the test case
		description string
		// table to create
		table api.Table
		// error assertion function
		expectError assert.ErrorAssertionFunc
		// expected sql statements
		expectedStatements []test.ExpectedStatement
		// optional setup function
		doBefore func() error
	}

	testCases := []testCase{
		{
			name:        "table with no name",
			description: "Should not allow to create a table with an empty name",
			table: api.Table{
				Name: "",
				Columns: []api.Column{
					{Name: "field1", Type: "varchar"},
				},
			},
			expectError: test.ErrorIs(api.ErrInvalidTableName),
		},
		{
			name:        "table with no columns",
			description: "Should not allow to create a table with no columns",
			table: api.Table{
				Name:    "mock_table_1",
				Columns: []api.Column{},
			},
			expectError: test.ErrorIs(api.ErrEmptyColumns),
		},
		{
			name:        "table with duplicate column name",
			description: "Should not allow to create a table with duplicate column names",
			expectError: test.ErrorIs(api.ErrDuplicateColumnName),
			table: api.Table{
				Name: "mock_table_2",
				Columns: []api.Column{
					{Name: "field1", Type: "varchar"},
					{Name: "field1", Type: "varchar"},
				},
			},
		},
		{
			name:        "table with invalid column name",
			description: "Should not allow to create a table with an invalid column name",
			expectError: test.ErrorIs(api.ErrInvalidColumnName),
			table: api.Table{
				Name: "mock_table_3",
				Columns: []api.Column{
					{Name: "field1", Type: "varchar"},
					{Name: " field2 ", Type: "varchar"},
				},
			},
		},
		{
			name:        "table with invalid column type",
			description: "Should not allow to create a table with an invalid column type",
			expectError: test.ErrorIs(api.ErrInvalidColumnType),
			table: api.Table{
				Name:    "mock_table_4",
				Columns: []api.Column{{Name: "field1", Type: "Bla"}},
			},
		},
		{
			name:        "already exists",
			description: "Should not allow to create a table that already exists",
			expectError: test.ErrorIs(api.ErrTableAlreadyExists),
			doBefore: func() error {
				_, err := s.Bench.DB.Exec(`CREATE TABLE "mock_table_6" ("bla" varchar)`)
				return err
			},
			table: api.Table{
				Name: "mock_table_6",
				Columns: []api.Column{
					{Name: "field1", Type: "varchar"},
				},
			},
			expectedStatements: []test.ExpectedStatement{
				{
					SQL:    `SELECT "name" FROM "sqlite_master" WHERE "type" = 'table' AND "name" = ?`,
					Params: []interface{}{"mock_table_6"},
				},
			},
		},
		{
			name:        "valid table",
			description: "Should allow to create a valid table",
			table: api.Table{
				Name: "test_create_valid_table",
				Columns: []api.Column{
					{Name: "field1", Type: "varchar"},
					{Name: "field2", Type: "varchar"},
				},
			},
			expectedStatements: []test.ExpectedStatement{
				{
					SQL:    `SELECT "name" FROM "sqlite_master" WHERE "type" = 'table' AND "name" = ?`,
					Params: []interface{}{"test_create_valid_table"},
				}, {
					SQL:    `CREATE TABLE IF NOT EXISTS "test_create_valid_table" ("_id" varchar, "_rev" varchar NOT NULL, "field1" varchar, "field2" varchar, PRIMARY KEY ("_id"))`,
					Params: []interface{}{},
				}, {
					SQL:    `CREATE TABLE IF NOT EXISTS "test_create_valid_table_history" ("_id" varchar, "_prev" varchar NOT NULL, "_rev" varchar NOT NULL, "_deleted" boolean DEFAULT false NOT NULL, "field1" varchar, "field2" varchar, PRIMARY KEY ("_id", "_rev"))`,
					Params: []interface{}{},
				},
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.doBefore != nil {
				err := tc.doBefore()
				if !assert.NoError(t, err) {
					return
				}
			}
			s.Bench.Recorder.Reset()
			_, err := s.Bench.Engine.CreateTable(ctx, tc.table)
			if tc.expectError != nil {
				tc.expectError(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tc.expectedStatements != nil {
				s.Bench.Recorder.AssertStatementsExecuted(t, tc.expectedStatements)
			}
		})
	}
}

func (s *Suite) TestEngineGetRecordExists() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if _, err := s.Bench.Engine.CreateTable(ctx, api.Table{
		Name: "test_get_record_exists",
		Columns: []api.Column{
			{Name: "field1", Type: "varchar"},
		},
	}); !assert.NoError(s.T(), err) {
		return
	}

	_, err := s.Bench.Engine.PutRecord(ctx, api.PutRecordRequest{
		Record: api.Record{
			Table: "test_get_record_exists",
			ID:    "mock_id_1",
			Attributes: map[string]api.Value{
				"field1": api.NewStringValue("value1", true),
			},
		},
	})
	assert.NoError(s.T(), err)

	found, err := s.Bench.Engine.GetRecord(ctx, api.GetRecordRequest{
		TableName: "test_get_record_exists",
		RecordID:  "mock_id_1",
	})
	assert.NoError(s.T(), err)

	s.T().Log(found)
}

func (s *Suite) TestEngineGetRecordDoesNotExist() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := s.Bench.Engine.CreateTable(ctx, api.Table{
		Name: "mock_table_2",
		Columns: []api.Column{
			{Name: "field1", Type: "varchar"},
		},
	})
	assert.NoError(s.T(), err)

	_, err = s.Bench.Engine.GetRecord(ctx, api.GetRecordRequest{
		TableName: "mock_table_2",
		RecordID:  "bla",
	})
	assert.ErrorIs(s.T(), err, api.ErrRecordNotFound)
}

func (s *Suite) TestEnginePutRecordMultipleRevisions() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resetClock := s.Bench.Clock.UseClock(&utils.Clock{})
	defer resetClock()

	_, err := s.Bench.Engine.CreateTable(ctx, api.Table{
		Name: "mock_table_3",
		Columns: []api.Column{
			{Name: "field1", Type: "varchar"},
			{Name: "field2", Type: "varchar"},
		},
	})

	request := api.PutRecordRequest{
		Record: api.Record{
			Table: "mock_table_3",
			ID:    "mock_id_1",
			Attributes: map[string]api.Value{
				"field1": api.NewStringValue("value1", true),
				"field2": api.NewStringValue("value2", true),
			},
		},
	}

	rec, err := s.Bench.Engine.PutRecord(ctx, request)
	assert.NoError(s.T(), err)

	s.T().Log(rec)

	var revisions []api.Revision
	revisions = append(revisions, rec.GetRevision())

	for i := 0; i < 5; i++ {
		record := request.Record
		record = record.SetFieldValue("field1", api.NewStringValue(fmt.Sprintf("value1_%d", i), true))
		rec2, err := s.Bench.Engine.PutRecord(ctx, api.PutRecordRequest{Record: record})
		assert.NoError(s.T(), err)
		revisions = append(revisions, rec2.GetRevision())
		rec = rec2
	}

	s.T().Log(revisions)
}
