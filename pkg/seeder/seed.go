package seeder

import (
	"context"
	"errors"
	"fmt"

	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/client"
)

type CountryContext = string

const (
	GlobalContext                       CountryContext = "Global"
	GlobalDatabaseName                                 = "Global"
	ColombiaContext                     CountryContext = "Colombia"
	ColombiaDatabaseName                               = "Colombia"
	GlobalBioDataFolderName                            = "Global Bio Information"
	GlobalIndividualFormName                           = "Individual"
	GlobalHouseholdFormName                            = "Household"
	GlobalIndividualBeneficiaryFormName                = "Individual Beneficiary"
)

type Seed struct {
	globalDatabase            types.Database
	globalBioDataFolder       types.Folder
	globalRootIndividualForm  types.FormDefinition
	globalRootHouseholdForm   types.FormDefinition
	globalRootBeneficiaryForm types.FormDefinition
	globalForms               []types.FormDefinition
	globalBeneficiaryRefField *types.FieldDefinition
}

func NewSeed(ctx context.Context, w client.Client) (*Seed, error) {

	var seed = &Seed{}

	var dbs types.DatabaseList
	if err := client.ListDatabases(ctx, &dbs); err != nil {
		return nil, err
	}

	var folders types.FolderList
	if err := client.ListFolders(ctx, &folders); err != nil {
		return nil, err
	}
	for _, folder := range folders.Items {
		if folder.Name == GlobalBioDataFolderName {
			seed.globalBioDataFolder = *folder
		}
	}

	// find existing Global DB and Folder, if any
	for _, db := range dbs.Items {
		if db.Name == GlobalDatabaseName {
			seed.globalDatabase = *db
		}
	}

	var forms types.FormDefinitionList
	if err := client.ListForms(ctx, &forms); err != nil {
		return nil, err
	}
	for _, form := range forms.Items {
		if form.DatabaseID == seed.globalDatabase.ID {
			if form.Name == GlobalIndividualFormName {
				seed.globalRootIndividualForm = *form
			}
			if form.Name == GlobalHouseholdFormName {
				seed.globalRootHouseholdForm = *form
			}
			if form.Name == GlobalIndividualBeneficiaryFormName {
				seed.globalRootBeneficiaryForm = *form
			}
			seed.globalForms = append(seed.globalForms, *form)
		}
	}
	return seed, nil
}

func (s *Seed) Seed(ctx context.Context, client client.Client, countryContext CountryContext) error {
	var globalWasSeeded = s.globalDatabase != types.Database{} &&
		s.globalBioDataFolder != types.Folder{} &&
		len(s.globalForms) > 0

	// if not global context, ensure global context has been seeded before proceeding
	if countryContext != GlobalContext && !globalWasSeeded {
		return errors.New("the global database has not been seeded")
	}

	switch countryContext {
	case GlobalContext:
		return s.seedGlobal(ctx, client)
	case ColombiaContext:
		return s.seedColombia(ctx, client)
	default:
		return fmt.Errorf("invalid country context")
	}
}
