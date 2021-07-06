package seeder

import (
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"strings"
)

func caseType(id, name, partyTypeID, teamID string, template *cms.CaseTemplate) cms.CaseType {
	ct := cms.CaseType{
		ID:          id,
		Name:        name,
		PartyTypeID: partyTypeID,
		TeamID:      teamID,
		Template:    template,
	}
	caseTypes = append(caseTypes, ct)
	return ct
}

func team(id, name string) iam.Team {
	t := iam.Team{
		ID:   id,
		Name: name,
	}
	teams = append(teams, t)
	return t
}

func individual(id, firstName, lastName string) iam.Individual {
	var i = iam.Individual{
		Party: &iam.Party{
			ID: id,
			PartyTypeIDs: []string{
				iam.IndividualPartyType.ID,
			},
			Attributes: map[string][]string{
				iam.FirstNameAttribute.ID: {firstName},
				iam.LastNameAttribute.ID:  {lastName},
				iam.EMailAttribute.ID:     {strings.ToLower(firstName) + "." + strings.ToLower(lastName) + "@email.com"},
			},
		},
	}
	individuals = append(individuals, i)
	return i
}

func staff(individual iam.Individual) iam.Individual {
	individual.AddPartyType(iam.StaffPartyType.ID)
	return individual
}

func membership(id string, individual iam.Individual, team iam.Team) iam.Membership {
	m := iam.Membership{
		ID:           id,
		TeamID:       team.ID,
		IndividualID: individual.ID,
	}
	memberships = append(memberships, m)
	return m
}

func kase(id, caseTypeID, createdByID, partyID, teamID, description string, done bool) cms.Case {
	k := cms.Case{
		ID:          id,
		CaseTypeID:  caseTypeID,
		CreatorID:   createdByID,
		PartyID:     partyID,
		TeamID:      teamID,
		Description: description,
		Done:        done,
	}
	cases = append(cases, k)
	return k
}

var (
	teams         []iam.Team
	individuals   []iam.Individual
	staffers      []iam.Staff
	memberships   []iam.Membership
	relationships []iam.Relationship
	caseTypes     []cms.CaseType
	cases         []cms.Case

	// Teams
	KampalaResponseTeam = team("ac9b8d7d-d04d-4850-9a7f-3f93324c0d1e", "Kampala Response Team")
	KampalaICLATeam     = team("a43f84d5-3f8a-48c4-a896-5fb0fcd3e42b", "Kampala ICLA Team")
	NairobiResponseTeam = team("814fc372-08a6-4e6b-809b-30ebb51cb268", "Nairobi Response Team")
	NairobiICLATeam     = team("80606eb4-b53a-4fda-be12-e9806e11d44a", "Nairobi ICLA Team")

	// Case Templates
	Legal, _ = cms.NewCaseTemplate(`{"formElements": [{"type":"dropdown","attributes":{"label":"Legal status","id":"legalStatus","description":"What is the beneficiary's current legal status?","options":["Citizen","Permanent resident","Accepted refugee","Asylum seeker","Undetermined"]},"validations":{"required":true}},{"type":"checkbox","attributes":{"label":"Qualified services","id":"qualifiedServices","description":"What services does the beneficiary qualify for?","checkboxOptions":[{"label":"Councelling"},{"label":"Representation"},{"label":"Arbitration"}]},"validations":{"required":true}},{"type":"textarea","attributes":{"label":"Notes","id":"notes","description":"Additional informations, observations, concerns, etc.","placeholder":"Type something here..."}}]}`)

	// Case Types
	GenderViolence     = caseType("920ca64f-66d0-4f00-af05-d8a50ce354e6", "Gender Violence", iam.IndividualPartyType.ID, KampalaResponseTeam.ID, Legal)
	Childcare          = caseType("039caa4e-1f93-40a9-a6ff-b5bf801dfd41", "Childcare", iam.IndividualPartyType.ID, NairobiICLATeam.ID, Legal)
	HousingRights      = caseType("0d91d175-673a-4c10-bf87-3b940585e4ac", "Housing Rights", iam.IndividualPartyType.ID, KampalaICLATeam.ID, Legal)
	FinancialAssistInd = caseType("785036d9-41ee-4413-987e-13d4af761737", "Financial Assistance", iam.IndividualPartyType.ID, KampalaResponseTeam.ID, Legal)
	FinancialAssistHH  = caseType("52a24b6b-ad10-4297-b030-263bbdcd5420", "Rent Subsidy", iam.HouseholdPartyType.ID, NairobiResponseTeam.ID, Legal)

	// Individuals
	JohnDoe     = individual("c529d679-3bb6-4a20-8f06-c096f4d9adc1", "John", "Doe")
	MaryPoppins = individual("bbf539fd-ebaa-4438-ae4f-8aca8b327f42", "Mary", "Poppins")
	BoDiddley   = individual("26335292-c839-48b6-8ad5-81271ee51e7b", "Bo", "Diddley")
	Howell      = staff(individual("066a0268-fdc6-495a-9e4b-d60cfae2d81a", "Howell", "Jorg"))
	Birdie      = staff(individual("ac9015ac-686f-4719-9c3d-bf3d1cae00ea", "Birdie", "Tifawt"))
	Charis      = staff(individual("ce7ae69c-9f6a-413b-96bf-5808d0da92cd", "Charis", "Timothy"))
	Danyial     = staff(individual("8da05c97-12c2-4b43-b022-dc79be7dc3a0", "Danyial", "Hrodebert"))
	Devi        = staff(individual("6414895a-ce60-4647-b491-baeb54a76f26", "Devi", "Malvina"))
	Levan       = staff(individual("5d4c4302-ad8e-45ab-bd4c-e1ac25ae972e", "Levan", "Elija"))
	Lisbeth     = staff(individual("a5d4dab0-90d3-474d-afe6-46d04ca3caba", "Lisbeth", "Furkan"))
	Liadan      = staff(individual("c9ce906d-87ba-4123-bb74-7a73664e6778", "Liadan", "Jordaan"))
	Muhammad    = staff(individual("818206ea-0b5e-4ed9-b47e-db31566d10c0", "Muhammad", "Annemarie"))
	Dardanos    = staff(individual("7921756a-8759-4589-8a83-ad98f8aa22c7", "Dardanos", "Rilla"))
	Jana        = staff(individual("c7ca3a4d-0e96-4e5c-8c32-6750d0312706", "Jana", "Nurul"))
	Simeon      = staff(individual("78663ffb-dbaa-4362-83b6-7319d6469caa", "Simeon", "Tumelo"))
	Sayen       = staff(individual("29a20d76-dd37-471f-b9ec-9ab08f61d1ed", "Sayen", "Gezabele"))
	Veniaminu   = staff(individual("051a46b2-1ef4-4c86-bd2f-9306daedec7e", "Veniaminu", "Ye-Jun"))
	Loan        = staff(individual("f2a5d586-6865-40ea-a3db-7c729516b32b", "Loan", "Daniel"))
	Reece       = staff(individual("bdeb7e66-9129-467e-abc0-51ab2df7f222", "Reece", "Hyakinthos"))
	Svetlana    = staff(individual("afdd8b5c-b9b4-41e1-a015-7e0beb33f10b", "Svetlana", "Cerdic"))
	Kyleigh     = staff(individual("12d6a293-d923-47c6-9bc1-441934bb79c5", "Kyleigh", "Jayma"))
	Hermina     = staff(individual("dafee423-49c0-4fbf-b2f9-a42276c0cfce", "Hermina", "Magnus"))
	Leela       = staff(individual("65410229-ad41-4c17-88f2-13e9a56a0fe8", "Leela", "Cynebald"))
	Jovan       = staff(individual("bf22e83b-cfef-4c8a-b74e-f0cef6b27147", "Jovan", "Lynette"))
	Bor         = staff(individual("e350e394-091f-469c-a217-488b27b113a3", "Bor", "Lora"))
	Aldwin      = staff(individual("fdb6a682-8eb6-4565-879b-835a76384fe0", "Aldwin", "Colin"))
	Trophimos   = staff(individual("bb800fe3-85a7-4c90-b8f2-cd0354825f56", "Trophimos", "Wiebke"))

	// Memberships
	HowellMembership    = membership("862690ee-87f0-4f95-aa1e-8f8a2f2fd54a", Howell, KampalaResponseTeam)
	BirdieMembership    = membership("5fa34696-80e2-4a3e-ad26-5cc06535f67c", Birdie, KampalaResponseTeam)
	CharisMembership    = membership("72db9abe-8297-4995-b30d-bafe79f01d90", Charis, KampalaResponseTeam)
	DanyialMembership   = membership("24f6acfd-d1dd-40f8-a5b7-2b9a74d4f70b", Danyial, KampalaResponseTeam)
	DeviMembership      = membership("9811ad8e-febd-4ea0-8dba-0188eec52b94", Devi, KampalaICLATeam)
	LevanMembership     = membership("341e0a25-352e-43cb-9e5a-ffc6ce373c61", Levan, KampalaICLATeam)
	LisbethMembership   = membership("102deee5-5cf9-49c9-a9a0-99b2bde85eae", Lisbeth, KampalaICLATeam)
	BorMembership       = membership("196ad5b6-3375-4acd-83ca-1b4d6f1de19c", Bor, KampalaICLATeam)
	LiadanMembership    = membership("7f6087dd-d4d3-4a92-8c22-90bddc3b28a8", Liadan, NairobiResponseTeam)
	MuhammadMembership  = membership("113a0595-b3e3-422c-8a14-0d60ff71bb17", Muhammad, NairobiResponseTeam)
	DardanosMembership  = membership("fbbd25e9-5a2b-46d3-a8b5-a52fab5801d7", Dardanos, NairobiResponseTeam)
	JanaMembership      = membership("da3795dc-dbd9-4213-bfaa-c10764c664ba", Jana, NairobiResponseTeam)
	SimeonMembership    = membership("a7f5ad21-dd00-4d6e-92be-68c186793935", Simeon, NairobiICLATeam)
	SayenMembership     = membership("9d99551b-5cd3-4948-8695-2ee73c79f13c", Sayen, NairobiICLATeam)
	VeniaminuMembership = membership("ea2b4a53-2968-405c-9c26-8618adba6540", Veniaminu, NairobiICLATeam)
	LoanMembership      = membership("340d8740-b029-41ef-9db6-2bdf991c3ed3", Loan, NairobiICLATeam)

	// Cases
	DomesticAbuse    = kase("dba43642-8093-4685-a197-f8848d4cbaaa", GenderViolence.ID, Birdie.ID, MaryPoppins.ID, KampalaResponseTeam.ID, "Domestic abuse", false)
	MonthlyAllowance = kase("47499762-c189-4a74-9156-7969f899073b", FinancialAssistInd.ID, Birdie.ID, JohnDoe.ID, KampalaResponseTeam.ID, "Monthly allowance", false)
	ChildCare        = kase("8fb5f755-85eb-4d91-97a9-fdf86c01df25", Childcare.ID, Birdie.ID, BoDiddley.ID, KampalaResponseTeam.ID, "Monthly stipend for Bo Diddley's child", true)
)
