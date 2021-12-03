package store

import (
	"context"
	"errors"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type DatabaseSuite struct {
	suite.Suite
	db        *gorm.DB
	dbFactory *factory
	store     *databaseStore
}

func (d *DatabaseSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Dialector{DSN: "file::memory:?cache=shared"}, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		d.FailNow(err.Error())
	}
	d.db = db
	if err := db.AutoMigrate(&Database{}, &Form{}, &Field{}); !assert.NoError(d.T(), err) {
		d.FailNow(err.Error())
	}
	dbFactory := &factory{
		db: db,
	}
	d.dbFactory = dbFactory

}

func (d *DatabaseSuite) SetupTest() {
	s := &databaseStore{
		createDatabase: func(db *gorm.DB, database *types.Database) error {
			return nil
		},
		deleteDatabase: func(db *gorm.DB, databaseId string) error {
			return nil
		},
		db: d.dbFactory,
	}
	d.store = s
}

func (d *DatabaseSuite) TestCreate() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	var databaseId = uuid.NewV4().String()
	var out *types.Database
	var in = &types.Database{ID: databaseId, Name: "my-database"}

	if out, err = d.store.Create(ctx, in); !assert.NoError(d.T(), err) {
		return
	}

	assert.Equal(d.T(), databaseId, out.ID)
	assert.Equal(d.T(), "my-database", out.Name)
}

func (d *DatabaseSuite) TestCreateShouldFailIfCreateSchemaFails() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d.store.createDatabase = func(db *gorm.DB, database *types.Database) error {
		return errors.New("mock error")
	}

	var in = &types.Database{ID: uuid.NewV4().String(), Name: "my-database"}
	if _, err := d.store.Create(ctx, in); !assert.Error(d.T(), err) {
		return
	}

	// ensure the database was not stored somehow (cancelling transaction error)
	if db, err := d.store.Get(ctx, in.ID); !assert.Error(d.T(), err) {
		d.T().Logf("%v", db)
		return
	}

}
func (d *DatabaseSuite) TestCreateShouldFailIfNonUniqueDatabaseID() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var in = &types.Database{ID: uuid.NewV4().String(), Name: "my-database"}
	if _, err := d.store.Create(ctx, in); !assert.NoError(d.T(), err) {
		return
	}
	_, err := d.store.Create(ctx, in)
	if !assert.Error(d.T(), err) {
		return
	}
	assert.Equal(d.T(), meta.StatusReasonAlreadyExists, meta.ReasonForError(err))
}

func (d *DatabaseSuite) TestGetDatabase() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	in, err := d.store.Create(ctx, &types.Database{ID: uuid.NewV4().String(), Name: "my-database"})
	if !assert.NoError(d.T(), err) {
		return
	}

	out, err := d.store.Get(ctx, in.ID)
	if !assert.NoError(d.T(), err) {
		return
	}

	assert.Equal(d.T(), in, out)
}

func (d *DatabaseSuite) TestGetDatabaseShouldFailIfDatabaseDoesNotExist() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := d.store.Get(ctx, uuid.NewV4().String())
	if !assert.Error(d.T(), err) {
		return
	}
	assert.Equal(d.T(), meta.StatusReasonNotFound, meta.ReasonForError(err))
}

func (d *DatabaseSuite) TestListDatabasesShouldReturnEmptyArrayIfNoDatabases() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := d.db.Where("id = id").Delete(&Database{}).Error; !assert.NoError(d.T(), err) {
		return
	}
	dbList, err := d.store.List(ctx)
	if !assert.NoError(d.T(), err) {
		return
	}
	assert.Equal(d.T(), &types.DatabaseList{Items: []*types.Database{}}, dbList)
}

func (d *DatabaseSuite) TestListShouldReturnDatabases() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db1, err := d.store.Create(ctx, &types.Database{ID: uuid.NewV4().String(), Name: "my-database"})
	if !assert.NoError(d.T(), err) {
		return
	}
	db2, err := d.store.Create(ctx, &types.Database{ID: uuid.NewV4().String(), Name: "my-database"})
	if !assert.NoError(d.T(), err) {
		return
	}
	dbList, err := d.store.List(ctx)
	if !assert.NoError(d.T(), err) {
		return
	}
	assert.Equal(d.T(), &types.DatabaseList{Items: []*types.Database{db1, db2}}, dbList)
}

func (d *DatabaseSuite) TestDeleteDatabase() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out, err := d.store.Create(ctx, &types.Database{ID: uuid.NewV4().String(), Name: "my-database"})
	if !assert.NoError(d.T(), err) {
		return
	}

	if err := d.store.Delete(ctx, out.ID); !assert.NoError(d.T(), err) {
		return
	}

	// test that we cannot get the database after deleting
	_, err = d.store.Get(ctx, out.ID)
	if !assert.Error(d.T(), err) {
		return
	}

	assert.Equal(d.T(), meta.StatusReasonNotFound, meta.ReasonForError(err))
}

func (d *DatabaseSuite) TestDeleteDatabaseFailsIfSchemaFailsToDelete() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out, err := d.store.Create(ctx, &types.Database{ID: uuid.NewV4().String(), Name: "my-database"})
	if !assert.NoError(d.T(), err) {
		return
	}

	d.store.deleteDatabase = func(db *gorm.DB, databaseId string) error {
		return errors.New("mock failure")
	}

	if err := d.store.Delete(ctx, out.ID); !assert.Error(d.T(), err) {
		return
	}

	// test that we can still get the database after deleting (ensure no transaction error)
	_, err = d.store.Get(ctx, out.ID)
	assert.NoError(d.T(), err)

}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}
