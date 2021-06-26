package seed

import (
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"strings"
)

var (

	// KAMPALA RESPONSE TEAM
	Howell           = individual("066a0268-fdc6-495a-9e4b-d60cfae2d81a", "Howell", "Jorg")
	HowellStaff      = staf("8a2497f9-1608-405a-adc9-6c29147f6a06", Howell)
	HowellMembership = mbship("", Howell, KampalaResponseTeam)

	Birdie           = individual("ac9015ac-686f-4719-9c3d-bf3d1cae00ea", "Birdie", "Tifawt")
	BirdieStaff      = staf("fec7ea88-c77c-4311-8824-8e6aa65c2702", Birdie)
	BirdieMembership = mbship("", Birdie, KampalaResponseTeam)

	Charis           = individual("ce7ae69c-9f6a-413b-96bf-5808d0da92cd", "Charis", "Timothy")
	CharisStaff      = staf("3cf8f80b-31b3-4744-8743-abfabbe6afc7", Charis)
	CharisMembership = mbship("", Charis, KampalaResponseTeam)

	Danyial           = individual("8da05c97-12c2-4b43-b022-dc79be7dc3a0", "Danyial", "Hrodebert")
	DanyialStaff      = staf("87d3cf68-5141-4838-bc8a-145deb0587ca", Danyial)
	DanyialMembership = mbship("", Danyial, KampalaResponseTeam)

	// KAMPALA ICLA TEAM
	Devi           = individual("6414895a-ce60-4647-b491-baeb54a76f26", "Devi", "Malvina")
	DeviStaff      = staf("0573c9a9-94ed-4353-a8bd-5c8422d8a66c", Devi)
	DeviMembership = mbship("", Devi, KampalaICLATeam)

	Levan           = individual("5d4c4302-ad8e-45ab-bd4c-e1ac25ae972e", "Levan", "Elija")
	LevanStaff      = staf("a6f32747-558e-4b2a-a767-d6313ea83a17", Levan)
	LevanMembership = mbship("", Levan, KampalaICLATeam)

	Lisbeth           = individual("a5d4dab0-90d3-474d-afe6-46d04ca3caba", "Lisbeth", "Furkan")
	LisbethStaff      = staf("3fce53f1-2583-438f-b5e6-e70efcaaa935", Lisbeth)
	LisbethMembership = mbship("", Lisbeth, KampalaICLATeam)

	Bor           = individual("e350e394-091f-469c-a217-488b27b113a3", "Bor", "Lora")
	BorStaff      = staf("ac7f8a3c-0f7e-49b3-b11d-5135cf3a3af0", Bor)
	BorMembership = mbship("", Bor, KampalaICLATeam)

	// NAIROBI RESPONSE TEAM
	Liadan           = individual("c9ce906d-87ba-4123-bb74-7a73664e6778", "Liadan", "Jordaan")
	LiadanStaff      = staf("168f3397-6f1c-4607-b0a2-e130f8dd9cc4", Liadan)
	LiadanMembership = mbship("", Liadan, NairobiResponseTeam)

	Muhammad           = individual("818206ea-0b5e-4ed9-b47e-db31566d10c0", "Muhammad", "Annemarie")
	MuhammadStaff      = staf("e559a5d1-7293-4c2f-8ac2-c4aa4ee12571", Muhammad)
	MuhammadMembership = mbship("", Muhammad, NairobiResponseTeam)

	Dardanos           = individual("7921756a-8759-4589-8a83-ad98f8aa22c7", "Dardanos", "Rilla")
	DardanosStaff      = staf("ef162397-b32d-4e4f-8418-606265b937b1", Dardanos)
	DardanosMembership = mbship("", Dardanos, NairobiResponseTeam)

	Jana           = individual("c7ca3a4d-0e96-4e5c-8c32-6750d0312706", "Jana", "Nurul")
	JanaStaff      = staf("58d10f3a-4e1d-4a2b-8b19-109d5fe540c7", Jana)
	JanaMembership = mbship("", Jana, NairobiResponseTeam)

	// NAIROBI ICLA TEAM
	Simeon           = individual("78663ffb-dbaa-4362-83b6-7319d6469caa", "Simeon", "Tumelo")
	SimeonStaff      = staf("79fcf4dd-2603-4d7f-891d-604529ecfe3c", Simeon)
	SimeonMembership = mbship("", Simeon, NairobiICLATeam)

	Sayen           = individual("29a20d76-dd37-471f-b9ec-9ab08f61d1ed", "Sayen", "Gezabele")
	SayenStaff      = staf("540fd5da-e152-4401-b30a-27c5a661a4f7", Sayen)
	SayenMembership = mbship("", Sayen, NairobiICLATeam)

	Veniaminu           = individual("051a46b2-1ef4-4c86-bd2f-9306daedec7e", "Veniaminu", "Ye-Jun")
	VeniaminuStaff      = staf("a2d7047b-b126-4213-a0a5-7ef4ecb09515", Veniaminu)
	VeniaminuMembership = mbship("", Veniaminu, NairobiICLATeam)

	Loan           = individual("f2a5d586-6865-40ea-a3db-7c729516b32b", "Loan", "Daniel")
	LoanStaff      = staf("304ae2f9-3cf7-4aed-a26d-24a400788bc1", Loan)
	LoanMembership = mbship("", Loan, NairobiICLATeam)

	// NO MEMBERSHIPS IN TEAMS (ENUMERATORS)
	Reece      = individual("bdeb7e66-9129-467e-abc0-51ab2df7f222", "Reece", "Hyakinthos")
	ReeceStaff = staf("b9c847b2-fefa-4223-8fdb-a048ef3153fc", Reece)

	Svetlana      = individual("afdd8b5c-b9b4-41e1-a015-7e0beb33f10b", "Svetlana", "Cerdic")
	SvetlanaStaff = staf("3b88b8a2-5cfd-403e-bc5f-724479418a94", Svetlana)

	Kyleigh      = individual("12d6a293-d923-47c6-9bc1-441934bb79c5", "Kyleigh", "Jayma")
	KyleighStaff = staf("13fba20e-07dc-49a8-bd41-bdcc796d53e0", Kyleigh)

	Hermina      = individual("dafee423-49c0-4fbf-b2f9-a42276c0cfce", "Hermina", "Magnus")
	HerminaStaff = staf("56b35e75-5703-46b6-9232-51570f0c73bc", Hermina)

	Leela      = individual("65410229-ad41-4c17-88f2-13e9a56a0fe8", "Leela", "Cynebald")
	LeelaStaff = staf("b0ca2145-c884-4aaf-88f0-39f2c8e9acdf", Leela)

	Jovan      = individual("bf22e83b-cfef-4c8a-b74e-f0cef6b27147", "Jovan", "Lynette")
	JovanStaff = staf("265ffaf6-fe44-4290-9755-17894b3c2e5f", Jovan)

	Aldwin      = individual("fdb6a682-8eb6-4565-879b-835a76384fe0", "Aldwin", "Colin")
	AldwinStaff = staf("31c1cdfc-cda2-4012-bd21-9ffe0939cda5", Aldwin)

	Trophimos      = individual("bb800fe3-85a7-4c90-b8f2-cd0354825f56", "Trophimos", "Wiebke")
	TrophimosStaff = staf("857e197c-19f8-4546-a2cb-769d42a0fd55", Trophimos)
)

var people []*iam.Party

func individual(id, firstName, lastName string) iam.Party {
	var party = iam.Party{
		ID: id,
		PartyTypeIDs: []string{
			iam.IndividualPartyType.ID,
		},
		Attributes: map[string][]string{
			iam.FirstNameAttribute.ID: {firstName},
			iam.LastNameAttribute.ID:  {lastName},
			iam.EMailAttribute.ID:     {strings.ToLower(firstName) + "." + strings.ToLower(lastName) + "@email.com"},
		},
	}
	people = append(people, &party)
	return party
}

var staffs []*iam.Staff

func staf(id string, individual iam.Party) *iam.Staff {
	s := &iam.Staff{
		ID:             id,
		OrganizationID: NRC.ID,
		IndividualID:   individual.ID,
	}
	staffs = append(staffs, s)
	return s
}

var mbships []*iam.Membership

func mbship(id string, individual iam.Party, team iam.Team) *iam.Membership {
	m := &iam.Membership{
		ID:           id,
		TeamID:       team.ID,
		IndividualID: individual.ID,
	}
	mbships = append(mbships, m)
	return m
}
