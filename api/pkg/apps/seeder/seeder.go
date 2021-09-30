package seeder

import (
	"context"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/apps/login"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func Clear(ctx context.Context, mongoClientFn utils.MongoClientFn, databaseName string) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mongoClient, err := mongoClientFn(ctx)
	if err != nil {
		return err
	}

	return mongoClient.Database(databaseName).Drop(ctx)
}

func Seed(ctx context.Context, mongoClientFn utils.MongoClientFn, databaseName string) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mongoClient, err := mongoClientFn(ctx)
	if err != nil {
		return err
	}

	err = initCollection(ctx, mongoClient, databaseName)
	if err != nil {
		return err
	}

	for _, obj := range teams {
		if err := seedMongo(ctx, mongoClient, databaseName, "parties", bson.M{"id": obj.ID}, iam.MapTeamToParty(&obj)); err != nil {
			return err
		}
	}

	for _, obj := range relationships {
		if err := seedMongo(ctx, mongoClient, databaseName, "relationships", bson.M{"id": obj.ID}, obj); err != nil {
			return err
		}
	}

	for _, obj := range memberships {
		if err := seedMongo(ctx, mongoClient, databaseName, "relationships", bson.M{"id": obj.ID}, iam.MapMembershipToRelationship(&obj)); err != nil {
			return err
		}
	}

	for _, obj := range identificationDocumentTypes {
		if err := seedMongo(ctx, mongoClient, databaseName, "identificationDocumentTypes", bson.M{"id": obj.ID}, obj); err != nil {
			return err
		}
	}

	for _, obj := range identificationDocuments {
		if err := seedMongo(ctx, mongoClient, databaseName, "identificationDocuments", bson.M{"id": obj.ID}, obj); err != nil {
			return err
		}
	}

	for _, obj := range caseTypes {
		if err := seedMongo(ctx, mongoClient, databaseName, "caseTypes", bson.M{"id": obj.ID}, obj); err != nil {
			return err
		}
	}

	for _, obj := range cases {
		if err := seedMongo(ctx, mongoClient, databaseName, "cases", bson.M{"id": obj.ID}, obj); err != nil {
			return err
		}
	}

	for _, obj := range countries {
		if err := seedMongo(ctx, mongoClient, databaseName, "parties", bson.M{"id": obj.ID}, iam.MapCountryToParty(&obj)); err != nil {
			return err
		}
	}

	for _, obj := range nationalities {
		if err := seedMongo(ctx, mongoClient, databaseName, "relationships", bson.M{"id": obj.ID}, iam.MapNationalityToRelationship(&obj)); err != nil {
			return err
		}
	}

	return nil
}

func initCollection(ctx context.Context, mongoClient *mongo.Client, databaseName string) error {
	for _, obj := range individuals {
		if err := seedMongo(ctx, mongoClient, databaseName, "parties", bson.M{"id": obj.ID}, obj.Party); err != nil {
			return err
		}
		if obj.HasPartyType(iam.StaffPartyType.ID) {
			hash, err := login.HashAndSalt(bcrypt.MinCost, []byte("password"))
			if err != nil {
				return err
			}
			if _, err := mongoClient.Database(databaseName).Collection("credentials").UpdateOne(ctx,
				bson.M{
					"partyId": obj.ID,
				},
				bson.M{
					"$set": bson.M{
						"partyId": obj.ID,
						"hash":    hash,
					},
				},
				options.Update().SetUpsert(true)); err != nil {
				return err
			}
		}
	}
	return nil
}

func seedMongo(ctx context.Context, mongoClient *mongo.Client, databaseName, collectionName string, filter interface{}, document interface{}) error {
	logrus.Infof("seeding collection %s.%s with object: %#v", databaseName, collectionName, document)
	collection := mongoClient.Database(databaseName).Collection(collectionName)
	if _, err := collection.InsertOne(ctx, document); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
		if _, err := collection.ReplaceOne(ctx, filter, document); err != nil {
			return err
		}
	}
	return nil
}
