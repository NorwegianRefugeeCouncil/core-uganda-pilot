package tests

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestFolderCreate() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var db types.Database
	if err := s.createDatabase(ctx, &db); !assert.NoError(s.T(), err) {
		return
	}

	var folder types.Folder
	if err := s.cli.CreateFolder(ctx, &types.Folder{
		Name:       "My Folder",
		DatabaseID: db.ID,
	}, &folder); !assert.NoError(s.T(), err) {
		return
	}
}
