package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PartyStore struct {
	getCollection func() (*mongo.Collection, error)
}

func newPartyStore(ctx context.Context, mongoClientSrc storage.MongoClientSrc, database string) (*PartyStore, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	store := &PartyStore{
		getCollection: storage.GetCollectionFn(mongoClientSrc, database, "parties"),
	}

	collection, err := store.getCollection()
	if err != nil {
		return nil, err
	}

	if _, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"id": 1},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, err
	}

	// full and display names full text index
	if _, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{
				Key: "attributes." + FullNameAttribute.ID, Value: "text",
			},
			{
				Key: "attributes." + DisplayNameAttribute.ID, Value: "text",
			},
		},
	}); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *PartyStore) get(ctx context.Context, id string) (*Party, error) {
	collection, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	res := collection.FindOne(ctx, bson.M{
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
		filterItems["partyTypeIds"] = bson.D{{Key: "$in", Value: BSONStringA(listOptions.PartyTypeIDs)}}
	}

	if len(listOptions.PartyIDs) != 0 {
		filterItems["id"] = bson.D{{Key: "$in", Value: BSONStringA(listOptions.PartyIDs)}}
	}

	if listOptions.Attributes != nil && len(listOptions.Attributes) > 0 {
		for key, value := range listOptions.Attributes {
			filterItems["attributes."+key] = value
		}
	}

	if len(listOptions.SearchParam) != 0 {
		filterItems["$text"] = bson.M{"$search": listOptions.SearchParam}
	}

	collection, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	res, err := collection.Find(ctx, filterItems)
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
	collection, err := s.getCollection()
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(ctx, bson.M{
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
	collection, err := s.getCollection()
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(ctx, party)
	if err != nil {
		return err
	}
	return nil
}

type FindOptions struct {
	Attributes map[string]string
}

// func (s *PartyStore) find(ctx context.Context, options FindOptions) (*Party, error) {
//	filter := bson.M{}
//	for key, value := range options.Attributes {
//		filter["attributes."+key] = value
//	}
//	collection, done, err := s.GetCollection(ctx)
//	if err != nil {
//		return nil, err
//	}
//	res := collection.FindOne(ctx, filter)
//	if res.Err() != nil {
//		return nil, res.Err()
//	}
//	var party *Party
//	if err := res.Decode(&party); err != nil {
//		return nil, err
//	}
//	return party, nil
//}
