package teams

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

func (s *Store) Get(ctx context.Context, id string) (*Team, error) {

	p, err := s.partyStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !p.HasPartyType(PartyType.ID) {
		return nil, fmt.Errorf("not found")
	}

	return mapPartyToTeam(p), nil
}

func mapPartyToTeam(p *parties.Party) *Team {
	name := p.Attributes.Get(TeamNameAttribute.ID)
	return &Team{
		ID:   p.ID,
		Name: name,
	}
}

func mapTeamToParty(team *Team) *parties.Party {
	return &parties.Party{
		ID: team.ID,
		PartyTypeIDs: []string{
			PartyType.ID,
		},
		Attributes: map[string][]string{
			TeamNameAttribute.ID: {team.Name},
		},
	}
}

func (s *Store) List(ctx context.Context) (*TeamList, error) {

	ps, err := s.partyStore.List(ctx, parties.ListOptions{
		PartyTypeID: PartyType.ID,
	})
	if err != nil {
		return nil, err
	}

	teams := make([]*Team, len(ps.Items))
	for i, item := range ps.Items {
		teams[i] = mapPartyToTeam(item)
	}

	return &TeamList{
		Items: teams,
	}, nil
}

func (s *Store) Update(ctx context.Context, team *Team) error {
	party := mapTeamToParty(team)
	if err := s.partyStore.Update(ctx, party); err != nil {
		return err
	}
	return nil
}

func (s *Store) Create(ctx context.Context, team *Team) error {
	party := mapTeamToParty(team)
	if err := s.partyStore.Create(ctx, party); err != nil {
		return err
	}
	return nil
}
