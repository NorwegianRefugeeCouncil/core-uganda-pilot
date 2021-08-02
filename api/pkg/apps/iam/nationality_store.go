package iam

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

type NationalityStore struct {
	relationshipStore *RelationshipStore
}

func NewNationalityStore(relationshipStore *RelationshipStore) *NationalityStore {
	return &NationalityStore{relationshipStore: relationshipStore}
}

func (s *NationalityStore) list(ctx context.Context, listOptions NationalityListOptions) (*NationalityList, error) {

	var relopts = RelationshipListOptions{
		RelationshipTypeID: NationalityRelationshipType.ID,
		FirstPartyID:       listOptions.TeamID,
		//SecondPartyID:      nil,
	}
	got, err := s.relationshipStore.list(ctx, relopts)
	if err != nil {
		return nil, err
	}

	var items = make([]*Nationality, len(got.Items))
	for i, item := range got.Items {
		items[i] = MapRelationshipToNationality(item)
	}

	return &NationalityList{
		Items: items,
	}, nil

}

func (s *NationalityStore) get(ctx context.Context, id string) (*Nationality, error) {
	got, err := s.relationshipStore.get(ctx, id)
	if err != nil {
		return nil, err
	}
	if got.RelationshipTypeID != NationalityRelationshipType.ID {
		return nil, fmt.Errorf("not found")
	}

	return MapRelationshipToNationality(got), nil

}

func (s *NationalityStore) find(ctx context.Context, teamId, countryId string) (*Nationality, error) {
	got, err := s.relationshipStore.list(ctx, RelationshipListOptions{
		RelationshipTypeID: NationalityRelationshipType.ID,
		FirstPartyID:       teamId,
		SecondPartyID:      countryId,
	})
	if err != nil {
		return nil, err
	}
	if len(got.Items) == 0 {
		return nil, err
	}
	return MapRelationshipToNationality(got.Items[0]), nil
}

func (s *NationalityStore) create(ctx context.Context, nationality *Nationality) error {
	got, err := s.find(ctx, nationality.TeamID, nationality.CountryID)
	if err != nil {
		return err
	}
	if got != nil {
		return nil
	}
	rel := MapNationalityToRelationship(nationality)
	if rel.ID == "" {
		rel.ID = uuid.NewV4().String()
	}
	return s.relationshipStore.create(ctx, rel)
}

func MapRelationshipToNationality(rel *Relationship) *Nationality {
	return &Nationality{
		ID:        rel.ID,
		CountryID: rel.SecondPartyID,
		TeamID:    rel.FirstPartyID,
	}
}

func MapNationalityToRelationship(nationality *Nationality) *Relationship {
	return &Relationship{
		ID:                 nationality.ID,
		RelationshipTypeID: NationalityRelationshipType.ID,
		FirstPartyID:       nationality.TeamID,
		SecondPartyID:      nationality.CountryID,
	}
}
