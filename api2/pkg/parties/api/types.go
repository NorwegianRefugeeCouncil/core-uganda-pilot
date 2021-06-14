package api

import (
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/expressions"
	"strings"
	"time"
)

type SubjectType string

const (
	BeneficiaryType SubjectType = "Beneficiary"
	HouseholdType   SubjectType = "Household"
)

type Attribute struct {
	ID                           string                 `json:"id" bson:"id"`
	Name                         string                 `json:"name" bson:"name"`
	ValueType                    expressions.ValueType  `json:"type" bson:"type"`
	SubjectType                  SubjectType            `json:"subjectType" bson:"subjectType"`
	IsPersonallyIdentifiableInfo bool                   `json:"isPii" bson:"isPii"`
	Translations                 []AttributeTranslation `json:"translations" bson:"translations"`
}

type AttributeTranslation struct {
	Locale           string `json:"locale" bson:"locale"`
	LongFormulation  string `json:"longFormulation" bson:"longFormulation"`
	ShortFormulation string `json:"shortFormulation" bson:"shortFormulation"`
}

type RelationshipTypeRule struct {
	PartyTypeRule
}

type PartyTypeRule struct {
	FirstPartyType  string `json:"firstPartyType"`
	SecondPartyType string `json:"secondPartyType"`
}

type RelationshipType struct {
	ID              string                 `json:"id" bson:"id"`
	Name            string                 `json:"name" bson:"name"`
	FirstPartyRole  string                 `json:"firstPartyRole" bson:"firstPartyRole"`
	SecondPartyRole string                 `json:"secondPartyRole" bson:"secondPartyRole"`
	Rules           []RelationshipTypeRule `json:"rules"`
}

type RelationshipTypeList struct {
	Items []*RelationshipType
}

type Party struct {
	ID         string   `json:"id" bson:"id"`
	PartyTypes []string `json:"partyTypes" bson:"partyTypes"`
	Attributes map[string]interface{}
}

type PartyList struct {
	Items []*Party `json:"items" bson:"items"`
}

type PartyType struct {
	ID        string `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
	IsBuiltIn bool   `json:"isBuiltIn" bson:"isBuiltIn"`
}

type PartyTypeSchema struct {
	ID    string                `json:"id" bson:"id"`
	Name  string                `json:"name" bson:"name"`
	Nodes []PartyTypeSchemaNode `json:"nodes" bson:"nodes"`
}

type PartyTypeSchemaList struct {
	Items []*PartyTypeSchema `json:"items" bson:"items"`
}

type PartyTypeSchemaNode struct {
	ID          string                `json:"id" bson:"id"`
	PartyTypeID string                `json:"partyTypeID" bson:"partyTypeID"`
	Children    []PartyTypeSchemaNode `json:"children" bson:"children"`
}

type PartyTypeList struct {
	Items []*PartyType
}

type Relationship struct {
	ID                  string     `json:"id" bson:"id"`
	RelationshipTypeID  string     `json:"relationshipTypeId" bson:"relationshipTypeId"`
	FirstParty          string     `json:"firstParty" bson:"firstParty"`
	SecondParty         string     `json:"secondParty" bson:"secondParty"`
	StartOfRelationship time.Time  `json:"startOfRelationship" bson:"startOfRelationship"`
	EndOfRelationship   *time.Time `json:"endOfRelationship" bson:"endOfRelationship"`
}

type RelationshipList struct {
	Items []*Relationship
}

type AttributeList struct {
	Items []*Attribute
}

type AttributeValue struct {
	Attribute
	Value interface{}
}

type Beneficiary struct {
	ID         string `json:"id" bson:"id"`
	Attributes map[string]*AttributeValue
}

func (b *Beneficiary) GetAttribute(name string) *AttributeValue {
	if v, ok := b.Attributes[name]; ok {
		return v
	}
	return nil
}

func (b *Beneficiary) HasAttribute(name string) bool {
	_, ok := b.Attributes[name]
	return ok
}

func (b *Beneficiary) FindAttributeValue(name string) interface{} {
	if v, ok := b.Attributes[name]; ok {
		return v.Value
	}
	return nil
}

func (b *Beneficiary) FindAge() *int {

	birthDate := b.FindAttributeValue("birthDate")
	if birthDate == nil {
		return nil
	}

	parsedDate, err := time.Parse("2006-01-02", birthDate.(string))
	if err != nil {
		return nil
	}

	diff := time.Now().UTC().Sub(parsedDate)
	years := diff.Hours() / 24 / 365
	rounded := int(years)

	return &rounded

}

func (b *Beneficiary) String() string {
	hasFirstName := b.HasAttribute("firstName")
	hasLastName := b.HasAttribute("lastName")

	if hasFirstName && hasLastName {
		return fmt.Sprintf("%s, %s",
			strings.ToUpper(b.GetAttribute("lastName").Value.(string)),
			b.GetAttribute("firstName").Value.(string),
		)
	}
	if hasFirstName {
		return b.GetAttribute("firstName").Value.(string)
	}
	if hasLastName {
		return strings.ToUpper(b.GetAttribute("lastName").Value.(string))
	}
	return b.ID
}

type BeneficiaryList struct {
	Items []*Beneficiary
}
