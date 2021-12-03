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

func (d *Suite) TestDatabaseCreate() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	var databaseId = uuid.NewV4().String()
	var out *types.Database
	var in = &types.Database{ID: databaseId, Name: "my-database"}

	if out, err = d.dbStore.Create(ctx, in); !assert.NoError(d.T(), err) {
		return
	}

	assert.Equal(d.T(), databaseId, out.ID)
	assert.Equal(d.T(), "my-database", out.Name)
}

func (d *Suite) TestDatabaseCreateShouldFailIfCreateSchemaFails() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d.dbStore.createDatabase = func(db *gorm.DB, database *types.Database) error {
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
func (d *Suite) TestDatabaseCreateShouldFailIfNonUniqueDatabaseID() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var in = &types.Database{ID: uuid.NewV4().String(), Name: "my-database"}
	if _, err := d.dbStore.Create(ctx, in); !assert.NoError(d.T(), err) {
		return
	}
	_, err := d.dbStore.Create(ctx, in)
	if !assert.Error(d.T(), err) {
		return
	}
	assert.Equal(d.T(), meta.StatusReasonAlreadyExists, meta.ReasonForError(err))
}

func (d *Suite) TestDatabaseGet() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	in, err := d.dbStore.Create(ctx, &types.Database{ID: uuid.NewV4().String(), Name: "my-database"})
	if !assert.NoError(d.T(), err) {
		return
	}

	out, err := d.dbStore.Get(ctx, in.ID)
	if !assert.NoError(d.T(), err) {
		return
	}

	assert.Equal(d.T(), in, out)
}

func (d *Suite) TestDatabaseGetShouldFailIfDatabaseDoesNotExist() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := d.dbStore.Get(ctx, uuid.NewV4().String())
	if !assert.Error(d.T(), err) {
		return
	}
	assert.Equal(d.T(), meta.StatusReasonNotFound, meta.ReasonForError(err))
}

func (d *Suite) TestDatabaseListShouldReturnEmptyArrayIfNoDatabases() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := d.db.Where("id = id").Delete(&Database{}).Error; !assert.NoError(d.T(), err) {
		return
	}
	dbList, err := d.dbStore.List(ctx)
	if !assert.NoError(d.T(), err) {
		return
	}
	assert.Equal(d.T(), &types.DatabaseList{Items: []*types.Database{}}, dbList)
}

func (d *Suite) TestDatabaseListShouldReturnDatabases() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db1, err := d.dbStore.Create(ctx, &types.Database{ID: uuid.NewV4().String(), Name: "my-database"})
	if !assert.NoError(d.T(), err) {
		return
	}
	db2, err := d.dbStore.Create(ctx, &types.Database{ID: uuid.NewV4().String(), Name: "my-database"})
	if !assert.NoError(d.T(), err) {
		return
	}
	dbList, err := d.dbStore.List(ctx)
	if !assert.NoError(d.T(), err) {
		return
	}

	assert.Contains(d.T(), dbList.Items, db1)
	assert.Contains(d.T(), dbList.Items, db2)
}

func (d *Suite) TestDatabaseDelete() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out, err := d.dbStore.Create(ctx, &types.Database{ID: uuid.NewV4().String(), Name: "my-database"})
	if !assert.NoError(d.T(), err) {
		return
	}

	if err := d.dbStore.Delete(ctx, out.ID); !assert.NoError(d.T(), err) {
		return
	}

	// test that we cannot get the database after deleting
	_, err = d.dbStore.Get(ctx, out.ID)
	if !assert.Error(d.T(), err) {
		return
	}

	assert.Equal(d.T(), meta.StatusReasonNotFound, meta.ReasonForError(err))
}

func (d *Suite) TestDatabaseDeleteFailsIfSchemaFailsToDelete() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out, err := d.dbStore.Create(ctx, &types.Database{ID: uuid.NewV4().String(), Name: "my-database"})
	if !assert.NoError(d.T(), err) {
		return
	}

	d.dbStore.deleteDatabase = func(db *gorm.DB, databaseId string) error {
		return errors.New("mock failure")
	}

	if err := d.dbStore.Delete(ctx, out.ID); !assert.Error(d.T(), err) {
		return
	}

	// test that we can still get the database after deleting (ensure no transaction error)
	_, err = d.dbStore.Get(ctx, out.ID)
	assert.NoError(d.T(), err)

}
