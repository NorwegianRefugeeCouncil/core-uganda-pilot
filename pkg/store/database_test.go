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
func (d *Suite) TestDatabaseCreate() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := d.mustCreateDatabase(ctx)
	got, err := d.dbStore.Get(ctx, db.ID)
	if !assert.NoError(d.T(), err) {
		return
	}

	assert.Equal(d.T(), db, got)
}

// TestDatabaseCreateShouldFailIfCreateSchemaFails tests that when creating a database,
// the transaction will be rolled out in case that we fail to create the actual
// SQL Schema for that database.
func (d *Suite) TestDatabaseCreateShouldFailIfCreateSchemaFails() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// mock failure to create the SQL Schema
	d.dbStore.createDatabaseSchema = func(db *gorm.DB, database *types.Database) error {
		return errors.New("mock error")
	}

	var in = &types.Database{ID: uuid.NewV4().String(), Name: "my-database"}
	if _, err := d.dbStore.Create(ctx, in); !assert.Error(d.T(), err) {
		return
	}

	// ensure the database was not stored somehow (cancelling transaction error)
	if db, err := d.dbStore.Get(ctx, in.ID); !assert.Error(d.T(), err) {
		d.T().Logf("%v", db)
		return
	}

}

// TestDatabaseCreateShouldFailIfNonUniqueDatabaseID tests that it's not possible
// to create two databases with the same ID, and that we get a
// meta.StatusReasonAlreadyExists error
func (d *Suite) TestDatabaseCreateShouldFailIfNonUniqueDatabaseID() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := d.mustCreateDatabase(ctx)

	// try recreating the same database
	_, err := d.dbStore.Create(ctx, db)
	if !assert.Error(d.T(), err) {
		return
	}

	assert.Equal(d.T(), meta.StatusReasonAlreadyExists, meta.ReasonForError(err))
}

// TestDatabaseGet tests that we can get a database after creating it
func (d *Suite) TestDatabaseGet() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := d.mustCreateDatabase(ctx)

	out, err := d.dbStore.Get(ctx, db.ID)
	if !assert.NoError(d.T(), err) {
		return
	}

	assert.Equal(d.T(), db, out)
}

// TestDatabaseGetShouldFailIfDatabaseDoesNotExist tests that we cannot get
// a database if it does not exist, and that we get a
// meta.StatusReasonNotFound
func (d *Suite) TestDatabaseGetShouldFailIfDatabaseDoesNotExist() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := d.dbStore.Get(ctx, uuid.NewV4().String())
	if !assert.Error(d.T(), err) {
		return
	}

	assert.Equal(d.T(), meta.StatusReasonNotFound, meta.ReasonForError(err))
}

// TestDatabaseListShouldReturnEmptyArrayIfNoDatabases tests that listing databases
// returns an empty list if there are no databases
func (d *Suite) TestDatabaseListShouldReturnEmptyArrayIfNoDatabases() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// clearing the databases beforehand
	if err := d.db.Where("id = id").Delete(&Database{}).Error; !assert.NoError(d.T(), err) {
		return
	}

	dbList, err := d.dbStore.List(ctx)
	if !assert.NoError(d.T(), err) {
		return
	}

	assert.Equal(d.T(), &types.DatabaseList{Items: []*types.Database{}}, dbList)
}

// TestDatabaseListShouldReturnDatabases tests that we can list the databases
func (d *Suite) TestDatabaseListShouldReturnDatabases() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db1 := d.mustCreateDatabase(ctx)
	db2 := d.mustCreateDatabase(ctx)

	dbList, err := d.dbStore.List(ctx)
	if !assert.NoError(d.T(), err) {
		return
	}

	assert.Contains(d.T(), dbList.Items, db1)
	assert.Contains(d.T(), dbList.Items, db2)
}

// TestDatabaseDelete tests that it's possible to delete a database, and that
// trying to GET the database should result in
// meta.StatusReasonNotFound
func (d *Suite) TestDatabaseDelete() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := d.mustCreateDatabase(ctx)

	if err := d.dbStore.Delete(ctx, db.ID); !assert.NoError(d.T(), err) {
		return
	}

	// test that we cannot get the database after deleting
	_, err := d.dbStore.Get(ctx, db.ID)
	if !assert.Error(d.T(), err) {
		return
	}

	assert.Equal(d.T(), meta.StatusReasonNotFound, meta.ReasonForError(err))
}

// TestDatabaseDeleteFailsIfSchemaFailsToDelete tests that when deleting a database,
// the transaction should roll back if deleting the SQL Schema fails to delete,
// and that we can still get the database afterwards
func (d *Suite) TestDatabaseDeleteFailsIfSchemaFailsToDelete() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := d.mustCreateDatabase(ctx)

	// mock that deleting the SQL Schema fails
	d.dbStore.deleteDatabaseSchema = func(db *gorm.DB, databaseId string) error {
		return errors.New("mock failure")
	}

	if err := d.dbStore.Delete(ctx, db.ID); !assert.Error(d.T(), err) {
		return
	}

	// test that we can still get the database after deleting (ensure no transaction error)
	_, err := d.dbStore.Get(ctx, db.ID)
	assert.NoError(d.T(), err)

}
