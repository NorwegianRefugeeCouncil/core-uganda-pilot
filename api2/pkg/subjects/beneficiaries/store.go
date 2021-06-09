package beneficiaries

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/subjects/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(mongoClient *mongo.Client) *Store {
	return &Store{
		collection: mongoClient.Database("core").Collection("beneficiaries"),
	}
}

func (s *Store) Create(ctx context.Context, beneficiary *api.Beneficiary) error {
	_, err := s.collection.InsertOne(ctx, beneficiary)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Upsert(ctx context.Context, ID string, attributes []*api.AttributeValue) error {
	bsonAttributes := bson.M{
		"id": ID,
	}
	for _, attribute := range attributes {
		bsonAttributes["attributes."+attribute.Name] = attribute
	}
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": ID,
	}, bson.M{
		"$set": bsonAttributes,
	}, options.Update().SetUpsert(true))

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Get(ctx context.Context, ID string) (*api.Beneficiary, error) {
	findResult := s.collection.FindOne(ctx, bson.M{"id": ID})
	if findResult.Err() != nil {
		return nil, findResult.Err()
	}
	var beneficiary *api.Beneficiary
	if err := findResult.Decode(&beneficiary); err != nil {
		return nil, err
	}
	return beneficiary, nil
}

func (s *Store) List(ctx context.Context) (*api.BeneficiaryList, error) {

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var items []*api.Beneficiary
	for {
		if !cursor.Next(ctx) {
			break
		}
		var b api.Beneficiary
		if err := cursor.Decode(&b); err != nil {
			return nil, err
		}
		items = append(items, &b)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}

	return &api.BeneficiaryList{
		Items: items,
	}, nil

}
