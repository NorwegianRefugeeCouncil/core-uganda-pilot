package iam

import (
	"context"
)

type ClientSet interface {
	Parties() PartyClient
	PartyTypes() PartyTypeClient
	Relationships() RelationshipClient
	RelationshipTypes() RelationshipTypeClient
	Attributes() AttributeClient
	Teams() TeamClient
	Organizations() OrganizationClient
}

type PartyListOptions struct {
	PartyTypeID string `json:"partyTypeId" bson:"partyTypeId"`
	SearchParam string `json:"searchParam" bson:"searchParam"`
}

type PartyClient interface {
	Get(ctx context.Context, id string) (*Party, error)
	Create(ctx context.Context, party *Party) (*Party, error)
	Update(ctx context.Context, party *Party) (*Party, error)
	List(ctx context.Context, listOptions PartyListOptions) (*PartyList, error)
}

type PartyTypeListOptions struct {
	RelationshipTypeID string
	FirstPartyId       string
	SecondParty        string
	EitherParty        string
}

type PartyTypeClient interface {
	Get(ctx context.Context, id string) (*PartyType, error)
	Create(ctx context.Context, party *PartyType) (*PartyType, error)
	Update(ctx context.Context, party *PartyType) (*PartyType, error)
	List(ctx context.Context, listOptions PartyTypeListOptions) (*PartyTypeList, error)
}

type RelationshipListOptions struct {
	RelationshipTypeID string
	FirstPartyId       string // TODO FirstPartyID
	SecondParty        string // TODO SecondPartyID
	EitherParty        string // TODO EitherPartyID
}

type RelationshipClient interface {
	Get(ctx context.Context, id string) (*Relationship, error)
	Create(ctx context.Context, party *Relationship) (*Relationship, error)
	Update(ctx context.Context, party *Relationship) (*Relationship, error)
	List(ctx context.Context, listOptions RelationshipListOptions) (*RelationshipList, error)
}

type RelationshipTypeListOptions struct {
	PartyType string
}

type RelationshipTypeClient interface {
	Get(ctx context.Context, id string) (*RelationshipType, error)
	Create(ctx context.Context, party *RelationshipType) (*RelationshipType, error)
	Update(ctx context.Context, party *RelationshipType) (*RelationshipType, error)
	List(ctx context.Context, listOptions RelationshipListOptions) (*RelationshipTypeList, error)
}

type AttributeListOptions struct {
	PartyTypeIDs []string
}

type AttributeClient interface {
	Get(ctx context.Context, id string) (*Attribute, error)
	Create(ctx context.Context, party *Attribute) (*Attribute, error)
	Update(ctx context.Context, party *Attribute) (*Attribute, error)
	List(ctx context.Context, listOptions AttributeListOptions) (*AttributeList, error)
}

type TeamListOptions struct {
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
	Create(ctx context.Context, party *Membership) (*Membership, error)
	Update(ctx context.Context, party *Membership) (*Membership, error)
	List(ctx context.Context, listOptions StaffListOptions) (*MembershipList, error)
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
