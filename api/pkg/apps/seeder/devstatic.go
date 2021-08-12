package seeder

import (
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/form"
)

// Dev static objects.
var (
	// Case Types for Dogfooding
	DTeamBugReportCaseType      = caseType("39b24aaa-02a3-4455-b3a6-fd05e6a59fef", "Report a bug in Core", iam.IndividualPartyType.ID, DTeam.ID, DTeamBugReport, false)
	DTeamFeatureRequestCaseType = caseType("95bf45fd-a703-4698-ae9c-12f1865b1a6f", "Request a feature/change in Core", iam.IndividualPartyType.ID, DTeam.ID, DTeamFeatureRequest, false)

	// D-Team (Dogfooding)
	TestIndividual = individual("3a38a57c-660d-4482-975c-e49e7618c689", "Test", "Individual", "01/01/1234", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Ludovic        = staff(individual("78b494dc-7461-42f5-bf2d-1c9695e63ba8", "Ludovic", "Cleroux", "12/02/1978", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Cassie         = staff(individual("dd65e4cf-c691-411a-a1f8-bed22c538480", "Cassie", "Seo", "12/02/1978", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Senyao         = staff(individual("1889acbb-5dbb-4998-a071-ab00c19c2b77", "Senyao", "Hou", "12/02/1978", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Robert         = staff(individual("3e8488eb-785a-49c4-95f1-2cc5c09e8ab9", "Robert", "Focke", "12/02/1978", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Nicolas        = staff(individual("7c1107b7-3fa7-4f49-acea-e953c5d8723f", "Nicolas", "Epstein", "12/02/1978", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Kristjan       = staff(individual("ae4d0fd5-bb03-4b9d-948d-c99754aca5ce", "Kristjan", "Thoroddsson", "12/02/1978", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))

	DTeam = team("5efecef3-a7a9-4705-84ca-70c89d7d783f", "D-Team")

	// D-Team Memberships
	LudovicMembership  = membership("156e7e3a-6cec-43ca-be28-94e8eb0bb27c", Ludovic, DTeam)
	CassieMembership   = membership("16f91d0b-d53a-41cc-a437-5124fd65656e", Cassie, DTeam)
	SenyaoMembership   = membership("fdad7109-5fde-41ca-8eee-3f699ad8e491", Senyao, DTeam)
	RobertMembership   = membership("2cc4d8e7-2087-41d0-af7f-90144820466f", Robert, DTeam)
	NicolasMembership  = membership("624291ac-d573-4866-b507-d9e83b9b2288", Nicolas, DTeam)
	KristjanMembership = membership("a6a7a318-64d8-4cfa-83c7-8710f1d12778", Kristjan, DTeam)

	// Dogfooding Case Templates
	DTeamBugReport = &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: "textarea",
				Attributes: form.FormElementAttributes{
					Label:       "What action were you undertaking in the application, when the error happened",
					Name:        "whatActionBeforeError",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: "textarea",
				Attributes: form.FormElementAttributes{
					Label:       "If the error had not happened, what would be your expected outcome for the action you were performing when the error happened",
					Name:        "expectedOutcome",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: "textarea",
				Attributes: form.FormElementAttributes{
					Label:       "List any error messages shown",
					Name:        "errorMessages",
					Description: "",
					Placeholder: "",
				},
			},
		},
	}

	DTeamFeatureRequest = &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: "textarea",
				Attributes: form.FormElementAttributes{
					Label:       "Describe the change or new functionality you would like in Core",
					Name:        "request",
					Description: "",
					Placeholder: "",
				},
			},
		},
	}

	// TestTemplate !!! It's important to keep this template up-to-date for e2e testing.
	// It MUST include an instance of each input type described in cms/iam.
	TestTemplate = &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: "dropdown",
				Attributes: form.FormElementAttributes{
					Label:       "Dropdown",
					Name:        "testDropdown",
					Description: "Dropdown description",
					Options:     []string{"0", "1", "2"},
				},
				Validation: form.FormElementValidation{Required: false},
			},
			{
				Type: "checkbox",
				Attributes: form.FormElementAttributes{
					Label:           "Checkbox",
					Name:            "testCheckbox",
					Description:     "Checkbox description",
					CheckboxOptions: []form.CheckboxOption{{Label: "0"}, {Label: "1"}, {Label: "2"}},
				},
				Validation: form.FormElementValidation{Required: false},
			},
			{
				Type: "textarea",
				Attributes: form.FormElementAttributes{
					Label:       "Textarea",
					Name:        "testTextarea",
					Description: "Textarea description",
					Placeholder: "Textarea placeholder",
				},
				Validation: form.FormElementValidation{Required: false},
			},
			{
				Type: "text",
				Attributes: form.FormElementAttributes{
					Label:       "Textinput",
					Name:        "testTextinput",
					Description: "Textinput description",
					Placeholder: "Textinput placeholder",
				},
				Validation: form.FormElementValidation{Required: false},
			},
		},
	}

	TestCaseType = caseType("05b3460e-9f20-4af7-9a74-06f632f7ae24", "Test Case Type", iam.IndividualPartyType.ID, DTeam.ID, TestTemplate, false)

	TestCase = kase("5827081c-e6b8-4b29-a1b1-f6780a65c28e", TestCaseType.ID, "D-Team", DTeam.ID, DTeam.ID, false, TestCaseType.Template, TestCaseType.IntakeCaseType)
)
