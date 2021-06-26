package seed

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/apps/cms"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

var (
	teams         []iam.Team
	orgs          []iam.Organization
	individuals   []iam.Individual
	staffers      []iam.Staff
	memberships   []iam.Membership
	relationships []iam.Relationship
	caseTypes     []cms.CaseType
	cases         []cms.Case

	// Organizations
	NRC = iam.Organization{
		Party: &iam.Party{
			PartyTypeIDs: []string{iam.OrganizationPartyType.ID},
			Attributes: map[string][]string{
				iam.LegalNameAttribute.ID: {"NRC"},
			},
		},
	}

	// Teams
	KampalaResponseTeam = team("ac9b8d7d-d04d-4850-9a7f-3f93324c0d1e", "Kampala Response Team")
	KampalaICLATeam     = team("a43f84d5-3f8a-48c4-a896-5fb0fcd3e42b", "Kampala ICLA Team")
	NairobiResponseTeam = team("814fc372-08a6-4e6b-809b-30ebb51cb268", "Nairobi Response Team")
	NairobiICLATeam     = team("80606eb4-b53a-4fda-be12-e9806e11d44a", "Nairobi ICLA Team")

	// Team Orgs
	KampalaResponseTeamOrg = relationship(iam.TeamOrganizationRelationshipType.ID, KampalaResponseTeam.ID, NRC.ID)
	KampalaICLATeamOrg     = relationship(iam.TeamOrganizationRelationshipType.ID, KampalaICLATeam.ID, NRC.ID)
	NairobiResponseTeamOrg = relationship(iam.TeamOrganizationRelationshipType.ID, NairobiResponseTeam.ID, NRC.ID)
	NairobiICLATeamOrg     = relationship(iam.TeamOrganizationRelationshipType.ID, NairobiICLATeam.ID, NRC.ID)

	// Case Types
	GenderViolence     = caseType("920ca64f-66d0-4f00-af05-d8a50ce354e6", "Gender Violence", iam.IndividualPartyType.ID, KampalaResponseTeam.ID)
	Childcare          = caseType("039caa4e-1f93-40a9-a6ff-b5bf801dfd41", "Childcare", iam.IndividualPartyType.ID, NairobiICLATeam.ID)
	HousingRights      = caseType("0d91d175-673a-4c10-bf87-3b940585e4ac", "Housing Rights", iam.IndividualPartyType.ID, KampalaICLATeam.ID)
	FinancialAssistInd = caseType("785036d9-41ee-4413-987e-13d4af761737", "Financial Assistance", iam.IndividualPartyType.ID, KampalaResponseTeam.ID)
	FinancialAssistHH  = caseType("52a24b6b-ad10-4297-b030-263bbdcd5420", "Financial Assistance", iam.IndividualPartyType.ID, NairobiResponseTeam.ID)

	// Individuals
	JohnDoe     = individual("c529d679-3bb6-4a20-8f06-c096f4d9adc1", "John", "Doe")
	MaryPoppins = individual("bbf539fd-ebaa-4438-ae4f-8aca8b327f42", "Mary", "Poppins")
	BoDiddley   = individual("26335292-c839-48b6-8ad5-81271ee51e7b", "Bo", "Diddley")
	Howell      = individual("066a0268-fdc6-495a-9e4b-d60cfae2d81a", "Howell", "Jorg")
	Birdie      = individual("ac9015ac-686f-4719-9c3d-bf3d1cae00ea", "Birdie", "Tifawt")
	Charis      = individual("ce7ae69c-9f6a-413b-96bf-5808d0da92cd", "Charis", "Timothy")
	Danyial     = individual("8da05c97-12c2-4b43-b022-dc79be7dc3a0", "Danyial", "Hrodebert")
	Devi        = individual("6414895a-ce60-4647-b491-baeb54a76f26", "Devi", "Malvina")
	Levan       = individual("5d4c4302-ad8e-45ab-bd4c-e1ac25ae972e", "Levan", "Elija")
	Lisbeth     = individual("a5d4dab0-90d3-474d-afe6-46d04ca3caba", "Lisbeth", "Furkan")
	Liadan      = individual("c9ce906d-87ba-4123-bb74-7a73664e6778", "Liadan", "Jordaan")
	Muhammad    = individual("818206ea-0b5e-4ed9-b47e-db31566d10c0", "Muhammad", "Annemarie")
	Dardanos    = individual("7921756a-8759-4589-8a83-ad98f8aa22c7", "Dardanos", "Rilla")
	Jana        = individual("c7ca3a4d-0e96-4e5c-8c32-6750d0312706", "Jana", "Nurul")
	Simeon      = individual("78663ffb-dbaa-4362-83b6-7319d6469caa", "Simeon", "Tumelo")
	Sayen       = individual("29a20d76-dd37-471f-b9ec-9ab08f61d1ed", "Sayen", "Gezabele")
	Veniaminu   = individual("051a46b2-1ef4-4c86-bd2f-9306daedec7e", "Veniaminu", "Ye-Jun")
	Loan        = individual("f2a5d586-6865-40ea-a3db-7c729516b32b", "Loan", "Daniel")
	Reece       = individual("bdeb7e66-9129-467e-abc0-51ab2df7f222", "Reece", "Hyakinthos")
	Svetlana    = individual("afdd8b5c-b9b4-41e1-a015-7e0beb33f10b", "Svetlana", "Cerdic")
	Kyleigh     = individual("12d6a293-d923-47c6-9bc1-441934bb79c5", "Kyleigh", "Jayma")
	Hermina     = individual("dafee423-49c0-4fbf-b2f9-a42276c0cfce", "Hermina", "Magnus")
	Leela       = individual("65410229-ad41-4c17-88f2-13e9a56a0fe8", "Leela", "Cynebald")
	Jovan       = individual("bf22e83b-cfef-4c8a-b74e-f0cef6b27147", "Jovan", "Lynette")
	Bor         = individual("e350e394-091f-469c-a217-488b27b113a3", "Bor", "Lora")
	Aldwin      = individual("fdb6a682-8eb6-4565-879b-835a76384fe0", "Aldwin", "Colin")
	Trophimos   = individual("bb800fe3-85a7-4c90-b8f2-cd0354825f56", "Trophimos", "Wiebke")

	// Staff
	HowellStaff    = staff("8a2497f9-1608-405a-adc9-6c29147f6a06", Howell)
	BirdieStaff    = staff("fec7ea88-c77c-4311-8824-8e6aa65c2702", Birdie)
	CharisStaff    = staff("3cf8f80b-31b3-4744-8743-abfabbe6afc7", Charis)
	DanyialStaff   = staff("87d3cf68-5141-4838-bc8a-145deb0587ca", Danyial)
	DeviStaff      = staff("0573c9a9-94ed-4353-a8bd-5c8422d8a66c", Devi)
	LevanStaff     = staff("a6f32747-558e-4b2a-a767-d6313ea83a17", Levan)
	LisbethStaff   = staff("3fce53f1-2583-438f-b5e6-e70efcaaa935", Lisbeth)
	BorStaff       = staff("ac7f8a3c-0f7e-49b3-b11d-5135cf3a3af0", Bor)
	LiadanStaff    = staff("168f3397-6f1c-4607-b0a2-e130f8dd9cc4", Liadan)
	MuhammadStaff  = staff("e559a5d1-7293-4c2f-8ac2-c4aa4ee12571", Muhammad)
	DardanosStaff  = staff("ef162397-b32d-4e4f-8418-606265b937b1", Dardanos)
	JanaStaff      = staff("58d10f3a-4e1d-4a2b-8b19-109d5fe540c7", Jana)
	SimeonStaff    = staff("79fcf4dd-2603-4d7f-891d-604529ecfe3c", Simeon)
	SayenStaff     = staff("540fd5da-e152-4401-b30a-27c5a661a4f7", Sayen)
	VeniaminuStaff = staff("a2d7047b-b126-4213-a0a5-7ef4ecb09515", Veniaminu)
	LoanStaff      = staff("304ae2f9-3cf7-4aed-a26d-24a400788bc1", Loan)
	ReeceStaff     = staff("b9c847b2-fefa-4223-8fdb-a048ef3153fc", Reece)
	SvetlanaStaff  = staff("3b88b8a2-5cfd-403e-bc5f-724479418a94", Svetlana)
	KyleighStaff   = staff("13fba20e-07dc-49a8-bd41-bdcc796d53e0", Kyleigh)
	HerminaStaff   = staff("56b35e75-5703-46b6-9232-51570f0c73bc", Hermina)
	LeelaStaff     = staff("b0ca2145-c884-4aaf-88f0-39f2c8e9acdf", Leela)
	JovanStaff     = staff("265ffaf6-fe44-4290-9755-17894b3c2e5f", Jovan)
	AldwinStaff    = staff("31c1cdfc-cda2-4012-bd21-9ffe0939cda5", Aldwin)
	TrophimosStaff = staff("857e197c-19f8-4546-a2cb-769d42a0fd55", Trophimos)

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
	DomesticAbuse    = kase("dba43642-8093-4685-a197-f8848d4cbaaa", GenderViolence.ID, MaryPoppins.ID, KampalaResponseTeam.ID, "Domestic abuse", false)
	MonthlyAllowance = kase("47499762-c189-4a74-9156-7969f899073b", FinancialAssistInd.ID, JohnDoe.ID, KampalaResponseTeam.ID, "Monthly allowance", false)
	ChildCare        = kase("8fb5f755-85eb-4d91-97a9-fdf86c01df25", Childcare.ID, BoDiddley.ID, KampalaResponseTeam.ID, "Monthly stipend for Bo Diddley's child", true)
)

func Seed(ctx context.Context, databaseName string, mongoClient *mongo.Client) error {

	for _, obj := range individuals {
		if err := seedMongo(ctx, mongoClient, databaseName, "parties", bson.M{"id": obj.ID}, obj.Party); err != nil {
			return err
		}
	}

	for _, obj := range orgs {
		if err := seedMongo(ctx, mongoClient, databaseName, "organizations", bson.M{"id": obj.ID}, obj.Party); err != nil {
			return err
		}
	}

	for _, obj := range teams {
		if err := seedMongo(ctx, mongoClient, databaseName, "parties", bson.M{"id": obj.ID}, iam.MapTeamToParty(&obj)); err != nil {
			return err
		}
	}

	for _, obj := range relationships {
		if err := seedMongo(ctx, mongoClient, databaseName, "relationships", bson.M{"id": obj.ID}, obj); err != nil {
			return err
		}
	}

	for _, obj := range memberships {
		if err := seedMongo(ctx, mongoClient, databaseName, "relationships", bson.M{"id": obj.ID}, iam.MapMembershipToRelationship(&obj)); err != nil {
			return err
		}
	}

	for _, obj := range caseTypes {
		if err := seedMongo(ctx, mongoClient, databaseName, "caseTypes", bson.M{"id": obj.ID}, obj); err != nil {
			return err
		}
	}

	for _, obj := range cases {
		if err := seedMongo(ctx, mongoClient, databaseName, "cases", bson.M{"id": obj.ID}, obj); err != nil {
			return err
		}
	}

	return nil
}

func caseType(id, name, partyTypeID, teamID string) cms.CaseType {
	ct := cms.CaseType{
		ID:          id,
		Name:        name,
		PartyTypeID: partyTypeID,
		TeamID:      teamID,
	}
	caseTypes = append(caseTypes, ct)
	return ct
}

func relationship(relationshipTypeID, firstPartyID, secondPartyID string) iam.Relationship {
	r := iam.Relationship{
		RelationshipTypeID: relationshipTypeID,
		FirstPartyID:       firstPartyID,
		SecondPartyID:      secondPartyID,
	}
	relationships = append(relationships, r)
	return r
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

func staff(id string, individual iam.Individual) iam.Staff {
	s := iam.Staff{
		ID:             id,
		OrganizationID: NRC.ID,
		IndividualID:   individual.ID,
	}
	staffers = append(staffers, s)
	return s
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

func kase(id, caseTypeID, partyID, teamID, description string, done bool) cms.Case {
	k := cms.Case{
		ID:          id,
		CaseTypeID:  caseTypeID,
		PartyID:     partyID,
		TeamID:      teamID,
		Description: description,
		Done:        done,
	}
	cases = append(cases, k)
	return k
}

func seedMongo(ctx context.Context, mongoClient *mongo.Client, databaseName, collectionName string, filter interface{}, document interface{}) error {
	logrus.Infof("seeding collection %s.%s with object: %#v", databaseName, collectionName, document)
	collection := mongoClient.Database(databaseName).Collection(collectionName)
	if _, err := collection.InsertOne(ctx, document); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
		if _, err := collection.ReplaceOne(ctx, filter, document); err != nil {
			return err
		}
	}
	return nil
}
