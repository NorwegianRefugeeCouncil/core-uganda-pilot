package iam

import (
	"context"
	"fmt"
)

type StaffStore struct {
	relationshipStore *RelationshipStore
}

func NewStaffStore(relationshipStore *RelationshipStore) *StaffStore {
	return &StaffStore{relationshipStore: relationshipStore}
}

func (s *StaffStore) Get(ctx context.Context, id string) (*Staff, error) {
	rel, err := s.relationshipStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if rel.RelationshipTypeID != StaffRelationshipType.ID {
		return nil, fmt.Errorf("not found")
	}
	return MapRelationshipToStaff(rel), nil
}

func (s *StaffStore) Find(ctx context.Context, individualId, organizationId string) (*Staff, error) {
	found, err := s.relationshipStore.List(ctx, RelationshipListOptions{
		RelationshipTypeID: StaffRelationshipType.ID,
		FirstPartyID:       individualId,
		SecondPartyID:      organizationId,
	})
	if err != nil {
		return nil, err
	}
	if len(found.Items) == 0 {
		return nil, err
	}

	return MapRelationshipToStaff(found.Items[0]), nil
}

func (s *StaffStore) Create(ctx context.Context, staff *Staff) error {
	found, err := s.Find(ctx, staff.IndividualID, staff.OrganizationID)
	if err != nil {
		return err
	}
	if found == nil {
		if err := s.relationshipStore.Create(ctx, MapStaffToRelationship(staff)); err != nil {
			return err
		}
	}
	return nil
}

func (s *StaffStore) List(ctx context.Context, listOptions StaffListOptions) (*StaffList, error) {
	list, err := s.relationshipStore.List(ctx, RelationshipListOptions{
		RelationshipTypeID: StaffRelationshipType.ID,
		FirstPartyID:       listOptions.IndividualID,
		SecondPartyID:      listOptions.OrganizationID,
	})
	if err != nil {
		return nil, err
	}
	var ret = make([]*Staff, len(list.Items))
	for i, item := range list.Items {
		ret[i] = MapRelationshipToStaff(item)
	}
	return &StaffList{
		Items: ret,
	}, nil
}

func MapStaffToRelationship(staff *Staff) *Relationship {
	return &Relationship{
		ID:                 staff.ID,
		RelationshipTypeID: StaffRelationshipType.ID,
		FirstPartyID:       staff.IndividualID,
		SecondPartyID:      staff.OrganizationID,
	}
}
func MapRelationshipToStaff(rel *Relationship) *Staff {
	return &Staff{
		ID:             rel.ID,
		OrganizationID: rel.SecondPartyID,
		IndividualID:   rel.FirstPartyID,
	}
}
