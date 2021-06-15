package api

import (
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"strings"
	"time"
)

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

type PartyAttributes map[string][]string

func (a PartyAttributes) Get(key string) string {
	if values, ok := a[key]; ok {
		if len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

func (a PartyAttributes) Set(key, value string) {
	a[key] = []string{value}
}

func (a PartyAttributes) Add(key, value string) {
	a[key] = append(a[key], value)
}

type Party struct {
	ID         string          `json:"id" bson:"id"`
	PartyTypes []string        `json:"partyTypes" bson:"partyTypes"`
	Attributes PartyAttributes `json:"attributes" bson:"attributes"`
}

func (p *Party) HasPartyType(partyType string) bool {
	for _, p := range p.PartyTypes {
		if p == partyType {
			return true
		}
	}
	return false
}

func (p *Party) AddPartyType(partyType string) {
	if p.HasPartyType(partyType) {
		return
	}
	p.PartyTypes = append(p.PartyTypes, partyType)
}

func (p *Party) String() string {
	if p.HasPartyType(partytypes.IndividualPartyType.ID) {
		firstName := p.Attributes.Get(attributes.FirstNameAttribute.ID)
		lastName := p.Attributes.Get(attributes.LastNameAttribute.ID)
		return strings.ToUpper(lastName) + ", " + firstName
	}
	return ""
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

func (p *Party) GetAttribute(name string) []string {
	if v, ok := p.Attributes[name]; ok {
		return v
	}
	return nil
}

func (p *Party) HasAttribute(name string) bool {
	_, ok := p.Attributes[name]
	return ok
}

func (p *Party) FindAttributeValue(name string) interface{} {
	if v, ok := p.Attributes[name]; ok {
		return v
	}
	return nil
}
