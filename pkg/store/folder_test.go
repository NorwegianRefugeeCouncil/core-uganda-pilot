package store

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type FolderSuite struct {
	suite.Suite
	db        *gorm.DB
	dbFactory *factory
	store     *folderStore
}

func (d *FolderSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Dialector{DSN: "file::memory:?cache=shared&_foreign_keys=1"}, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	db = db.Debug()
	if err != nil {
		d.FailNow(err.Error())
	}
	d.db = db
	if err := db.AutoMigrate(&Database{}, &Folder{}); !assert.NoError(d.T(), err) {
		d.FailNow(err.Error())
	}
	dbFactory := &factory{db: db}
	d.dbFactory = dbFactory

}

func (d *FolderSuite) SetupTest() {
	s := &folderStore{
		db: d.dbFactory,
	}
	d.store = s
}

func (s *FolderSuite) TestCreateFolder() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	folder := &types.Folder{
		ID:         uuid.NewV4().String(),
		DatabaseID: "abc",
		ParentID:   "",
		Name:       "my-folder",
	}
	got, err := s.store.Create(ctx, folder)
	if !assert.NoError(s.T(), err) {
		return
	}
	actual, err := s.store.Get(ctx, folder.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), actual, got)
}

func TestFolderSuite(t *testing.T) {
	suite.Run(t, new(FolderSuite))
}
