package store

import (
	"context"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/stretchr/testify/assert"
)

// TestFolderCreate tests that we can create a simple folder in a database
func (s *Suite) TestFolderCreate() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := s.mustCreateDatabase(ctx)
	folder, err := s.createFolder(ctx, db.ID)
	actual, err := s.folderStore.Get(ctx, folder.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), actual, folder)
}

// TestFolderCreateWithDuplicateID tests that it is not possible to create 2 folders with the same ID
func (s *Suite) TestFolderCreateWithDuplicateID() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := s.mustCreateDatabase(ctx)
	folder := s.mustCreateFolder(ctx, db.ID)
	_, err := s.folderStore.Create(ctx, folder)
	if !assert.Error(s.T(), err) {
		return
	}
	assert.Equal(s.T(), meta.StatusReasonAlreadyExists, meta.ReasonForError(err))
}

// TestFolderCreateOnNotExistingDatabase tests that it is not possible to create a folder on a database that does not exist
func (s *Suite) TestFolderCreateOnNotExistingDatabase() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := s.createFolder(ctx, "non-existing")
	if !assert.Error(s.T(), err) {
		return
	}
	assert.Equal(s.T(), meta.StatusReasonInvalid, meta.ReasonForError(err))
	assert.True(s.T(), meta.HasCauseForError(err, "databaseId", meta.CauseTypeFieldValueNotFound))
}

// TestFolderCreateWithParent tests that it is possible to create a folder with a parent folder
func (s *Suite) TestFolderCreateWithParent() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := s.mustCreateDatabase(ctx)
	parent := s.mustCreateFolder(ctx, db.ID)
	folder := s.mustCreateFolderWithParent(ctx, db.ID, parent.ID)
	assert.Equal(s.T(), parent.ID, folder.ParentID)
}

// TestFolderCreateWithNonExistingParent tests that it is not possible to create a folder with a parent that does not exist
func (s *Suite) TestFolderCreateWithNonExistingParent() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := s.mustCreateDatabase(ctx)
	_, err := s.createFolderWithParent(ctx, db.ID, "non-existing")
	if !assert.Error(s.T(), err) {
		return
	}
	assert.Equal(s.T(), meta.StatusReasonInvalid, meta.ReasonForError(err))
	assert.True(s.T(), meta.HasCauseForError(err, "parentId", meta.CauseTypeFieldValueNotFound))
}

// TestFolderDelete tests that it is possible to delete a folder
func (s *Suite) TestFolderDelete() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := s.mustCreateDatabase(ctx)
	folder := s.mustCreateFolder(ctx, db.ID)
	if err := s.folderStore.Delete(ctx, folder.ID); !assert.NoError(s.T(), err) {
		return
	}
}

// TestFolderParentDelete tests that it is possible to delete a parent folder
func (s *Suite) TestFolderParentDelete() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := s.mustCreateDatabase(ctx)
	parent := s.mustCreateFolder(ctx, db.ID)
	child := s.mustCreateFolderWithParent(ctx, db.ID, parent.ID)
	if err := s.folderStore.Delete(ctx, parent.ID); !assert.NoError(s.T(), err) {
		return
	}
	_, getChildErr := s.folderStore.Get(ctx, child.ID)
	if !assert.Error(s.T(), getChildErr) {
		return
	}
	assert.Equal(s.T(), meta.StatusReasonNotFound, meta.ReasonForError(getChildErr))
	_, getParentErr := s.folderStore.Get(ctx, parent.ID)
	if !assert.Error(s.T(), getParentErr) {
		return
	}
	assert.Equal(s.T(), meta.StatusReasonNotFound, meta.ReasonForError(getParentErr))
}

// TestFolderDeleteNonExisting tests that it is not possible to delete a folder that does not exist
func (s *Suite) TestFolderDeleteNonExisting() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := s.folderStore.Delete(ctx, "non-existing")
	if !assert.Error(s.T(), err) {
		return
	}
	assert.Equal(s.T(), meta.StatusReasonNotFound, meta.ReasonForError(err))
}

// TestListFolders tests that it is possible to list folders
func (s *Suite) TestListFolders() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := s.mustCreateDatabase(ctx)
	folder1 := s.mustCreateFolder(ctx, db.ID)
	folder2 := s.mustCreateFolder(ctx, db.ID)
	list, err := s.folderStore.List(ctx)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, folder1)
	assert.Contains(s.T(), list.Items, folder2)
}
