package organizations

import (
	"context"
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
)

type Store struct {
	partyStore *parties.Store
}

func NewStore(partyStore *parties.Store) *Store {
	return &Store{
		partyStore: partyStore,
	}
}

func (s *Store) Get(ctx context.Context, id string) (*Organization, error) {
	party, err := s.partyStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !party.HasPartyType(PartyType.ID) {
		return nil, fmt.Errorf("not found")
	}

	return &Organization{
		Party: party,
	}, nil
}

func (s *Store) Update(ctx context.Context, organization *Organization) error {

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

func (s *Store) Create(ctx context.Context, organization *Organization) error {

	// Make sure the party has the Organization party type
	organization.AddPartyType(PartyType.ID)

	if err := s.partyStore.Create(ctx, organization.Party); err != nil {
		return err
	}
	return nil
}
