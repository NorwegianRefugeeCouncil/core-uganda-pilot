package main

import (
	"github.com/nrc-no/core/pkg/cms"
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/iam"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
	"os"
)


func main() {
	outputFile := os.Args[1]
	if len(outputFile) == 0 {
		println("Usage: ./gotypes2ts FILE\nThe program will generate types in the given filepath")
	}
	converter := typescriptify.New()
	// export interfaces and not classes
	converter.CreateInterface = true
	for _, typ := range types {
		converter = converter.Add(typ)
	}
	err := converter.WithBackupDir("").ConvertToFile(outputFile)
	if err != nil {
		panic(err)
	}
}

var	types = []interface{}{
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
	iam.IdentificationDocument{},
	iam.IdentificationDocumentList{},
	iam.IdentificationDocumentListOptions{},
	iam.IdentificationDocumentType{},
	iam.IdentificationDocumentTypeList{},
	iam.IdentificationDocumentTypeListOptions{},

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
