package test

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestDatabaseApi() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli := s.PublicClient(ctx)

	created := &types.Database{}
	if err := cli.CreateDatabase(ctx, &types.Database{
		Name: "my-db",
	}, created); !assert.NoError(s.T(), err) {
		return
	}
	defer cli.DeleteDatabase(ctx, created.ID)

	got := &types.Database{}
	if err := cli.GetDatabase(ctx, created.ID, got); !assert.NoError(s.T(), err) {
		return
	}

	var dbs types.DatabaseList
	if err := cli.ListDatabases(context.Background(), &dbs); !assert.NoError(s.T(), err) {
		return
	}

	found := false
	for _, item := range dbs.Items {
		if item.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		assert.Fail(s.T(), "database not found")
	}

}
