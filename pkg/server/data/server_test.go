package data

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nrc-no/core/pkg/server/data/api"
	"github.com/nrc-no/core/pkg/server/data/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	ctx    context.Context
	cancel context.CancelFunc
	cli    client.HTTPClient
}

func (s *Suite) SetupSuite() {
	server, err := NewServer(Options{})
	if err != nil {
		s.T().Fatal(err)
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	server.Start(s.ctx)

	s.cli = client.NewClient(fmt.Sprintf("http://localhost:%d", server.Port()))
}

func (s *Suite) MustCreateTable(table api.Table) api.Table {
	ret, err := s.cli.CreateTable(s.ctx, table)
	if !s.NoError(err) {
		s.T().Fatal(err)
	}
	return ret
}

func (s *Suite) MustPutRecord(record api.PutRecordRequest) api.Record {
	ret, err := s.cli.PutRecord(s.ctx, record)
	if !s.NoError(err) {
		s.T().Fatal(err)
	}
	return ret
}

func (s *Suite) MustGetRecord(request api.GetRecordRequest) api.Record {
	ret, err := s.cli.GetRecord(s.ctx, request)
	if !s.NoError(err) {
		s.T().Fatal(err)
	}
	return ret
}

func (s *Suite) MustGetChanges(request api.GetChangesRequest) api.Changes {
	ret, err := s.cli.GetChanges(s.ctx, request)
	if !s.NoError(err) {
		s.T().Fatal(err)
	}
	return ret
}

func (s *Suite) AssertHasChangesSince(since int64, count int) {
	changes := s.MustGetChanges(api.GetChangesRequest{Since: since})
	assert.Equal(s.T(), count, len(changes.Items))
	s.T().Log(changes)
}

func (s *Suite) TearDownSuite() {
	s.cancel()
}

func (s *Suite) TestServer() {

	var (
		table         api.Table
		createdRecord api.Record
		updatedRecord api.Record
		foundRecord   api.Record
		changes       api.Changes
	)

	s.AssertHasChangesSince(0, 0)

	// Create table
	table = s.MustCreateTable(api.Table{
		Name:    "bla",
		Columns: []api.Column{{Name: "bli", Type: "varchar"}},
	})

	s.T().Log(table)

	// Create record
	createdRecord = s.MustPutRecord(api.PutRecordRequest{
		Record: api.Record{
			ID:         "1",
			Table:      table.Name,
			Attributes: api.NewAttributes().WithString("bli", "bla"),
		},
	})
	assert.Equal(s.T(), 1, createdRecord.Revision.Num)

	s.AssertHasChangesSince(0, 1)
	s.T().Log(createdRecord)

	changes = s.MustGetChanges(api.GetChangesRequest{Since: 0})
	assert.Equal(s.T(), 1, len(changes.Items))
	s.T().Log(changes)

	// Get record
	foundRecord = s.MustGetRecord(api.GetRecordRequest{
		RecordID:  createdRecord.ID,
		TableName: table.Name,
	})

	s.AssertHasChangesSince(0, 1)
	s.T().Log(foundRecord)

	assert.Equal(s.T(), createdRecord, foundRecord)

	// Update the record
	foundRecord.Attributes.WithString("bli", "blub")
	updatedRecord = s.MustPutRecord(api.PutRecordRequest{
		Record: foundRecord,
	})
	assert.Equal(s.T(), 2, updatedRecord.Revision.Num)
	assert.Equal(s.T(), api.NewStringValue("blub", true), updatedRecord.Attributes["bli"])

	s.AssertHasChangesSince(0, 2)
	s.T().Log(updatedRecord)

	// Get record
	foundRecord = s.MustGetRecord(api.GetRecordRequest{
		RecordID:  createdRecord.ID,
		TableName: table.Name,
	})

	s.T().Log(foundRecord)

	assert.Equal(s.T(), updatedRecord, foundRecord)

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
