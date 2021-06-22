package staff

import (
	"context"
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/parties/relationships"
)

type Store struct {
	relationshipStore *relationships.Store
}

func NewStore(relationshipStore *relationships.Store) *Store {
	return &Store{relationshipStore: relationshipStore}
}

func (s *Store) Get(ctx context.Context, id string) (*Staff, error) {
	rel, err := s.relationshipStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if rel.RelationshipTypeID != RelationshipType.ID {
		return nil, fmt.Errorf("not found")
	}
	return mapRelationshipToStaff(rel), nil
}

func (s *Store) Find(ctx context.Context, individualId, organizationId string) (*Staff, error) {
	found, err := s.relationshipStore.List(ctx, relationships.ListOptions{
		RelationshipTypeID: RelationshipType.ID,
		FirstPartyId:       individualId,
		SecondParty:        organizationId,
	})
	if err != nil {
		return nil, err
	}
	if len(found.Items) == 0 {
		return nil, err
	}

	return mapRelationshipToStaff(found.Items[0]), nil
}

func (s *Store) Create(ctx context.Context, staff *Staff) error {
	found, err := s.Find(ctx, staff.IndividualID, staff.OrganizationID)
	if err != nil {
		return err
	}
	if found == nil {
		if err := s.relationshipStore.Create(ctx, mapStaffToRelationship(staff)); err != nil {
			return err
		}
	}
	return nil
}

type ListOptions struct {
	IndividualID   string
	OrganizationID string
}

func (s *Store) List(ctx context.Context, listOptions ListOptions) (*StaffList, error) {
	list, err := s.relationshipStore.List(ctx, relationships.ListOptions{
		RelationshipTypeID: RelationshipType.ID,
		FirstPartyId:       listOptions.IndividualID,
		SecondParty:        listOptions.OrganizationID,
	})
	if err != nil {
		return nil, err
	}
	var ret = make([]*Staff, len(list.Items))
	for i, item := range list.Items {
		ret[i] = mapRelationshipToStaff(item)
	}
	return &StaffList{
		Items: ret,
	}, nil
}

func mapStaffToRelationship(staff *Staff) *relationships.Relationship {
	return &relationships.Relationship{
		ID:                 staff.ID,
		RelationshipTypeID: RelationshipType.ID,
		FirstParty:         staff.IndividualID,
		SecondParty:        staff.OrganizationID,
	}
}
func mapRelationshipToStaff(rel *relationships.Relationship) *Staff {
	return &Staff{
		ID:             rel.ID,
		OrganizationID: rel.SecondParty,
		IndividualID:   rel.FirstParty,
	}
}
