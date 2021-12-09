package seeder

import (
	"github.com/nrc-no/core/pkg/api/types"
)

func (s *Seed) seedGlobal() error {
	var dbConfig = &types.Database{
		Name: "Global",
	}
	globalDB, err := s.dbStore.Create(s.ctx, dbConfig)
	if err != nil {
		return err
	}

	globalBioDataFolder, err := s.folderStore.Create(s.ctx, &types.Folder{
		DatabaseID: globalDB.ID,
		Name:       GlobalBioDataFolderName,
	})

	s.globalDatabaseId = globalDB.ID
	s.globalBioDataFolderId = globalBioDataFolder.ID
	return s.seedGlobalForms()
}

func (s *Seed) seedGlobalForms() error {
	// Root entities ---------------------------------
	rootIndividualForm, err := s.formStore.Create(s.ctx, newFormDefinition(
		s.globalDatabaseId,
		s.globalBioDataFolderId,
		"Individual",
		[]*types.FieldDefinition{
			newFieldDefinition("Full Name", "The full name of the individual", false, true, types.FieldType{
				Text: &types.FieldTypeText{},
			}),
			newFieldDefinition("Preferred Name", "The name which will be used to refer to the beneficiary within Core", false, true, types.FieldType{
				Text: &types.FieldTypeText{},
			}),
		},
	))
	if err != nil {
		return err
	}

	rootHouseholdForm, err := s.formStore.Create(s.ctx, newFormDefinition(
		s.globalDatabaseId,
		s.globalBioDataFolderId,
		"Household",
		[]*types.FieldDefinition{
			newFieldDefinition("Household Name", "", false, true, types.FieldType{
				Text: &types.FieldTypeText{},
			}),
		},
	))
	if err != nil {
		return err
	}

	// Global Intake ---------------------------------
	rootIndividualBeneficiaryForm, err := s.formStore.Create(s.ctx, newFormDefinition(
		s.globalDatabaseId,
		s.globalBioDataFolderId,
		"Individual Beneficiary",
		[]*types.FieldDefinition{
			newFieldDefinition("Individual", "Individual who is being registered as a beneficiary", false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: s.globalDatabaseId,
					FormID:     rootIndividualForm.ID,
				},
			}),
			yesNoQuestion("Consent"),
			newFieldDefinition("Consent URL", "Link to proof of consent", false, true, types.FieldType{
				Text: &types.FieldTypeText{},
			}),
			yesNoQuestion("Anonymous"),
			newFieldDefinition("Minor", "Is this beneficiary a minor", false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{yesNoChoice},
			}),
			yesNoQuestion("Protection Concern"),
			yesNoQuestion("Physical Disability"),
			newFieldDefinition("Physical Disability - Intensity", "How would you define the intensity of such challenges?", false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{wgShortSet},
			}),
			yesNoQuestion("Sensory Disability"),
			newFieldDefinition("Sensory Disability - Intensity", "How would you define the intensity of such challenges?", false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{wgShortSet},
			}),
			yesNoQuestion("Mental Disability"),
			newFieldDefinition("Mental Disability - Intensity", "How would you define the intensity of such challenges?", false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{wgShortSet},
			}),
			newFieldDefinition("Displacement Status", "", false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{globalDisplacementStatuses},
			}),
			newFieldDefinition("Gender", "", false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{globalGenders},
			}),
			newFieldDefinition("Affiliated Household", "Household to which this beneficiary belongs", false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: s.globalDatabaseId,
					FormID:     rootHouseholdForm.ID,
				},
			}),
			yesNoQuestion("Representative of Household"),
		},
	))
	if err != nil {
		return err
	}

	s.globalRootIndividualFormId = rootIndividualForm.ID
	s.globalRootHouseholdFormId = rootHouseholdForm.ID
	s.globalRootBeneficiaryFormId = rootIndividualBeneficiaryForm.ID

	return nil
}
