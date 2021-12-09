package tests

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestDatabaseCreateGetList() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var out types.Database
	if err := s.cli.CreateDatabase(ctx, &types.Database{Name: "My Database"}, &out); !assert.NoError(s.T(), err) {
		return
	}

	var got types.Database
	if err := s.cli.GetDatabase(ctx, out.ID, &got); !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), out, got)

	var list types.DatabaseList
	if err := s.cli.ListDatabases(ctx, &list); !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, &got)

	if err := s.cli.DeleteDatabase(ctx, out.ID); !assert.NoError(s.T(), err) {
		return
	}

}
