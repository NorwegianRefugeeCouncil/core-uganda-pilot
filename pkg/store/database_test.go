package store

import (
	"context"
	"errors"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestDatabaseCreate tests that we can create a database successfully
func (s *Suite) TestDatabaseCreate() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := s.mustCreateDatabase(ctx)
	got, err := s.dbStore.Get(ctx, db.ID)
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), db, got)
}

// TestDatabaseCreateShouldFailIfCreateSchemaFails tests that when creating a database,
// the transaction will be rolled out in case that we fail to create the actual
// SQL Schema for that database.
func (s *Suite) TestDatabaseCreateShouldFailIfCreateSchemaFails() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// mock failure to create the SQL Schema
	s.dbStore.createDatabaseSchema = func(db *gorm.DB, database *types.Database) error {
		return errors.New("mock error")
	}

	var in = &types.Database{ID: uuid.NewV4().String(), Name: "my-database"}
	if _, err := s.dbStore.Create(ctx, in); !assert.Error(s.T(), err) {
		return
	}

	// ensure the database was not stored somehow (cancelling transaction error)
	if db, err := s.dbStore.Get(ctx, in.ID); !assert.Error(s.T(), err) {
		s.T().Logf("%v", db)
		return
	}

}

// TestDatabaseCreateShouldFailIfNonUniqueDatabaseID tests that it's not possible
// to create two databases with the same ID, and that we get a
// meta.StatusReasonAlreadyExists error
func (s *Suite) TestDatabaseCreateShouldFailIfNonUniqueDatabaseID() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := s.mustCreateDatabase(ctx)

	// try recreating the same database
	_, err := s.dbStore.Create(ctx, db)
	if !assert.Error(s.T(), err) {
		return
	}

	assert.Equal(s.T(), meta.StatusReasonAlreadyExists, meta.ReasonForError(err))
}

// TestDatabaseGet tests that we can get a database after creating it
func (s *Suite) TestDatabaseGet() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := s.mustCreateDatabase(ctx)

	out, err := s.dbStore.Get(ctx, db.ID)
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), db, out)
}

// TestDatabaseGetShouldFailIfDatabaseDoesNotExist tests that we cannot get
// a database if it does not exist, and that we get a
// meta.StatusReasonNotFound
func (s *Suite) TestDatabaseGetShouldFailIfDatabaseDoesNotExist() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := s.dbStore.Get(ctx, uuid.NewV4().String())
	if !assert.Error(s.T(), err) {
		return
	}

	assert.Equal(s.T(), meta.StatusReasonNotFound, meta.ReasonForError(err))
}

// TestDatabaseListShouldReturnEmptyArrayIfNoDatabases tests that listing databases
// returns an empty list if there are no databases
func (s *Suite) TestDatabaseListShouldReturnEmptyArrayIfNoDatabases() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// clearing the databases beforehand
	if err := s.db.Where("id = id").Delete(&Database{}).Error; !assert.NoError(s.T(), err) {
		return
	}

	dbList, err := s.dbStore.List(ctx)
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), &types.DatabaseList{Items: []*types.Database{}}, dbList)
}

// TestDatabaseListShouldReturnDatabases tests that we can list the databases
func (s *Suite) TestDatabaseListShouldReturnDatabases() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db1 := s.mustCreateDatabase(ctx)
	db2 := s.mustCreateDatabase(ctx)

	dbList, err := s.dbStore.List(ctx)
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Contains(s.T(), dbList.Items, db1)
	assert.Contains(s.T(), dbList.Items, db2)
}

// TestDatabaseDelete tests that it's possible to delete a database, and that
// trying to GET the database should result in
// meta.StatusReasonNotFound
func (s *Suite) TestDatabaseDelete() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := s.mustCreateDatabase(ctx)

	if err := s.dbStore.Delete(ctx, db.ID); !assert.NoError(s.T(), err) {
		return
	}

	// test that we cannot get the database after deleting
	_, err := s.dbStore.Get(ctx, db.ID)
	if !assert.Error(s.T(), err) {
		return
	}

	assert.Equal(s.T(), meta.StatusReasonNotFound, meta.ReasonForError(err))
}

// TestDatabaseDeleteFailsIfSchemaFailsToDelete tests that when deleting a database,
// the transaction should roll back if deleting the SQL Schema fails to delete,
// and that we can still get the database afterwards
func (s *Suite) TestDatabaseDeleteFailsIfSchemaFailsToDelete() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := s.mustCreateDatabase(ctx)

	// mock that deleting the SQL Schema fails
	s.dbStore.deleteDatabaseSchema = func(db *gorm.DB, databaseId string) error {
		return errors.New("mock failure")
	}

	if err := s.dbStore.Delete(ctx, db.ID); !assert.Error(s.T(), err) {
		return
	}

	// test that we can still get the database after deleting (ensure no transaction error)
	_, err := s.dbStore.Get(ctx, db.ID)
	assert.NoError(s.T(), err)

}
