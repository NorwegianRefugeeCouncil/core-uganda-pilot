package seeder

import (
	"context"

	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/client"
)

func (s *Seed) seedColombia(ctx context.Context, client client.Client) error {
	var err error

	var co_db types.Database

	err = client.CreateDatabase(ctx, &types.Database{Name: ColombiaDatabaseName}, &co_db)

	if err != nil {
		return err
	}

	// keep a reference to the global beneficiary form as a reference field for use in the forms
	s.globalBeneficiaryRefField = &types.FieldDefinition{
		Name:     "Individual Beneficiary",
		Key:      true,
		Required: true,
		FieldType: types.FieldType{
			Reference: &types.FieldTypeReference{
				DatabaseID: s.globalDatabase.ID,
				FormID:     s.globalRootBeneficiaryForm.ID,
			},
		},
	}

	err = s.seedCoIntake(ctx, client, co_db.ID)

	if err != nil {
		return err
	}

	err = s.seedCoConsent(ctx, client, co_db.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Seed) seedCoIntake(ctx context.Context, client client.Client, dbID string) error {
	var (
		err                    error
		intake_folder          types.Folder
		intake_pii             *types.FormDefinition
		intake_id_details      *types.FormDefinition
		intake_nrc_details     *types.FormDefinition
		intake_contact_info    *types.FormDefinition
		intake_ind_cc_specific *types.FormDefinition
		intake_hh_information  *types.FormDefinition
	)

	err = client.CreateFolder(ctx, &types.Folder{
		Name:       "Intake",
		DatabaseID: dbID,
	}, &intake_folder)

	if err != nil {
		return err
	}

	intake_pii = &types.FormDefinition{
		Name:       "Colombia Individual",
		DatabaseID: dbID,
		FolderID:   intake_folder.ID,
		Type:       types.RecipientFormType,
		Fields: types.FieldDefinitions{
			s.globalBeneficiaryRefField,
			date("Date of birth", true),
			dropdown("Nationality (1)", co_nationality1, false),
			dropdown("Nationality (2)", co_nationality2, false),
			dropdown("Marital status", co_marital_status, true),
			dropdown("Relationship to HH representative", co_relationship_to_hh, true),
			dropdown("Beneficiary type", co_beneficiary_type, true),
			dropdown("Ethnicity", co_ethnicity, true),
			yesNo("Do you have a job or entrepreneurship"),
			text("What sector? (commerce, production, service, agro)", false),
			text("Entrepreneurship time", false),
			text("Type of job (contract type)", false),
			textarea("Types of income sources in the family", false),
			text("Name and surname of the legal representative", false),
			text("Additional information about the legal representative", false),
			textarea("Reasons for the representation", false),
			yesNo("Is the guardianship legal under national law?"),
			textarea("Add the legal evaluation (if the answer is positive)", false),
			textarea("Offer assistance in identifying a recognized legal representative (if the answer is negative)", false),
			yesNo("Can the person give their consent legally?"),
			quantity("Age of head of household (if another person, and if the household is not registered)", false),
			quantity("Monthly household income / per capita", false),
			yesNo("Is the head of household is pregnant or nursing?"),
			yesNo("Does the head of household have a chronic illness?"),
		},
	}

	// TODO: This should be in the global "Identification Documents"
	intake_id_details = &types.FormDefinition{
		Name:       "ID Details",
		DatabaseID: dbID,
		FolderID:   intake_folder.ID,
		Fields: types.FieldDefinitions{
			s.globalBeneficiaryRefField,
			dropdown("Type of identification", co_identification_type, true),
			ifOtherPleaseSpecify,
			text("Identification number", true),
			text("Additional ID #1", false),
			text("Additional ID #2", false),
		},
	}

	// TODO: This should be a SubForm ?
	intake_nrc_details = &types.FormDefinition{
		Name:       "NRC Details",
		DatabaseID: dbID,
		FolderID:   intake_folder.ID,
		Fields: types.FieldDefinitions{
			s.globalBeneficiaryRefField,
			dropdown("Source of identification", co_identification_source, true),
			text("Location name", true),
			date("Registration date", false),
			dropdown("Country", co_countries, true),
			dropdown("District", co_admin_2, true),
			text("Subcounty", false),
			dropdown("Type of settlement", co_settlement_type, false),
			text("Parish", false),
			text("Village", false),
			yesNo("Emergency attention?"),
			yesNo("Sustainable solutions?"),
			yesNo("Hard to reach area?"),
			yesNo("COVID-19 emergency?"),
			text("How did you learn about NRC?", true),
		},
	}

	// TODO: Should this be a SubForm ?
	intake_contact_info = &types.FormDefinition{
		Name:       "Contact Information",
		DatabaseID: dbID,
		FolderID:   intake_folder.ID,
		Fields: types.FieldDefinitions{
			s.globalBeneficiaryRefField,
			textarea("Spoken languages", true),
			text("Preferred language", true),
			text("Email address", false),
			textarea("Physical address", false),
			text("Phone number (1)", false),
			text("Phone number (2)", false),
			text("Phone number (3)", false),
			dropdown("Preferred means of contact", co_means_of_contact, true),
			yesNo("Requires an interpreter?"),
			text("Interpreter name", false),
		},
	}

	// TODO: Should this be a SubForm ?
	intake_ind_cc_specific = &types.FormDefinition{
		Name:       "Individual CC Specific",
		DatabaseID: dbID,
		FolderID:   intake_folder.ID,
		Fields: types.FieldDefinitions{
			s.globalBeneficiaryRefField,
			text("Case number", true),
			text("Caseworker/Paralegal's ID number", false),
			dropdown("Modality of service delivery (taxonomy)", co_service_delivery_modalities, false),
			text("Living situation", true),
			textarea("Comment on living situation", false),
			dropdown("How did you learn about our services?", co_intro_source, true),
		},
	}

	intake_hh_information = &types.FormDefinition{
		Name:       "Household Information",
		DatabaseID: dbID,
		FolderID:   intake_folder.ID,
		Fields: types.FieldDefinitions{
			s.globalBeneficiaryRefField,
			quantity("Number of 60+ years old males", true),
			text("Household/OPM Number", false),
			text("Scholarship level of head of household", false),
			quantity("Total number of people who are part of the household", true),
			quantity("# of household members who have a disability", false),
			quantity("# of household members who are pregnant or nursing", false),
			quantity("# of household members who have a chronic illness", false),
		},
	}

	var forms = []*types.FormDefinition{
		intake_pii,
		intake_id_details,
		intake_nrc_details,
		intake_contact_info,
		intake_ind_cc_specific,
		intake_hh_information,
	}

	for _, form := range forms {
		err = client.CreateForm(ctx, form, nil)
		if err != nil {
			return err
		}
	}

	return err
}

func (s *Seed) seedCoConsent(ctx context.Context, client client.Client, dbID string) error {
	var (
		err                error
		consent_folder     types.Folder
		consent_individual *types.FormDefinition
	)

	err = client.CreateFolder(ctx, &types.Folder{
		Name:       "Consent",
		DatabaseID: dbID,
	}, &consent_folder)

	if err != nil {
		return err
	}

	consent_individual = &types.FormDefinition{
		Name:       "Individual Consent",
		DatabaseID: dbID,
		FolderID:   consent_folder.ID,
		Fields: types.FieldDefinitions{
			s.globalBeneficiaryRefField,
			text("Upload the consent form signed by the beneficiary", true),
			yesNo("Can NRC staff initiate contact with the beneficiary?"),
			yesNo("Format of the act of commitment to the program"),
			yesNo("Format of delivery for assets and inputs"),
			textarea("Other consent information", false),
			yesNo("Can NRC share the beneficiary's data?"),
		},
	}

	err = client.CreateForm(ctx, consent_individual, nil)

	return err
}
