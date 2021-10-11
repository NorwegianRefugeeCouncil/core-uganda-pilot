package iam

import (
	"context"
	"github.com/nrc-no/core/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RelationshipStore struct {
	getCollection utils.MongoCollectionFn
}

func newRelationshipStore(ctx context.Context, mongoClientFn utils.MongoClientFn, database string) (*RelationshipStore, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	store := &RelationshipStore{
		getCollection: utils.GetCollectionFn(database, "relationships", mongoClientFn),
	}

	collection, done, err := store.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	defer done()

	if _, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"id": 1},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, err
	}

	return store, nil

}

func (s *RelationshipStore) create(ctx context.Context, relationship *Relationship) error {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	defer done()

	_, err = collection.InsertOne(ctx, relationship)
	if err != nil {
		return err
	}
	return nil
}

func (s *RelationshipStore) get(ctx context.Context, id string) (*Relationship, error) {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	defer done()

	res := collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r Relationship
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *RelationshipStore) update(ctx context.Context, relationship *Relationship) error {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	defer done()

	_, err = collection.UpdateOne(ctx, bson.M{
		"id": relationship.ID,
	}, bson.M{
		"$set": bson.M{
			"firstParty":         relationship.FirstPartyID,
			"secondParty":        relationship.SecondPartyID,
			"relationshipTypeId": relationship.RelationshipTypeID,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *RelationshipStore) list(ctx context.Context, listOptions RelationshipListOptions) (*RelationshipList, error) {

	filter := bson.M{}

	if len(listOptions.RelationshipTypeID) > 0 {
		filter["relationshipTypeId"] = listOptions.RelationshipTypeID
	}

	if len(listOptions.EitherPartyID) > 0 {
		filter["$or"] = bson.A{
			bson.M{"firstParty": listOptions.EitherPartyID},
			bson.M{"secondParty": listOptions.EitherPartyID},
		}
	} else {
		if len(listOptions.FirstPartyID) > 0 {
			filter["firstParty"] = listOptions.FirstPartyID
		}
		if len(listOptions.SecondPartyID) > 0 {
			filter["secondParty"] = listOptions.SecondPartyID
		}
	}

	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	defer done()

	res, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*Relationship
	for {
		if !res.Next(ctx) {
			break
		}
		var r Relationship
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	ret := RelationshipList{
		Items: items,
	}
	return &ret, nil
}

func (s *RelationshipStore) delete(ctx context.Context, id string) error {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	defer done()

	_, err = collection.DeleteOne(ctx, bson.M{
		"id": id,
	})
	if err != nil {
		return err
	}
	return nil
}
