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
	ID                           string                 `json:"id"`
	Name                         string                 `json:"name"`
	ValueType                    expressions.ValueType  `json:"type"`
	SubjectType                  SubjectType            `json:"subjectType"`
	IsPersonallyIdentifiableInfo bool                   `json:"isPii"`
	Translations                 []AttributeTranslation `json:"translations"`
}

type AttributeTranslation struct {
	Locale           string `json:"locale"`
	LongFormulation  string `json:"longFormulation"`
	ShortFormulation string `json:"shortFormulation"`
}

type RelationshipType struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	FirstPartyRole  string `json:"firstPartyRole"`
	SecondPartyRole string `json:"secondPartyRole"`
}

type RelationshipTypeList struct {
	Items []*RelationshipType
}

type Relationship struct {
	ID                  string     `json:"id"`
	FirstParty          string     `json:"firstParty"`
	SecondParty         string     `json:"secondParty"`
	StartOfRelationship time.Time  `json:"startOfRelationship"`
	EndOfRelationship   *time.Time `json:"endOfRelationship"`
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
	ID         string `json:"id"`
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
