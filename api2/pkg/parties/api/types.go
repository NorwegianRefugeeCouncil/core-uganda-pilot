package api

import (
	"github.com/nrc-no/core-kafka/pkg/expressions"
	"time"
)

type Attribute struct {
	ID                           string                 `json:"id" bson:"id"`
	Name                         string                 `json:"name" bson:"name"`
	ValueType                    expressions.ValueType  `json:"type" bson:"type"`
	PartyTypes                   []string               `json:"partyTypes" bson:"partyTypes"`
	IsPersonallyIdentifiableInfo bool                   `json:"isPii" bson:"isPii"`
	Translations                 []AttributeTranslation `json:"translations" bson:"translations"`
}

type AttributeTranslation struct {
	Locale           string `json:"locale" bson:"locale"`
	LongFormulation  string `json:"longFormulation" bson:"longFormulation"`
	ShortFormulation string `json:"shortFormulation" bson:"shortFormulation"`
}

type RelationshipTypeRule struct {
	PartyTypeRule `json:",inline" bson:",inline"`
}

type PartyTypeRule struct {
	FirstPartyType  string `json:"firstPartyType" bson:"firstPartyType"`
	SecondPartyType string `json:"secondPartyType" bson:"secondPartyType"`
}

type RelationshipType struct {
	ID              string                 `json:"id" bson:"id"`
	Name            string                 `json:"name" bson:"name"`
	FirstPartyRole  string                 `json:"firstPartyRole" bson:"firstPartyRole"`
	SecondPartyRole string                 `json:"secondPartyRole" bson:"secondPartyRole"`
	Rules           []RelationshipTypeRule `json:"rules"`
}

type RelationshipTypeList struct {
	Items []*RelationshipType `json:"items" bson:"items"`
}

type Party struct {
	ID         string              `json:"id" bson:"id"`
	PartyTypes []string            `json:"partyTypes" bson:"partyTypes"`
	Attributes map[string][]string `json:"attributes" bson:"attributes"`
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
	Items []*PartyType `json:"items" bson:"items"`
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
	Items []*Relationship `json:"items" bson:"items"`
}

type AttributeList struct {
	Items []*Attribute `json:"items" bson:"items"`
}

type AttributeValue struct {
	Attribute
	Value interface{}
}

func (b *Party) GetAttribute(name string) []string {
	if v, ok := b.Attributes[name]; ok {
		return v
	}
	return nil
}

func (b *Party) HasAttribute(name string) bool {
	_, ok := b.Attributes[name]
	return ok
}

func (b *Party) FindAttributeValue(name string) interface{} {
	if v, ok := b.Attributes[name]; ok {
		return v
	}
	return nil
}
