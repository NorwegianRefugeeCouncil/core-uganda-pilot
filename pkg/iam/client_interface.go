package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/rest"
	"net/url"
	"strconv"
	"strings"
)

type (
	Interface interface {
		Parties() PartyClient
		PartyTypes() PartyTypeClient
		Relationships() RelationshipClient
		RelationshipTypes() RelationshipTypeClient
		PartyAttributeDefinitions() PartyAttributeDefinitionClient
		Teams() TeamClient
		Memberships() MembershipClient
		Individuals() IndividualClient
		Countries() CountryClient
		Nationalities() NationalityClient
		IdentificationDocuments() IdentificationDocumentClient
		IdentificationDocumentTypes() IdentificationDocumentTypeClient
	}
)

type PartyListOptions struct {
	PartyTypeID string
	SearchParam string
	Attributes  map[string]string
}

func (a *PartyListOptions) MarshalQueryParameters() (url.Values, error) {
	ret := url.Values{}
	if len(a.PartyTypeID) > 0 {
		ret.Set("partyTypeId", a.PartyTypeID)
	}
	if len(a.SearchParam) > 0 {
		ret.Set("searchParam", a.SearchParam)
	}
	if a.Attributes != nil {
		for key, value := range a.Attributes {
			if len(value) == 0 {
				continue
			}
			ret.Set("attributes["+key+"]", value)
		}
	}
	return ret, nil
}

func (a *PartyListOptions) UnmarshalQueryParameters(values url.Values) error {
	a.PartyTypeID = values.Get("partyTypeId")
	a.SearchParam = values.Get("searchParam")
	for key, values := range values {
		if strings.HasPrefix(key, "attributes[") && strings.HasSuffix(key, "]") {
			if a.Attributes == nil {
				a.Attributes = map[string]string{}
			}
			attrKey := key[11 : len(key)-1]
			a.Attributes[attrKey] = values[0]
		}
	}
	return nil
}

type PartySearchOptions struct {
	PartyIDs     []string          `json:"partyIds"`
	PartyTypeIDs []string          `json:"partyTypeIds"`
	Attributes   map[string]string `json:"attributes"`
	SearchParam  string            `json:"searchParam"`
}

type PartyClient interface {
	Get(ctx context.Context, id string) (*Party, error)
	Create(ctx context.Context, party *Party) (*Party, error)
	Update(ctx context.Context, party *Party) (*Party, error)
	List(ctx context.Context, listOptions PartyListOptions) (*PartyList, error)
	Search(ctx context.Context, listOptions PartySearchOptions) (*PartyList, error)
}

type PartyTypeListOptions struct {
}

func (a *PartyTypeListOptions) MarshalQueryParameters() (url.Values, error) {
	return url.Values{}, nil
}

func (a *PartyTypeListOptions) UnmarshalQueryParameters(values url.Values) error {
	return nil
}

type PartyTypeClient interface {
	Get(ctx context.Context, id string) (*PartyType, error)
	Create(ctx context.Context, party *PartyType) (*PartyType, error)
	Update(ctx context.Context, party *PartyType) (*PartyType, error)
	List(ctx context.Context, listOptions PartyTypeListOptions) (*PartyTypeList, error)
}

type RelationshipListOptions struct {
	RelationshipTypeID string
	FirstPartyID       string
	SecondPartyID      string
	EitherPartyID      string
}

func (a *RelationshipListOptions) MarshalQueryParameters() (url.Values, error) {
	ret := url.Values{}
	if len(a.RelationshipTypeID) > 0 {
		ret.Set("relationshipTypeId", a.RelationshipTypeID)
	}
	if len(a.FirstPartyID) > 0 {
		ret.Set("firstPartyId", a.FirstPartyID)
	}
	if len(a.SecondPartyID) > 0 {
		ret.Set("secondPartyId", a.SecondPartyID)
	}
	if len(a.EitherPartyID) > 0 {
		ret.Set("eitherPartyId", a.EitherPartyID)
	}
	return ret, nil
}

func (a *RelationshipListOptions) UnmarshalQueryParameters(values url.Values) error {
	a.RelationshipTypeID = values.Get("relationshipTypeId")
	a.FirstPartyID = values.Get("firstPartyId")
	a.SecondPartyID = values.Get("secondPartyId")
	a.EitherPartyID = values.Get("eitherPartyId")
	return nil
}

type RelationshipClient interface {
	Get(ctx context.Context, id string) (*Relationship, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, party *Relationship) (*Relationship, error)
	Update(ctx context.Context, party *Relationship) (*Relationship, error)
	List(ctx context.Context, listOptions RelationshipListOptions) (*RelationshipList, error)
}

type RelationshipTypeListOptions struct {
	PartyTypeID string
}

func (a *RelationshipTypeListOptions) MarshalQueryParameters() (url.Values, error) {
	urlValues := url.Values{}
	if len(a.PartyTypeID) > 0 {
		urlValues.Set("partyTypeId", a.PartyTypeID)
	}
	return urlValues, nil
}

func (a *RelationshipTypeListOptions) UnmarshalQueryParameters(values url.Values) error {
	a.PartyTypeID = values.Get("partyTypeId")
	return nil
}

type RelationshipTypeClient interface {
	Get(ctx context.Context, id string) (*RelationshipType, error)
	Create(ctx context.Context, party *RelationshipType) (*RelationshipType, error)
	Update(ctx context.Context, party *RelationshipType) (*RelationshipType, error)
	List(ctx context.Context, listOptions RelationshipTypeListOptions) (*RelationshipTypeList, error)
}

type PartyAttributeDefinitionListOptions struct {
	PartyTypeIDs []string
	CountryIDs   []string
}

func (a *PartyAttributeDefinitionListOptions) MarshalQueryParameters() (url.Values, error) {
	ret := url.Values{}
	for _, partyTypeID := range a.PartyTypeIDs {
		ret.Add("partyTypeId", partyTypeID)
	}
	for _, countryID := range a.CountryIDs {
		ret.Add("countryId", countryID)
	}
	return ret, nil
}

func (a *PartyAttributeDefinitionListOptions) UnmarshalQueryParameters(values url.Values) error {
	partyTypeIDs := values["partyTypeId"]
	countryIDs := values["countryId"]
	a.PartyTypeIDs = append(a.PartyTypeIDs, partyTypeIDs...)
	a.CountryIDs = append(a.CountryIDs, countryIDs...)
	return nil
}

var _ rest.URLValuer = &PartyAttributeDefinitionListOptions{}

type PartyAttributeDefinitionClient interface {
	Get(ctx context.Context, id string) (*PartyAttributeDefinition, error)
	Create(ctx context.Context, create *PartyAttributeDefinition) (*PartyAttributeDefinition, error)
	Update(ctx context.Context, update *PartyAttributeDefinition) (*PartyAttributeDefinition, error)
	List(ctx context.Context, listOptions PartyAttributeDefinitionListOptions) (*PartyAttributeDefinitionList, error)
}

type IdentificationDocumentTypeClient interface {
	Get(ctx context.Context, id string) (*IdentificationDocumentType, error)
	Create(ctx context.Context, create *IdentificationDocumentType) (*IdentificationDocumentType, error)
	Update(ctx context.Context, update *IdentificationDocumentType) (*IdentificationDocumentType, error)
	List(ctx context.Context, listOptions IdentificationDocumentTypeListOptions) (*IdentificationDocumentTypeList, error)
}

type IdentificationDocumentTypeListOptions struct {
}

func (a *IdentificationDocumentTypeListOptions) MarshalQueryParameters() (url.Values, error) {
	return url.Values{}, nil
}

func (a *IdentificationDocumentTypeListOptions) UnmarshalQueryParameters(values url.Values) error {
	return nil
}

type IdentificationDocumentClient interface {
	Get(ctx context.Context, id string) (*IdentificationDocument, error)
	Create(ctx context.Context, create *IdentificationDocument) (*IdentificationDocument, error)
	Update(ctx context.Context, update *IdentificationDocument) (*IdentificationDocument, error)
	List(ctx context.Context, listOptions IdentificationDocumentListOptions) (*IdentificationDocumentList, error)
	Delete(ctx context.Context, id string) error
}

type IdentificationDocumentListOptions struct {
	PartyIDs []string
}

func (a *IdentificationDocumentListOptions) MarshalQueryParameters() (url.Values, error) {
	ret := url.Values{}
	for _, partyTypeID := range a.PartyIDs {
		ret.Add("partyId", partyTypeID)
	}
	return ret, nil
}

func (a *IdentificationDocumentListOptions) UnmarshalQueryParameters(values url.Values) error {
	partyIDs := values["partyId"]
	a.PartyIDs = append(a.PartyIDs, partyIDs...)
	return nil
}

type TeamListOptions struct {
}

func (a *TeamListOptions) MarshalQueryParameters() (url.Values, error) {
	return url.Values{}, nil
}

func (a *TeamListOptions) UnmarshalQueryParameters(values url.Values) error {
	return nil
}

type TeamClient interface {
	Get(ctx context.Context, id string) (*Team, error)
	Create(ctx context.Context, party *Team) (*Team, error)
	Update(ctx context.Context, party *Team) (*Team, error)
	List(ctx context.Context, listOptions TeamListOptions) (*TeamList, error)
}

// Country

type CountryListOptions struct {
}

func (a *CountryListOptions) MarshalQueryParameters() (url.Values, error) {
	return url.Values{}, nil
}

func (a *CountryListOptions) UnmarshalQueryParameters(values url.Values) error {
	return nil
}

type CountryClient interface {
	Get(ctx context.Context, id string) (*Country, error)
	Create(ctx context.Context, party *Country) (*Country, error)
	Update(ctx context.Context, party *Country) (*Country, error)
	List(ctx context.Context, listOptions CountryListOptions) (*CountryList, error)
}

type StaffListOptions struct {
	IndividualID   string
	OrganizationID string
}

func (a *StaffListOptions) MarshalQueryParameters() (url.Values, error) {
	values := url.Values{}
	if len(a.IndividualID) > 0 {
		values.Set("individualId", a.IndividualID)
	}
	if len(a.OrganizationID) > 0 {
		values.Set("organizationId", a.OrganizationID)
	}
	return values, nil
}

func (a *StaffListOptions) UnmarshalQueryParameters(values url.Values) error {
	a.OrganizationID = values.Get("organizationId")
	a.IndividualID = values.Get("individualId")
	return nil
}

type StaffClient interface {
	Get(ctx context.Context, id string) (*Staff, error)
	Create(ctx context.Context, party *Staff) (*Staff, error)
	Update(ctx context.Context, party *Staff) (*Staff, error)
	List(ctx context.Context, listOptions StaffListOptions) (*StaffList, error)
}
type MembershipListOptions struct {
	IndividualID string
	TeamID       string
}

type NationalityListOptions struct {
	TeamID    string
	CountryID string
}

func (a *NationalityListOptions) MarshalQueryParameters() (url.Values, error) {
	values := url.Values{}
	if len(a.TeamID) > 0 {
		values.Set("teamId", a.TeamID)
	}
	if len(a.CountryID) > 0 {
		values.Set("countryId", a.CountryID)
	}
	return values, nil
}

func (a *NationalityListOptions) UnmarshalQueryParameters(values url.Values) error {
	a.CountryID = values.Get("countryId")
	a.TeamID = values.Get("teamId")
	return nil
}

func (a *MembershipListOptions) MarshalQueryParameters() (url.Values, error) {
	values := url.Values{}
	if len(a.IndividualID) > 0 {
		values.Set("individualId", a.IndividualID)
	}
	if len(a.TeamID) > 0 {
		values.Set("teamId", a.TeamID)
	}
	return values, nil
}

func (a *MembershipListOptions) UnmarshalQueryParameters(values url.Values) error {
	a.TeamID = values.Get("teamId")
	a.IndividualID = values.Get("individualId")
	return nil
}

type MembershipClient interface {
	Get(ctx context.Context, id string) (*Membership, error)
	Create(ctx context.Context, create *Membership) (*Membership, error)
	Update(ctx context.Context, update *Membership) (*Membership, error)
	List(ctx context.Context, listOptions MembershipListOptions) (*MembershipList, error)
}

type NationalityClient interface {
	Get(ctx context.Context, id string) (*Nationality, error)
	Create(ctx context.Context, create *Nationality) (*Nationality, error)
	Update(ctx context.Context, update *Nationality) (*Nationality, error)
	List(ctx context.Context, listOptions NationalityListOptions) (*NationalityList, error)
}

type IndividualListOptions struct {
	PartyTypeIDs []string
	Attributes   map[string]string
	SearchParam  string
	Page         int
	PerPage      int
	Sort         string
}

func (a *IndividualListOptions) MarshalQueryParameters() (url.Values, error) {
	values := url.Values{}
	for _, partyTypeID := range a.PartyTypeIDs {
		values.Add("partyTypeId", partyTypeID)
	}
	values.Add("searchParam", a.SearchParam)
	values.Add("page", strconv.Itoa(a.Page))
	values.Add("perPage", strconv.Itoa(a.PerPage))
	values.Add("sort", a.Sort)
	if a.Attributes != nil {
		for key, value := range a.Attributes {
			if len(value) == 0 {
				continue
			}
			values.Set("attributes["+key+"]", value)
		}
	}
	return values, nil
}

func (a *IndividualListOptions) UnmarshalQueryParameters(values url.Values) error {
	a.PartyTypeIDs = values["partyTypeId"]
	for key, values := range values {
		if containAttributes(key) {
			setAttributesValue(a, key, values)
		}
		switch key {
		case "searchParam":
			a.SearchParam = values[0]
		case "page":
			str := values[0]
			page, err := strconv.Atoi(str)
			if err != nil {
				return err
			}
			a.Page = page
		case "perPage":
			str := values[0]
			perPage, err := strconv.Atoi(str)
			if err != nil {
				return err
			}
			a.PerPage = perPage
		case "sort":
			a.Sort = values[0]
		default:
			break
		}
	}
	return nil
}

func setAttributesValue(a *IndividualListOptions, key string, values []string) {
	if a.Attributes == nil {
		a.Attributes = map[string]string{}
	}
	attrKey := key[11 : len(key)-1]
	a.Attributes[attrKey] = values[0]
}

func containAttributes(key string) bool {
	return strings.HasPrefix(key, "attributes[") && strings.HasSuffix(key, "]")
}

type IndividualClient interface {
	Get(ctx context.Context, id string) (*Individual, error)
	Create(ctx context.Context, party *Individual) (*Individual, error)
	Update(ctx context.Context, party *Individual) (*Individual, error)
	List(ctx context.Context, listOptions IndividualListOptions) (*IndividualList, error)
}