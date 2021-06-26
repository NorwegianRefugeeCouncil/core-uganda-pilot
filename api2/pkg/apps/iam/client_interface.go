package iam

import (
	"context"
	"net/url"
)

type Interface interface {
	Parties() PartyClient
	PartyTypes() PartyTypeClient
	Relationships() RelationshipClient
	RelationshipTypes() RelationshipTypeClient
	Attributes() AttributeClient
	Teams() TeamClient
	Organizations() OrganizationClient
	Staff() StaffClient
	Memberships() MembershipClient
	Individuals() IndividualClient
}

type PartyListOptions struct {
	PartyTypeID string `json:"partyTypeId" bson:"partyTypeId"`
	SearchParam string `json:"searchParam" bson:"searchParam"`
}

func (a *PartyListOptions) MarshalQueryParameters() (url.Values, error) {
	ret := url.Values{}
	if len(a.PartyTypeID) > 0 {
		ret.Set("partyTypeId", a.PartyTypeID)
	}
	if len(a.SearchParam) > 0 {
		ret.Set("searchParam", a.SearchParam)
	}
	return ret, nil
}

func (a *PartyListOptions) UnmarshalQueryParameters(values url.Values) error {
	a.PartyTypeID = values.Get("partyTypeId")
	a.SearchParam = values.Get("searchParam")
	return nil
}

type PartyClient interface {
	Get(ctx context.Context, id string) (*Party, error)
	Create(ctx context.Context, party *Party) (*Party, error)
	Update(ctx context.Context, party *Party) (*Party, error)
	List(ctx context.Context, listOptions PartyListOptions) (*PartyList, error)
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

type AttributeListOptions struct {
	PartyTypeIDs []string
}

func (a *AttributeListOptions) MarshalQueryParameters() (url.Values, error) {
	ret := url.Values{}
	for _, partyTypeID := range a.PartyTypeIDs {
		ret.Add("partyTypeId", partyTypeID)
	}
	return ret, nil
}

func (a *AttributeListOptions) UnmarshalQueryParameters(values url.Values) error {
	partyTypeIDs := values["partyTypeId"]
	for _, partyTypeID := range partyTypeIDs {
		a.PartyTypeIDs = append(a.PartyTypeIDs, partyTypeID)
	}
	return nil
}

var _ UrlValuer = &AttributeListOptions{}

type AttributeClient interface {
	Get(ctx context.Context, id string) (*Attribute, error)
	Create(ctx context.Context, create *Attribute) (*Attribute, error)
	Update(ctx context.Context, update *Attribute) (*Attribute, error)
	List(ctx context.Context, listOptions AttributeListOptions) (*AttributeList, error)
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

type OrganizationListOptions struct {
}

type OrganizationClient interface {
	Get(ctx context.Context, id string) (*Organization, error)
	Create(ctx context.Context, party *Organization) (*Organization, error)
	Update(ctx context.Context, party *Organization) (*Organization, error)
	List(ctx context.Context, listOptions OrganizationListOptions) (*OrganizationList, error)
}

type StaffListOptions struct {
	IndividualID   string
	OrganizationID string
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

type MembershipClient interface {
	Get(ctx context.Context, id string) (*Membership, error)
	Create(ctx context.Context, create *Membership) (*Membership, error)
	Update(ctx context.Context, update *Membership) (*Membership, error)
	List(ctx context.Context, listOptions MembershipListOptions) (*MembershipList, error)
}

type IndividualListOptions struct {
	PartyTypeIDs []string `json:"partyTypeIds" bson:"partyTypeIds"`
}

type IndividualClient interface {
	Get(ctx context.Context, id string) (*Individual, error)
	Create(ctx context.Context, party *Individual) (*Individual, error)
	Update(ctx context.Context, party *Individual) (*Individual, error)
	List(ctx context.Context, listOptions IndividualListOptions) (*IndividualList, error)
}
