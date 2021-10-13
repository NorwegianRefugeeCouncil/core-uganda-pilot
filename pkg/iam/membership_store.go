package iam

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

type MembershipStore struct {
	relationshipStore *RelationshipStore
}

func NewMembershipStore(relationshipStore *RelationshipStore) *MembershipStore {
	return &MembershipStore{relationshipStore: relationshipStore}
}

func (s *MembershipStore) list(ctx context.Context, listOptions MembershipListOptions) (*MembershipList, error) {
	got, err := s.relationshipStore.list(ctx, RelationshipListOptions{
		RelationshipTypeID: MembershipRelationshipType.ID,
		FirstPartyID:       listOptions.IndividualID,
		SecondPartyID:      listOptions.TeamID,
	})
	if err != nil {
		return nil, err
	}

	var items = make([]*Membership, len(got.Items))
	for i, item := range got.Items {
		items[i] = MapRelationshipToMembership(item)
	}

	return &MembershipList{
		Items: items,
	}, nil
}

func (s *MembershipStore) get(ctx context.Context, id string) (*Membership, error) {
	got, err := s.relationshipStore.get(ctx, id)
	if err != nil {
		return nil, err
	}
	if got.RelationshipTypeID != MembershipRelationshipType.ID {
		return nil, fmt.Errorf("not found")
	}

	return MapRelationshipToMembership(got), nil
}

func (s *MembershipStore) find(ctx context.Context, individualID, teamId string) (*Membership, error) {
	got, err := s.relationshipStore.list(ctx, RelationshipListOptions{
		RelationshipTypeID: MembershipRelationshipType.ID,
		FirstPartyID:       individualID,
		SecondPartyID:      teamId,
	})
	if err != nil {
		return nil, err
	}
	if len(got.Items) == 0 {
		return nil, err
	}
	return MapRelationshipToMembership(got.Items[0]), nil
}

func (s *MembershipStore) create(ctx context.Context, membership *Membership) error {
	got, err := s.find(ctx, membership.IndividualID, membership.TeamID)
	if err != nil {
		return err
	}
	if got != nil {
		return nil
	}
	rel := MapMembershipToRelationship(membership)
	if rel.ID == "" {
		rel.ID = uuid.NewV4().String()
	}
	return s.relationshipStore.create(ctx, rel)
}

func MapRelationshipToMembership(rel *Relationship) *Membership {
	return &Membership{
		ID:           rel.ID,
		TeamID:       rel.SecondPartyID,
		IndividualID: rel.FirstPartyID,
	}
}

func MapMembershipToRelationship(membership *Membership) *Relationship {
	return &Relationship{
		ID:                 membership.ID,
		RelationshipTypeID: MembershipRelationshipType.ID,
		FirstPartyID:       membership.IndividualID,
		SecondPartyID:      membership.TeamID,
	}
}