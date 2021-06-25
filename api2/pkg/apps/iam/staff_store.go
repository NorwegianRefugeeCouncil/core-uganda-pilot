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
	return mapRelationshipToStaff(rel), nil
}

func (s *StaffStore) Find(ctx context.Context, individualId, organizationId string) (*Staff, error) {
	found, err := s.relationshipStore.List(ctx, RelationshipListOptions{
		RelationshipTypeID: StaffRelationshipType.ID,
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

func (s *StaffStore) Create(ctx context.Context, staff *Staff) error {
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

func (s *StaffStore) List(ctx context.Context, listOptions StaffListOptions) (*StaffList, error) {
	list, err := s.relationshipStore.List(ctx, RelationshipListOptions{
		RelationshipTypeID: StaffRelationshipType.ID,
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

func mapStaffToRelationship(staff *Staff) *Relationship {
	return &Relationship{
		ID:                 staff.ID,
		RelationshipTypeID: StaffRelationshipType.ID,
		FirstParty:         staff.IndividualID,
		SecondParty:        staff.OrganizationID,
	}
}
func mapRelationshipToStaff(rel *Relationship) *Staff {
	return &Staff{
		ID:             rel.ID,
		OrganizationID: rel.SecondParty,
		IndividualID:   rel.FirstParty,
	}
}
