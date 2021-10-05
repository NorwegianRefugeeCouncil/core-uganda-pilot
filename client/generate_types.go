package main

import (
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/form"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func GenerateTypescriptTypes() error {
	types := []interface{}{
		// iam types
		iam.PartyAttributeDefinition{},
		iam.PartyAttributeDefinitionList{},
		iam.PartyAttributeDefinitionListOptions{},
		iam.AttributeMap{},
		iam.Party{},
		iam.PartyList{},
		iam.PartyType{},
		iam.PartyTypeList{},
		iam.Relationship{},
		iam.RelationshipList{},
		iam.RelationshipTypeRule{},
		iam.PartyTypeRule{},
		iam.RelationshipType{},
		iam.RelationshipTypeList{},
		iam.Individual{},
		iam.IndividualList{},
		iam.Team{},
		iam.TeamList{},
		iam.Country{},
		iam.CountryList{},
		iam.Staff{},
		iam.StaffList{},
		iam.Membership{},
		iam.MembershipList{},
		iam.Nationality{},
		iam.NationalityList{},
		iam.PartyListOptions{},
		iam.PartySearchOptions{},
		iam.PartyTypeListOptions{},
		iam.RelationshipListOptions{},
		iam.RelationshipTypeListOptions{},
		iam.TeamListOptions{},
		iam.CountryListOptions{},
		iam.StaffListOptions{},
		iam.MembershipListOptions{},
		iam.NationalityListOptions{},
		iam.IndividualListOptions{},

		// cms types
		cms.Case{},
		cms.CaseList{},
		cms.CaseType{},
		cms.CaseTypeList{},
		cms.Comment{},
		cms.CommentList{},
		cms.CaseListOptions{},
		cms.CaseTypeListOptions{},
		cms.CommentListOptions{},

		// other types
		form.Form{},
	}
	converter := typescriptify.New()
	converter.CreateInterface = true
	for _, typ := range types {
		converter = converter.Add(typ)
	}
	err := converter.ConvertToFile("types/models.ts")
	return err
}

func main() {
	GenerateTypescriptTypes()
}
