package engine

import (
	"context"
	"testing"

	"github.com/nrc-no/core/pkg/server/data/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ReconcilerSuite struct {
	suite.Suite
	ctx         context.Context
	cancel      context.CancelFunc
	source      *Bench
	destination *Bench
}

func (s *ReconcilerSuite) SetupSuite() {
	s.source = NewTestBench("file::memory:")
	s.destination = NewTestBench("file::memory:")
}

func (s *ReconcilerSuite) SetupTest() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	if err := s.source.Reset(); err != nil {
		panic(err)
	}
	if err := s.destination.Reset(); err != nil {
		panic(err)
	}
}

func (s *ReconcilerSuite) TearDownTest() {
	s.cancel()
}

func (s *ReconcilerSuite) TestReconcileEmptySourceAndDestination() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := &reconciler{}
	if err := r.Reconcile(ctx, s.source.Engine, s.destination.Engine); !assert.NoError(s.T(), err) {
		return
	}
}

func (s *ReconcilerSuite) TestReconcileRecordFromSource() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := &reconciler{}

	for _, e := range []api.Engine{s.source.Engine, s.destination.Engine} {
		if err := e.CreateTable(ctx, api.Table{
			Name: "test",
			Columns: []api.Column{
				{
					Name: "field1",
					Type: "varchar",
				},
			},
		}); !assert.NoError(s.T(), err) {
			return
		}
	}

	srcRec, err := s.source.Engine.PutRecord(ctx, api.PutRecordRequest{
		Record: api.Record{
			Table: "test",
			ID:    "mock_id_1",
			Attributes: map[string]api.Value{
				"field1": api.NewStringValue("mock_value_1", true),
			},
		},
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	if err := r.Reconcile(ctx, s.source.Engine, s.destination.Engine); !assert.NoError(s.T(), err) {
		return
	}

	// find record in destination
	destRec, err := s.source.Engine.GetRecord(ctx, api.GetRecordRequest{
		TableName: "test",
		RecordID:  srcRec.GetID(),
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), srcRec, destRec)

	s.T().Log(srcRec, destRec)
}

func TestReconcilerSuite(t *testing.T) {
	suite.Run(t, new(ReconcilerSuite))
}
