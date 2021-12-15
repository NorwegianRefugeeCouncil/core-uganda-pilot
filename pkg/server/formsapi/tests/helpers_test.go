package tests

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) createDatabase(ctx context.Context, out *types.Database) error {
	if err := s.cli.CreateDatabase(ctx, &types.Database{Name: "A Database"}, out); !assert.NoError(s.T(), err) {
		return err
	}
	return nil
}

func (s *Suite) createFolder(ctx context.Context, dbId string, out *types.Folder) error {
	if err := s.cli.CreateFolder(ctx, &types.Folder{
		Name:       "My Folder",
		DatabaseID: dbId,
	}, out); err != nil {
		return err
	}
	return nil
}
