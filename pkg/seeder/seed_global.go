package seeder

import (
	"context"

	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/client"
)

func (s *Seed) seedGlobal(ctx context.Context, client client.Client) error {

	emptyDb := types.Database{}
	if s.globalDatabase == emptyDb {
		if err := client.CreateDatabase(ctx, &types.Database{
			Name: "Global",
		}, &s.globalDatabase); err != nil {
			return err
		}
	}

	emptyFolder := types.Folder{}
	if s.globalBioDataFolder == emptyFolder {
		if err := client.CreateFolder(ctx, &types.Folder{
			DatabaseID: s.globalDatabase.ID,
			Name:       GlobalBioDataFolderName,
		}, &s.globalBioDataFolder); err != nil {
			return err
		}
	}

	return s.seedGlobalForms(ctx, client)
}

func (s *Seed) seedGlobalForms(ctx context.Context, client client.Client) error {
	// Root entities ---------------------------------

	if len(s.globalRootIndividualForm.ID) == 0 {
		if err := client.CreateForm(ctx, &types.FormDefinition{
			DatabaseID: s.globalDatabase.ID,
			FolderID:   s.globalBioDataFolder.ID,
			Name:       GlobalIndividualFormName,
			Fields: types.FieldDefinitions{
				{
					Name:        "Full Name",
					Description: "The full name of the individual",
					Required:    true,
					FieldType: types.FieldType{
						Text: &types.FieldTypeText{},
					},
				}, {
					Name:        "Preferred Name",
					Description: "The name which will be used to refer to the beneficiary within Core",
					Required:    true,
					FieldType: types.FieldType{
						Text: &types.FieldTypeText{},
					},
				},
			},
		}, &s.globalRootIndividualForm); err != nil {
			return err
		}
	}

	if len(s.globalRootHouseholdForm.ID) == 0 {
		if err := client.CreateForm(ctx, &types.FormDefinition{
			DatabaseID: s.globalDatabase.ID,
			FolderID:   s.globalBioDataFolder.ID,
			Name:       GlobalHouseholdFormName,
			Fields: types.FieldDefinitions{
				{
					Name:     "Household Name",
					Required: true,
					FieldType: types.FieldType{
						Text: &types.FieldTypeText{},
					},
				},
			},
		}, &s.globalRootHouseholdForm); err != nil {
			return err
		}
	}
	if len(s.globalRootBeneficiaryForm.ID) == 0 {
		if err := client.CreateForm(ctx, &types.FormDefinition{
			DatabaseID: s.globalDatabase.ID,
			FolderID:   s.globalBioDataFolder.ID,
			Name:       GlobalIndividualBeneficiaryFormName,
			Fields: types.FieldDefinitions{
				{
					Name:     "Individual",
					Required: true,
					Key:      true,
					FieldType: types.FieldType{
						Reference: &types.FieldTypeReference{
							DatabaseID: s.globalDatabase.ID,
							FormID:     s.globalRootIndividualForm.ID,
						},
					},
				},
				yesNoQuestion("Has the beneficiary consented to NRC using their data?"),
				newFieldDefinition("URL to proof of beneficiary consent", "", false, true, types.FieldType{
					Text: &types.FieldTypeText{},
				}),
				yesNoQuestion("Beneficiary prefers to remain anonymous?"),
				yesNoQuestion("Is the beneficiary a minor?"),
				yesNoQuestion("Beneficiary presents protection concerns?"),
				yesNoQuestion("Would you say you experience some form of physical challenges?"),
				newFieldDefinition("How would you define the intensity of such challenges?", "", false, true, types.FieldType{
					SingleSelect: &types.FieldTypeSingleSelect{
						Options: wgShortSet,
					},
				}),
				yesNoQuestion("Would you say you experience some form of sensory challenges?"),
				newFieldDefinition("How would you define the intensity of such challenges?", "", false, true, types.FieldType{
					SingleSelect: &types.FieldTypeSingleSelect{
						Options: wgShortSet,
					},
				}),
				yesNoQuestion("Would you say you experience some form of mental challenges?"),
				newFieldDefinition("How would you define the intensity of such challenges?", "", false, true, types.FieldType{
					SingleSelect: &types.FieldTypeSingleSelect{
						Options: wgShortSet,
					},
				}),
				newFieldDefinition("Displacement Status", "", false, true, types.FieldType{
					SingleSelect: &types.FieldTypeSingleSelect{
						Options: globalDisplacementStatuses,
					},
				}),
				newFieldDefinition("Gender", "", false, true, types.FieldType{
					SingleSelect: &types.FieldTypeSingleSelect{
						Options: globalGenders,
					},
				}),
				newFieldDefinition("Affiliated Household", "", false, true, types.FieldType{
					Reference: &types.FieldTypeReference{
						DatabaseID: s.globalDatabase.ID,
						FormID:     s.globalRootHouseholdForm.ID,
					},
				}),
				yesNoQuestion("Are you a representative for the household?"),
			},
		}, &s.globalRootBeneficiaryForm); err != nil {
			return err

		}
	}

	return nil
}
