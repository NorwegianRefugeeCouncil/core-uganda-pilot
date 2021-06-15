package beneficiaries

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

func NewStore(mongoClient *mongo.Client) *Store {
	return &Store{
		collection: mongoClient.Database("core").Collection("parties"),
	}
}

func (s *Store) Create(ctx context.Context, beneficiary *Beneficiary) error {
	beneficiary.AddPartyType(partytypes.IndividualPartyType.ID)
	_, err := s.collection.InsertOne(ctx, beneficiary)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Upsert(ctx context.Context, beneficiary *Beneficiary) error {
	beneficiary.AddPartyType(partytypes.IndividualPartyType.ID)
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": beneficiary.ID,
	}, bson.M{
		"$set": bson.M{
			"attributes": beneficiary.Attributes,
			"partyTypes": beneficiary.PartyTypes,
		},
	}, options.Update().SetUpsert(true))

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Get(ctx context.Context, ID string) (*Beneficiary, error) {
	findResult := s.collection.FindOne(ctx, bson.M{"id": ID})
	if findResult.Err() != nil {
		return nil, findResult.Err()
	}
	var beneficiary *Beneficiary
	if err := findResult.Decode(&beneficiary); err != nil {
		return nil, err
	}
	return beneficiary, nil
}

func (s *Store) List(ctx context.Context, listOptions ListOptions) (*BeneficiaryList, error) {

	filter := bson.M{}

	if len(listOptions.PartyTypes) > 0 {
		filter["partyTypes"] = bson.M{
			"$all": listOptions.PartyTypes,
		}
	}

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*Beneficiary
	for {
		if !cursor.Next(ctx) {
			break
		}
		var b Beneficiary
		if err := cursor.Decode(&b); err != nil {
			return nil, err
		}
		items = append(items, &b)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}
	if items == nil {
		items = []*Beneficiary{}
	}
	return &BeneficiaryList{
		Items: items,
	}, nil
}
