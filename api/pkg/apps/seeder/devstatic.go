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
	Ludovic  = staff(individual("78b494dc-7461-42f5-bf2d-1c9695e63ba8", "Ludovic Cleroux", "Ludovic Cleroux", "12/02/1978", "ludovic.cleroux", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Cassie   = staff(individual("dd65e4cf-c691-411a-a1f8-bed22c538480", "Cassie Seo", "Cassie Seo", "12/02/1978", "cassie.seo", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Senyao   = staff(individual("1889acbb-5dbb-4998-a071-ab00c19c2b77", "Senyao Hou", "Senyao Hou", "12/02/1978", "senyao.hou", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Robert   = staff(individual("3e8488eb-785a-49c4-95f1-2cc5c09e8ab9", "Robert Focke", "Robert Focke", "12/02/1978", "robert.focke", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Nicolas  = staff(individual("7c1107b7-3fa7-4f49-acea-e953c5d8723f", "Nicolas Epstein", "Nicolas Epstein", "12/02/1978", "nicolas.epstein", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Kristjan = staff(individual("ae4d0fd5-bb03-4b9d-948d-c99754aca5ce", "Kristjan Thoroddsson", "Kristjan Thoroddsson", "12/02/1978", "kristjan.thoroddsson", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))

	DTeam = team("5efecef3-a7a9-4705-84ca-70c89d7d783f", "D-Team")

	globalCountry = country(iam.GlobalCountry.ID, iam.GlobalCountry.Name)
	DTeamNationality  = nationality("8f6be182-f64c-4096-ba8e-4562506dae6d", DTeam, globalCountry)

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
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "What action were you undertaking in the application, when the error happened",
					Name:        "whatActionBeforeError",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "If the error had not happened, what would be your expected outcome for the action you were performing when the error happened",
					Name:        "expectedOutcome",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
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
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Describe the change or new functionality you would like in Core",
					Name:        "request",
					Description: "",
					Placeholder: "",
				},
			},
		},
	}
)
