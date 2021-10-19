package seeder

var (
	// Identification Document Types
	DriversLicense = identificationDocumentType("75c41c5f-bf7e-4b45-a242-5e0f875e3044", "Drivers License")
	NationalID     = identificationDocumentType("8910a1ea-4bfe-4321-aa5b-15922b09ad4d", "National ID")
	UNHCRID        = identificationDocumentType("6833cb6d-593f-4f3f-926d-498be74352d1", "UNHCR ID")
	Passport       = identificationDocumentType("567d04e5-abf4-4899-848f-0395264309f0", "Passport")
	OPMID		   = identificationDocumentType("bcb85d3d-7474-44e4-bbed-ac3d91737df0", "OPM ID")

	_ = identificationDocument("20d194d6-a1ac-483e-8c24-38b5efbaca6f", BoDiddley.ID, "A0JBODIDDLEY129", Passport.ID)
	_ = identificationDocument("0244b59e-5d5c-4e13-af96-da1ccf4e9499", MaryPoppins.ID, "LLP987MARYPOPPINS99", UNHCRID.ID)
	_ = identificationDocument("4c9477c9-c149-4db7-928c-f5e5f915e018", JohnDoe.ID, "B811HJOHNDOE01", NationalID.ID)
)
