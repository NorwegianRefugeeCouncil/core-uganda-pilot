package memberships

import (
	"context"
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/parties/relationships"
	uuid "github.com/satori/go.uuid"
)

type Store struct {
	relationshipStore *relationships.Store
}

func NewStore(relationshipStore *relationships.Store) *Store {
	return &Store{relationshipStore: relationshipStore}
}

type ListOptions struct {
	IndividualID string `json:"individualId"`
	TeamID       string `json:"teamId"`
}

func (s *Store) List(ctx context.Context, listOptions ListOptions) (*MembershipList, error) {

	got, err := s.relationshipStore.List(ctx, relationships.ListOptions{
		RelationshipTypeID: RelationshipType.ID,
		FirstPartyId:       listOptions.IndividualID,
		SecondParty:        listOptions.TeamID,
	})
	if err != nil {
		return nil, err
	}

	var items = make([]*Membership, len(got.Items))
	for i, item := range got.Items {
		items[i] = mapRelationshipToMembership(item)
	}

	return &MembershipList{
		Items: items,
	}, nil

}

func (s *Store) Get(ctx context.Context, id string) (*Membership, error) {
	got, err := s.relationshipStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if got.RelationshipTypeID != RelationshipType.ID {
		return nil, fmt.Errorf("not found")
	}

	return mapRelationshipToMembership(got), nil

}

func (s *Store) Find(ctx context.Context, individualId, teamId string) (*Membership, error) {
	got, err := s.relationshipStore.List(ctx, relationships.ListOptions{
		RelationshipTypeID: RelationshipType.ID,
		FirstPartyId:       individualId,
		SecondParty:        teamId,
	})
	if err != nil {
		return nil, err
	}
	if len(got.Items) == 0 {
		return nil, err
	}
	return mapRelationshipToMembership(got.Items[0]), nil
}

func (s *Store) Create(ctx context.Context, membership *Membership) error {
	got, err := s.Find(ctx, membership.IndividualID, membership.TeamID)
	if err != nil {
		return err
	}
	if got != nil {
		return nil
	}
	rel := mapMembershipToRelationship(membership)
	rel.ID = uuid.NewV4().String()
	return s.relationshipStore.Create(ctx, rel)
}

func mapRelationshipToMembership(rel *relationships.Relationship) *Membership {
	return &Membership{
		ID:           rel.ID,
		TeamID:       rel.SecondParty,
		IndividualID: rel.FirstParty,
	}
}

func mapMembershipToRelationship(membership *Membership) *relationships.Relationship {
	return &relationships.Relationship{
		ID:                 membership.ID,
		RelationshipTypeID: RelationshipType.ID,
		FirstParty:         membership.IndividualID,
		SecondParty:        membership.TeamID,
	}
}
