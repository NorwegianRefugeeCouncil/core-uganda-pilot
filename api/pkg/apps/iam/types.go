package iam

import (
	"github.com/nrc-no/core/pkg/generic/pagination"
	"strings"
	"time"
)

// Attribute represents an attribute that can be attached to a Party
type Attribute struct {
	// ID is the unique ID of an attribute
	ID string `json:"id" bson:"id"`

	// Name is the name of an attribute
	Name string `json:"name" bson:"name"`

	// PartyTypeIDs represents the type of Party that can have this attribute
	PartyTypeIDs []string `json:"partyTypeIds" bson:"partyTypeIds"`

	// IsPersonallyIdentifiableInfo represents wether or not the Attribute is considered
	// as personally identifiable information
	IsPersonallyIdentifiableInfo bool `json:"isPii" bson:"isPii"`

	// Translations represent the localized descriptions of the Attribute
	Translations []AttributeTranslation `json:"translations" bson:"translations"`
}

// AttributeTranslation represents a localized description of an Attribute
type AttributeTranslation struct {

	// Locale represents the locale of the AttributeTranslation
	Locale string `json:"locale" bson:"locale"`

	// LongFormulation is a longer formulation for display purposes
	LongFormulation string `json:"longFormulation" bson:"longFormulation"`

	// ShortFormulation is a short formulation for display purposes
	ShortFormulation string `json:"shortFormulation" bson:"shortFormulation"`
}

// AttributeList represents a list of Attribute
type AttributeList struct {
	Items []*Attribute `json:"items" bson:"items"`
}

// FindByID finds an Attribute by ID
func (l *AttributeList) FindByID(id string) *Attribute {
	for _, item := range l.Items {
		if item.ID == id {
			return item
		}
	}
	return nil
}

// TODO COR-158
// PartyAttributes contains the Attribute values of a Party
type PartyAttributes map[string][]string

// Get returns the first value of a Party Attribute
func (a PartyAttributes) Get(key string) string {
	if values, ok := a[key]; ok {
		if len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

// Set sets the value of an Attribute
func (a PartyAttributes) Set(key, value string) {
	a[key] = []string{value}
}

// Add adds an Attribute value
func (a PartyAttributes) Add(key, value string) {
	a[key] = append(a[key], value)
}

// Party represents a physical or moral person
type Party struct {

	// ID is the unique ID of a Party
	ID string `json:"id" bson:"id"`

	// PartyTypeIDs represent the different PartyType that this Party has
	PartyTypeIDs []string `json:"partyTypeIds" bson:"partyTypeIds"`

	// Attributes represent the Attribute values
	Attributes PartyAttributes `json:"attributes" bson:"attributes"`
}

// HasPartyType checks if the Party has the given PartyType
func (p *Party) HasPartyType(partyType string) bool {
	for _, p := range p.PartyTypeIDs {
		if p == partyType {
			return true
		}
	}
	return false
}

// AddPartyType adds the given PartyType to the list of PartyTypes
func (p *Party) AddPartyType(partyType string) {
	if p.HasPartyType(partyType) {
		return
	}
	p.PartyTypeIDs = append(p.PartyTypeIDs, partyType)
}

func (p *Party) String() string {

	// Staff
	if p.HasPartyType(StaffPartyType.ID) {
		return p.Attributes.Get(FirstNameAttribute.ID) +
			" " +
			p.Attributes.Get(LastNameAttribute.ID)
	}

	// Individual
	if p.HasPartyType(IndividualPartyType.ID) {
		return p.Attributes.Get(FirstNameAttribute.ID) +
			" " +
			p.Attributes.Get(LastNameAttribute.ID)
	}

	// Team
	if p.HasPartyType(TeamPartyType.ID) {
		return p.Attributes.Get(TeamNameAttribute.ID)
	}

	// Country
	if p.HasPartyType(CountryPartyType.ID) {
		return p.Attributes.Get(CountryNameAttribute.ID)
	}
	// Default
	return p.ID
}

// GetAttribute returns the value of an Attribute
func (p *Party) GetAttribute(name string) []string {
	if v, ok := p.Attributes[name]; ok {
		return v
	}
	return nil
}

// Get returns the first value of an Attribute
func (p *Party) Get(name string) string {
	if v, ok := p.Attributes[name]; ok {
		if len(v) > 0 {
			return v[0]
		}
	}
	return ""
}

// HasAttribute checks if the Party has the given Attribute
func (p *Party) HasAttribute(name string) bool {
	_, ok := p.Attributes[name]
	return ok
}

// FindAttributeValue returns the value of an Attribute
func (p *Party) FindAttributeValue(name string) interface{} {
	if v, ok := p.Attributes[name]; ok {
		return v
	}
	return nil
}

// PartyList is a list of Party
type PartyList struct {
	Items []*Party `json:"items" bson:"items"`
}

// FindByID finds a Party by ID
func (pl *PartyList) FindByID(id string) *Party {
	for _, v := range pl.Items {
		if id == v.ID {
			return v
		}
	}
	return nil
}

// PartyType represents a type of Party
type PartyType struct {
	// ID represents the unique ID of this PartyType
	ID string `json:"id" bson:"id"`

	// Name is the name of the PartyType
	Name string `json:"name" bson:"name"`

	// IsBuiltIn indicates whether this is a system-managed PartyType
	IsBuiltIn bool `json:"isBuiltIn" bson:"isBuiltIn"`
}

// PartyTypeList is a list of PartyType
type PartyTypeList struct {

	// Items the list of PartyType
	Items []*PartyType
}

// FindByID finds a PartyType by ID
func (p *PartyTypeList) FindByID(id string) *PartyType {
	for _, item := range p.Items {
		if item.ID == id {
			return item
		}
	}
	return nil
}

// Relationship represents a relationship between two Party
type Relationship struct {

	// ID the unique ID of the Relationship
	ID string `json:"id" bson:"id"`

	// RelationshipTypeID the RelationshipType of that Relationship
	RelationshipTypeID string `json:"relationshipTypeId" bson:"relationshipTypeId"`

	// FirstPartyID the first Party part of that Relationship
	FirstPartyID string `json:"firstParty" bson:"firstParty"`

	// SecondPartyID the second Party of that relationship
	SecondPartyID string `json:"secondParty" bson:"secondParty"`
}

// RelationshipList represents a list of Relationship
type RelationshipList struct {

	// Items the list of Relationship
	Items []*Relationship `json:"items" bson:"items"`
}

// RelationshipTypeRule represents a rule that restricts which kind of
// Party can be on either side of the Relationship
type RelationshipTypeRule struct {
	PartyTypeRule *PartyTypeRule `json:",inline" bson:",inline"`
}

// Mirror returns the mirror of a Relationship (inverses the sides)
func (r RelationshipTypeRule) Mirror() RelationshipTypeRule {
	ret := RelationshipTypeRule{}
	if r.PartyTypeRule != nil {
		rev := r.PartyTypeRule.Mirror()
		ret.PartyTypeRule = &rev
	}
	return ret
}

type PartyTypeRule struct {
	FirstPartyTypeID  string `json:"firstPartyTypeId" bson:"firstPartyTypeId"`
	SecondPartyTypeID string `json:"secondPartyTypeId" bson:"secondPartyTypeId"`
}

// Mirror returns the mirror image of a PartyTypeRule
func (p PartyTypeRule) Mirror() PartyTypeRule {
	return PartyTypeRule{
		FirstPartyTypeID:  p.SecondPartyTypeID,
		SecondPartyTypeID: p.FirstPartyTypeID,
	}
}

// RelationshipType represents a type of Relationship (marriage, employment, etc.)
type RelationshipType struct {
	// ID is the ID of the Relationship
	ID string `json:"id" bson:"id"`
	// IsDirectional indicates that the relationship sides are different.
	// NonDirectional indicate that a relationship is semantically equal when it is inversed
	IsDirectional bool `json:"isDirectional" bson:"isDirectional"`
	// The name of the Relationship type
	Name string `json:"name" bson:"name"`
	// The role of the first Party in the relationship
	FirstPartyRole string `json:"firstPartyRole" bson:"firstPartyRole"`
	// The role of the second Party in the relationship
	SecondPartyRole string `json:"secondPartyRole" bson:"secondPartyRole"`
	// The relationship rules
	Rules []RelationshipTypeRule `json:"rules"`
}

// RelationshipTypeList represents a list of Relationship
type RelationshipTypeList struct {
	Items []*RelationshipType `json:"items" bson:"items"`
}

// Mirror returns the mirror image of a RelationshipType
func (r *RelationshipType) Mirror() *RelationshipType {
	rules := r.Rules
	for i, rule := range rules {
		rules[i] = rule.Mirror()
	}
	return &RelationshipType{
		ID:              r.ID,
		IsDirectional:   r.IsDirectional,
		Name:            r.Name,
		FirstPartyRole:  r.SecondPartyRole,
		SecondPartyRole: r.FirstPartyRole,
		Rules:           rules,
	}
}

type Individual struct {
	*Party `json:",inline" bson:",inline"`
}

type IndividualList struct {
	Items    []*Individual         `json:"items" bson:"items"`
	Metadata pagination.Pagination `json:"metadata"`
}

func NewIndividual(ID string) *Individual {
	return &Individual{
		Party: &Party{
			ID: ID,
			PartyTypeIDs: []string{
				IndividualPartyType.ID,
			},
			Attributes: map[string][]string{},
		},
	}
}

func (b *Individual) FindAge() *int {

	birthDate := b.FindAttributeValue(BirthDateAttribute.ID)
	if birthDate == nil {
		return nil
	}

	birthDateStrs, ok := birthDate.([]string)
	if !ok {
		return nil
	}

	if len(birthDateStrs) == 0 {
		return nil
	}

	birthDateStr := birthDateStrs[0]
	parsedDate, err := time.Parse("2006-01-02", birthDateStr)
	if err != nil {
		return nil
	}

	diff := time.Now().UTC().Sub(parsedDate)
	years := diff.Hours() / 24 / 365
	rounded := int(years)

	return &rounded

}

func (b *Individual) String() string {

	firstNames, hasFirstNames := b.Attributes[FirstNameAttribute.ID]
	lastNames, hasLastNames := b.Attributes[LastNameAttribute.ID]

	if hasFirstNames && hasLastNames {
		return strings.ToUpper(lastNames[0]) + ", " + firstNames[0]
	}

	return b.ID
}

type Team struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

func (t *Team) String() string {
	return t.Name
}

type TeamList struct {
	Items []*Team `json:"items"`
}

func (l *TeamList) FindByID(id string) *Team {
	for _, team := range l.Items {
		if team.ID == id {
			return team
		}
	}
	return nil
}

//Country
type Country struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

func (c *Country) String() string {
	return c.Name
}

type CountryList struct {
	Items []*Country `json:"items"`
}

func (l *CountryList) FindByID(id string) *Country {
	for _, country := range l.Items {
		if country.ID == id {
			return country
		}
	}
	return nil
}

// Staff is a relationship between an organization and an individual
// that represents that the individual is working for that organization
type Staff struct {
	ID           string `json:"id"`
	IndividualID string `json:"individualId"`
}

// StaffList is a list of Staff
type StaffList struct {
	Items []*Staff `json:"items"`
}

type Membership struct {
	ID           string `json:"id"`
	TeamID       string `json:"teamId"`
	IndividualID string `json:"individualId"`
}

type MembershipList struct {
	Items []*Membership `json:"items"`
}

type Nationality struct {
	ID        string `json:"id"`
	CountryID string `json:"CountryId"`
	TeamID    string `json:"teamId"`
}

type NationalityList struct {
	Items []*Nationality `json:"items"`
}
