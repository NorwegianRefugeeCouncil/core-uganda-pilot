package iam

import (
	"context"
	"fmt"
)

type OrganizationStore struct {
	partyStore *PartyStore
}

func NewOrganizationStore(partyStore *PartyStore) *OrganizationStore {
	return &OrganizationStore{
		partyStore: partyStore,
	}
}

func (s *OrganizationStore) Get(ctx context.Context, id string) (*Organization, error) {
	party, err := s.partyStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !party.HasPartyType(OrganizationPartyType.ID) {
		return nil, fmt.Errorf("not found")
	}

	return &Organization{
		Party: party,
	}, nil
}

func (s *OrganizationStore) Update(ctx context.Context, organization *Organization) error {

	// Make sure that the party is an organization
	_, err := s.Get(ctx, organization.ID)
	if err != nil {
		return err
	}

	if err := s.partyStore.Update(ctx, organization.Party); err != nil {
		return err
	}

	return nil
}

func (s *OrganizationStore) Create(ctx context.Context, organization *Organization) error {

	// Make sure the party has the Organization party type
	organization.AddPartyType(OrganizationPartyType.ID)

	if err := s.partyStore.Create(ctx, organization.Party); err != nil {
		return err
	}
	return nil
}
