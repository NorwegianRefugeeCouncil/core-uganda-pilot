package iam

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PartyStore struct {
	Collection *mongo.Collection
}

func newPartyStore(ctx context.Context, mongoClient *mongo.Client, database string) (*PartyStore, error) {
	store := &PartyStore{
		Collection: mongoClient.Database(database).Collection("parties"),
	}

	if _, err := store.Collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"id": 1},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, err
	}

	// first name and last name full text index
	if _, err := store.Collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{
				"attributes." + FirstNameAttribute.ID, "text",
			},
			{
				"attributes." + LastNameAttribute.ID, "text",
			},
		},
	}); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *PartyStore) get(ctx context.Context, id string) (*Party, error) {
	res := s.Collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r Party
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func BSONStringA(strSlice []string) (result bson.A) {
	result = bson.A{}
	for _, s := range strSlice {
		result = append(result, s)
	}
	return
}

func (s *PartyStore) list(ctx context.Context, listOptions PartySearchOptions) (*PartyList, error) {
	filterItems := bson.M{}

	if len(listOptions.PartyTypeIDs) != 0 {
		filterItems["partyTypeIds"] = bson.D{{"$in", BSONStringA(listOptions.PartyTypeIDs)}}
	}

	if len(listOptions.PartyIDs) != 0 {
		filterItems["id"] = bson.D{{"$in", BSONStringA(listOptions.PartyIDs)}}
	}

	if listOptions.Attributes != nil && len(listOptions.Attributes) > 0 {
		for key, value := range listOptions.Attributes {
			filterItems["attributes."+key] = value
		}
	}

	if len(listOptions.SearchParam) != 0 {
		filterItems["$text"] = bson.M{"$search": listOptions.SearchParam}
	}

	res, err := s.Collection.Find(ctx, filterItems)
	if err != nil {
		return nil, err
	}
	var items []*Party
	for {
		if !res.Next(ctx) {
			break
		}
		var r Party
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	if items == nil {
		items = []*Party{}
	}
	ret := PartyList{
		Items: items,
	}
	return &ret, nil
}

func (s *PartyStore) update(ctx context.Context, party *Party) error {
	_, err := s.Collection.UpdateOne(ctx, bson.M{
		"id": party.ID,
	}, bson.M{
		"$set": bson.M{
			"attributes":   party.Attributes,
			"partyTypeIds": party.PartyTypeIDs,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *PartyStore) create(ctx context.Context, party *Party) error {
	_, err := s.Collection.InsertOne(ctx, party)
	if err != nil {
		return err
	}
	return nil
}

type FindOptions struct {
	Attributes map[string]string
}

func (s *PartyStore) find(ctx context.Context, options FindOptions) (*Party, error) {
	filter := bson.M{}
	for key, value := range options.Attributes {
		filter["attributes."+key] = value
	}
	res := s.Collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}
	var party *Party
	if err := res.Decode(&party); err != nil {
		return nil, err
	}
	return party, nil
}
