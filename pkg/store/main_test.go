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

type Suite struct {
	suite.Suite
	db          *gorm.DB
	dbFactory   *factory
	dbStore     *databaseStore
	folderStore *folderStore
}

func (s *Suite) SetupSuite() {
	db, err := gorm.Open(sqlite.Dialector{DSN: "file::memory:?cache=shared&_foreign_keys=1"}, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	db = db.Debug()
	if err != nil {
		s.FailNow(err.Error())
	}
	s.db = db
	if err := db.AutoMigrate(&Database{}, &Form{}, &Field{}); !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	dbFactory := &factory{
		db: db,
	}
	s.dbFactory = dbFactory
}

func (s *Suite) SetupTest() {
	s.dbStore = &databaseStore{
		createDatabaseSchema: func(db *gorm.DB, database *types.Database) error {
			return nil
		},
		deleteDatabaseSchema: func(db *gorm.DB, databaseId string) error {
			return nil
		},
		db: s.dbFactory,
	}
	s.folderStore = &folderStore{
		db: s.dbFactory,
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) createDatabase(ctx context.Context) (*types.Database, error) {
	db := &types.Database{
		ID:   uuid.NewV4().String(),
		Name: "my-db",
	}
	return s.dbStore.Create(ctx, db)
}

func (s *Suite) mustCreateDatabase(ctx context.Context) *types.Database {
	db, err := s.createDatabase(ctx)
	if !assert.NoError(s.T(), err) {
		s.T().Fail()
		return nil
	}
	return db
}

func (s *Suite) mustCreateFolder(ctx context.Context, databaseId string) *types.Folder {
	return s.mustCreateFolderWithParent(ctx, databaseId, "")
}

func (s *Suite) createFolder(ctx context.Context, databaseId string) (*types.Folder, error) {
	return s.createFolderWithParent(ctx, databaseId, "")
}

func (s *Suite) mustCreateFolderWithParent(ctx context.Context, databaseId, parentId string) *types.Folder {
	folder, err := s.createFolderWithParent(ctx, databaseId, parentId)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	return folder
}

func (s *Suite) createFolderWithParent(ctx context.Context, databaseId, parentId string) (*types.Folder, error) {
	folder := &types.Folder{
		ID:         uuid.NewV4().String(),
		DatabaseID: databaseId,
		Name:       "my-folder",
		ParentID:   parentId,
	}
	return s.folderStore.Create(ctx, folder)
}
