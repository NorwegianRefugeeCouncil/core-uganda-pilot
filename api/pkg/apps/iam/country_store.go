package iam

import (
	"context"
	"fmt"
)

type CountryStore struct {
	partyStore *PartyStore
}

func NewCountryStore(partyStore *PartyStore) *CountryStore {
	return &CountryStore{
		partyStore: partyStore,
	}
}

func (s *CountryStore) Get(ctx context.Context, id string) (*Country, error) {

	p, err := s.partyStore.get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !p.HasPartyType(CountryPartyType.ID) {
		return nil, fmt.Errorf("not found")
	}

	return MapPartyToCountry(p), nil
}

func MapPartyToCountry(p *Party) *Country {
	name := p.Attributes.Get(CountryNameAttribute.ID)
	return &Country{
		ID:   p.ID,
		Name: name,
	}
}

func MapCountryToParty(country *Country) *Party {
	return &Party{
		ID: country.ID,
		PartyTypeIDs: []string{
			CountryPartyType.ID,
		},
		Attributes: map[string][]string{
			CountryNameAttribute.ID: {country.Name},
		},
	}
}

func (s *CountryStore) List(ctx context.Context) (*CountryList, error) {
	ps, err := s.partyStore.list(ctx, PartySearchOptions{
		PartyTypeIDs: []string{CountryPartyType.ID},
	})
	if err != nil {
		return nil, err
	}

	countries := make([]*Country, len(ps.Items))
	for i, item := range ps.Items {
		countries[i] = MapPartyToCountry(item)
	}

	return &CountryList{
		Items: countries,
	}, nil
}

func (s *CountryStore) Update(ctx context.Context, country *Country) error {
	party := MapCountryToParty(country)
	if err := s.partyStore.update(ctx, party); err != nil {
		return err
	}
	return nil
}

func (s *CountryStore) Create(ctx context.Context, country *Country) error {
	party := MapCountryToParty(country)
	if err := s.partyStore.create(ctx, party); err != nil {
		return err
	}
	return nil
}
