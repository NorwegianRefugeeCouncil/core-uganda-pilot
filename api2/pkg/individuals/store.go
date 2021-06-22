package individuals

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(mongoClient *mongo.Client, database string) *Store {
	return &Store{
		collection: mongoClient.Database(database).Collection("parties"),
	}
}

func (s *Store) Create(ctx context.Context, individual *Individual) error {
	individual.AddPartyType(partytypes.IndividualPartyType.ID)
	_, err := s.collection.InsertOne(ctx, individual)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Upsert(ctx context.Context, individual *Individual) error {
	individual.AddPartyType(partytypes.IndividualPartyType.ID)
	_, err := s.collection.UpdateOne(ctx, bson.M{
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

func (s *Store) Get(ctx context.Context, ID string) (*Individual, error) {
	findResult := s.collection.FindOne(ctx, bson.M{"id": ID})
	if findResult.Err() != nil {
		return nil, findResult.Err()
	}
	var individual *Individual
	if err := findResult.Decode(&individual); err != nil {
		return nil, err
	}
	return individual, nil
}

func (s *Store) List(ctx context.Context, listOptions ListOptions) (*IndividualList, error) {

	filter := bson.M{}

	if len(listOptions.PartyTypeIDs) > 0 {
		filter["partyTypeIds"] = bson.M{
			"$all": listOptions.PartyTypeIDs,
		}
	}

	cursor, err := s.collection.Find(ctx, filter)
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
