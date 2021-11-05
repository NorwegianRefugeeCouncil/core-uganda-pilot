package seeder

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CountryContext = string

const (
	GlobalContext CountryContext = "global"
	UgandaContext CountryContext = "uganda"
)

func Seed(countryContext CountryContext, dbFactory store.Factory) error {

	db, err := gorm.Open(sqlite.Dialector{DSN: "file::memory:?cache=shared&_foreign_keys=1"}, &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&store.Database{}, &store.Form{}, &store.Field{}); err != nil {
		return err
	}
	return err
}

dbStore := store.NewDatabaseStore(dbFactory)
folderStore := store.NewFolderStore(dbFactory)
formStore := store.NewFormStore(dbFactory)

switch countryContext {
case GlobalContext:
err = seedGlobal(ctx, dbStore, folderStore, formStore)
case UgandaContext:
err = seedUganda(ctx, dbStore, folderStore, formStore)
}

return err
}

func seedGlobal(ctx context.Context, dbStore store.DatabaseStore, folderStore store.FolderStore, formStore store.FormStore) error {
	// TODO
	return nil
}
func seedUganda(ctx context.Context, dbStore store.DatabaseStore, folderStore store.FolderStore, formStore store.FormStore) error {
	var dbConfig = &types.Database{
		Name: "Uganda",
	}
	ugDB, err := dbStore.Create(ctx, dbConfig)
	if err != nil {
		return err
	}

	coFolder, err := folderStore.Create(ctx, &types.Folder{
		DatabaseID: ugDB.ID,
		Name:       "Uganda Bio Information Folder",
	})
	if err != nil {
		return err
	}

	iclaFolder, err := folderStore.Create(ctx, &types.Folder{
		DatabaseID: ugDB.ID,
		Name:       "Uganda ICLA Folder",
	})
	if err != nil {
		return err
	}

	intakeFolder, err := folderStore.Create(ctx, &types.Folder{
		DatabaseID: ugDB.ID,
		Name:       "Uganda Intake Folder",
	})
	if err != nil {
		return err
	}

	protectionFolder, err := folderStore.Create(ctx, &types.Folder{
		DatabaseID: ugDB.ID,
		Name:       "Uganda Protection Folder",
	})
	if err != nil {
		return err
	}

	if err := seedUgandaForms(ctx, formStore, ugDB.ID, coFolder.ID, iclaFolder.ID, intakeFolder.ID, protectionFolder.ID); err != nil {
		return err
	}
	return nil
}
