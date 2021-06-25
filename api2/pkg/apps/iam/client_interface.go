package iam

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/organizations"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"github.com/nrc-no/core-kafka/pkg/parties/relationships"
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	"github.com/nrc-no/core-kafka/pkg/teams"
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
	Get(ctx context.Context, id string) (*parties.Party, error)
	Create(ctx context.Context, party *parties.Party) (*parties.Party, error)
	Update(ctx context.Context, party *parties.Party) (*parties.Party, error)
	List(ctx context.Context, listOptions PartyListOptions) (*parties.PartyList, error)
}

type PartyTypeListOptions struct {
	RelationshipTypeID string
	FirstPartyId       string
	SecondParty        string
	EitherParty        string
}

type PartyTypeClient interface {
	Get(ctx context.Context, id string) (*partytypes.PartyType, error)
	Create(ctx context.Context, party *partytypes.PartyType) (*partytypes.PartyType, error)
	Update(ctx context.Context, party *partytypes.PartyType) (*partytypes.PartyType, error)
	List(ctx context.Context, listOptions PartyTypeListOptions) (*partytypes.PartyTypeList, error)
}

type RelationshipListOptions struct {
}

type RelationshipClient interface {
	Get(ctx context.Context, id string) (*relationships.Relationship, error)
	Create(ctx context.Context, party *relationships.Relationship) (*relationships.Relationship, error)
	Update(ctx context.Context, party *relationships.Relationship) (*relationships.Relationship, error)
	List(ctx context.Context, listOptions RelationshipListOptions) (*relationships.RelationshipList, error)
}

type RelationshipTypeListOptions struct {
	PartyType string
}

type RelationshipTypeClient interface {
	Get(ctx context.Context, id string) (*relationshiptypes.RelationshipType, error)
	Create(ctx context.Context, party *relationshiptypes.RelationshipType) (*relationshiptypes.RelationshipType, error)
	Update(ctx context.Context, party *relationshiptypes.RelationshipType) (*relationshiptypes.RelationshipType, error)
	List(ctx context.Context, listOptions RelationshipListOptions) (*relationshiptypes.RelationshipTypeList, error)
}

type AttributeListOptions struct {
	PartyTypeIDs []string
}

type AttributeClient interface {
	Get(ctx context.Context, id string) (*attributes.Attribute, error)
	Create(ctx context.Context, party *attributes.Attribute) (*attributes.Attribute, error)
	Update(ctx context.Context, party *attributes.Attribute) (*attributes.Attribute, error)
	List(ctx context.Context, listOptions AttributeListOptions) (*attributes.AttributeList, error)
}

type TeamListOptions struct {
}

type TeamClient interface {
	Get(ctx context.Context, id string) (*teams.Team, error)
	Create(ctx context.Context, party *teams.Team) (*teams.Team, error)
	Update(ctx context.Context, party *teams.Team) (*teams.Team, error)
	List(ctx context.Context, listOptions TeamListOptions) (*teams.TeamList, error)
}

type OrganizationListOptions struct {
}

type OrganizationClient interface {
	Get(ctx context.Context, id string) (*organizations.Organization, error)
	Create(ctx context.Context, party *organizations.Organization) (*organizations.Organization, error)
	Update(ctx context.Context, party *organizations.Organization) (*organizations.Organization, error)
	List(ctx context.Context, listOptions OrganizationListOptions) (*organizations.OrganizationList, error)
}
