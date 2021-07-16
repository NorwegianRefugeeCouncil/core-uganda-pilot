package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IndividualStore struct {
	getCollection utils.MongoCollectionFn
}

func NewIndividualStore(mongoClientFn utils.MongoClientFn, database string) *IndividualStore {
	store := &IndividualStore{
		getCollection: utils.GetCollectionFn(database, "parties", mongoClientFn),
	}
	return store
}

func (s *IndividualStore) create(ctx context.Context, individual *Individual) error {
	individual.AddPartyType(IndividualPartyType.ID)
	collection, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, individual)
	if err != nil {
		return err
	}
	return nil
}

func (s *IndividualStore) get(ctx context.Context, ID string) (*Individual, error) {
	collection, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	findResult := collection.FindOne(ctx, bson.M{"id": ID})
	if findResult.Err() != nil {
		return nil, findResult.Err()
	}
	var individual *Individual
	if err := findResult.Decode(&individual); err != nil {
		return nil, err
	}
	return individual, nil
}

func (s *IndividualStore) upsert(ctx context.Context, individual *Individual) error {
	collection, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	individual.AddPartyType(IndividualPartyType.ID)
	_, err = collection.UpdateOne(ctx, bson.M{
		"id": individual.ID,
	}, bson.M{
		"$set": bson.M{
			"attributes":   individual.Attributes,
			"partyTypeIds": individual.PartyTypeIDs,
		},
	}, options.Update().SetUpsert(true))

	if err != nil {
		return err
	}
	return nil
}

func (s *IndividualStore) list(ctx context.Context, listOptions IndividualListOptions) (*IndividualList, error) {

	includesIndividualPartyType := false

	for _, partyTypeId := range listOptions.PartyTypeIDs {
		if partyTypeId == IndividualPartyType.ID {
			includesIndividualPartyType = true
		}
	}

	if !includesIndividualPartyType {
		listOptions.PartyTypeIDs = append(listOptions.PartyTypeIDs, IndividualPartyType.ID)
	}

	filter := bson.M{
		"partyTypeIds": bson.M{
			"$all": listOptions.PartyTypeIDs,
		},
	}

	for key, value := range listOptions.Attributes {
		filter["attributes."+key] = value
	}

	if len(listOptions.SearchParam) != 0 {
		filter["$text"] = bson.M{"$search": listOptions.SearchParam}
	}

	collection, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*Individual
	for {
		if !cursor.Next(ctx) {
			break
		}
		var b Individual
		if err := cursor.Decode(&b); err != nil {
			return nil, err
		}
		items = append(items, &b)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}
	if items == nil {
		items = []*Individual{}
	}
	return &IndividualList{
		Items: items,
	}, nil
}
