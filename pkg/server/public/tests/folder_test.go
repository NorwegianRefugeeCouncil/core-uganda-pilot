package tests

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestFolderCreateGetList() {
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

	var got types.Folder
	if err := s.cli.GetFolder(ctx, folder.ID, &got); !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), folder, got)

	var list types.FolderList
	if err := s.cli.ListFolders(ctx, &list); !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, &got)
}
