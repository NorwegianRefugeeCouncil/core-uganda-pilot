package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AttributeStore struct {
	getCollection utils.MongoCollectionFn
}

func newAttributeStore(ctx context.Context, mongoClientFn utils.MongoClientFn, database string) (*AttributeStore, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	store := &AttributeStore{
		getCollection: utils.GetCollectionFn(database, "attributes", mongoClientFn),
	}

	collection, done, err := store.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	defer done()

	if _, err := collection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return nil, err
	}

	return store, nil
}

// list returns an AttributeList. If AttributeListOptions are supplied, list returns a filtered list containing
// only those items whose Attribute.PartyTypeIDs field contains all the elements given in the query.
func (s *AttributeStore) list(ctx context.Context, listOptions AttributeListOptions) (*AttributeList, error) {

	filter := bson.M{}

	if len(listOptions.PartyTypeIDs) > 0 {
		filter["partyTypeIds"] = bson.M{
			"$all": listOptions.PartyTypeIDs,
		}
	}

	logrus.Infof("CountryIDs len: %d", len(listOptions.CountryIDs))
	if len(listOptions.CountryIDs) > 0 {
		filter["countryId"] = bson.M{
			"$in": listOptions.CountryIDs,
		}
	}

	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	defer done()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var list []*Attribute
	for {
		if !cursor.Next(ctx) {
			break
		}
		var a Attribute
		if err := cursor.Decode(&a); err != nil {
			return nil, err
		}
		list = append(list, &a)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}

	return &AttributeList{
		Items: list,
	}, nil

}

func (s *AttributeStore) create(ctx context.Context, attribute *Attribute) error {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	defer done()

	_, err = collection.InsertOne(ctx, attribute)
	if err != nil {
		return err
	}
	return nil
}

func (s *AttributeStore) get(ctx context.Context, id string) (*Attribute, error) {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	defer done()

	result := collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if result.Err() != nil {
		return nil, result.Err()
	}
	var a Attribute
	if err := result.Decode(&a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AttributeStore) update(ctx context.Context, attribute *Attribute) error {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	defer done()

	_, err = collection.UpdateOne(ctx, bson.M{
		"id": attribute.ID,
	}, bson.M{
		"$set": bson.M{
			"name":         attribute.Name,
			"translations": attribute.Translations,
			"isPii":        attribute.IsPersonallyIdentifiableInfo,
			"partyTypeIds": attribute.PartyTypeIDs,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
