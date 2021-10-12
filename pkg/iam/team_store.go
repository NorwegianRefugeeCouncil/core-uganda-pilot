package iam

import (
	"context"
	"fmt"
)

type TeamStore struct {
	partyStore *PartyStore
}

func NewTeamStore(partyStore *PartyStore) *TeamStore {
	return &TeamStore{
		partyStore: partyStore,
	}
}

func (s *TeamStore) Get(ctx context.Context, id string) (*Team, error) {
	p, err := s.partyStore.get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !p.HasPartyType(TeamPartyType.ID) {
		return nil, fmt.Errorf("not found")
	}

	return MapPartyToTeam(p), nil
}

func MapPartyToTeam(p *Party) *Team {
	name := p.Attributes.Get(TeamNameAttribute.ID)
	return &Team{
		ID:   p.ID,
		Name: name,
	}
}

func MapTeamToParty(team *Team) *Party {
	return &Party{
		ID: team.ID,
		PartyTypeIDs: []string{
			TeamPartyType.ID,
		},
		Attributes: map[string][]string{
			TeamNameAttribute.ID: {team.Name},
		},
	}
}

func (s *TeamStore) List(ctx context.Context) (*TeamList, error) {
	ps, err := s.partyStore.list(ctx, PartySearchOptions{
		PartyTypeIDs: []string{TeamPartyType.ID},
	})
	if err != nil {
		return nil, err
	}

	teams := make([]*Team, len(ps.Items))
	for i, item := range ps.Items {
		teams[i] = MapPartyToTeam(item)
	}

	return &TeamList{
		Items: teams,
	}, nil
}

func (s *TeamStore) Update(ctx context.Context, team *Team) error {
	party := MapTeamToParty(team)
	return s.partyStore.update(ctx, party)
}

func (s *TeamStore) Create(ctx context.Context, team *Team) error {
	party := MapTeamToParty(team)
	return s.partyStore.create(ctx, party)
}
