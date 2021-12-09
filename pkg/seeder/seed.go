package seeder

import (
	"context"
	"errors"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CountryContext = string

const (
	GlobalContext           CountryContext = "Global"
	UgandaContext           CountryContext = "Uganda"
	GlobalDatabaseName                     = "Global"
	GlobalBioDataFolderName                = "Global Bio Information"
)

type Seed struct {
	ctx                         context.Context
	dbStore                     store.DatabaseStore
	folderStore                 store.FolderStore
	formStore                   store.FormStore
	globalDatabaseId            string
	globalBioDataFolderId       string
	globalRootIndividualFormId  string
	globalRootHouseholdFormId   string
	globalRootBeneficiaryFormId string
	globalForms                 *types.FormDefinitionList
}

func NewSeed(ctx context.Context, dbFactory store.Factory) (*Seed, error) {
	dbStore := store.NewDatabaseStore(dbFactory)
	folderStore := store.NewFolderStore(dbFactory)
	formStore := store.NewFormStore(dbFactory)

	dbs, err := dbStore.List(ctx)
	if err != nil {
		return nil, err
	}

	var seed = &Seed{
		ctx:         ctx,
		dbStore:     dbStore,
		folderStore: folderStore,
		formStore:   formStore,
	}

	// find existing Global DB and Folder, if any
	for _, db := range dbs.Items {
		if db.Name == GlobalDatabaseName {
			seed.globalDatabaseId = db.ID
			folders, err := folderStore.List(ctx)
			if err != nil {
				return nil, err
			}
			for _, folder := range folders.Items {
				if folder.Name == GlobalBioDataFolderName {
					seed.globalBioDataFolderId = folder.ID
				}
			}
			forms, err := formStore.List(ctx)
			seed.globalForms = forms
		}
	}

	return seed, nil
}

func (s *Seed) Seed(countryContext CountryContext) error {
	db, err := gorm.Open(sqlite.Dialector{DSN: "file::memory:?cache=shared&_foreign_keys=1"}, &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&store.Database{}, &store.Form{}, &store.Field{}); err != nil {
		return err
	}

	var globalWasSeeded = len(s.globalDatabaseId) > 0 && len(s.globalBioDataFolderId) > 0 && s.globalForms.Len() > 0

	// if not global context, ensure global context has been seeded before proceeding
	if countryContext != GlobalContext && !globalWasSeeded {
		return errors.New("the global database hase not been seeded")
	}

	switch countryContext {
	case GlobalContext:
		err = s.seedGlobal()
	case UgandaContext:
		err = s.seedUganda()
	}

	return err
}
