package store

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"reflect"
)

type genericStore struct {
	databaseName   string
	collectionName string
	clientFn       func() (*mongo.Client, error)
}

type GetOptions struct {
	Version string
}

type Object interface {
	GetVersion() int
	SetVersion(version int)
}

func (s *genericStore) Get(ctx context.Context, key string, options GetOptions, out Object) error {

	filter := bson.M{
		"id": key,
	}

	if len(options.Version) > 0 {
		filter["version"] = options.Version
	}

	return s.get(ctx, filter, out)

}

type UpdateFunc func(in Object) (Object, error)

type UpdateOptions struct {
	IgnoreNotFound bool
}

func (s *genericStore) Update(ctx context.Context, key string, updateOptions UpdateOptions, out Object, update UpdateFunc) error {

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	client, err := s.getClient()
	if err != nil {
		return err
	}

	session, err := client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start mongo session: %v", err)
	}
	defer session.EndSession(ctx)

	res, err := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {

		collection := s.getCollection(client)

		filter := bson.M{
			"id":                key,
			"_isCurrentVersion": true,
		}

		// instantiate new object
		oldObj := reflect.New(reflect.ValueOf(out).Elem().Type()).Interface().(Object)

		// find current version
		oldResp := collection.FindOne(sessCtx, filter)
		if oldResp.Err() != nil {
			return nil, fmt.Errorf("failed to get object: %v", oldResp.Err())
		}
		if err := oldResp.Decode(oldObj); err != nil {
			return nil, fmt.Errorf("failed to decode object: %v", err)
		}

		oldVersion := oldObj.GetVersion()

		// update the object using the provided callback
		updatedObj, err := update(oldObj)
		if err != nil {
			return nil, err
		}

		// update old version
		if _, err := collection.UpdateOne(sessCtx, filter, bson.M{"$set": bson.M{"_isCurrentVersion": false}}); err != nil {
			return nil, err
		}

		// insert new version
		updatedObj.SetVersion(oldVersion + 1)
		insertRes, err := collection.InsertOne(sessCtx, updatedObj)
		if err != nil {
			return nil, fmt.Errorf("failed to insert new version: %v", err)
		}

		// set _isCurrentVersion to true
		if _, err = collection.UpdateByID(sessCtx, insertRes.InsertedID, bson.M{"$set": bson.M{"_isCurrentVersion": true}}); err != nil {
			return nil, fmt.Errorf("failed to set current version: %v", err)
		}

		return updatedObj, nil

	}, txnOpts)

	if err != nil {
		return err
	}

	if out != nil {
		reflect.ValueOf(out).Elem().Set(reflect.ValueOf(res).Elem())
	}

	return nil
}

type CreateOptions struct {
}

func (s *genericStore) Create(ctx context.Context, obj Object, out Object, createOptions CreateOptions) error {

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	client, err := s.getClient()
	if err != nil {
		return err
	}

	session, err := client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start mongo session: %v", err)
	}
	defer session.EndSession(ctx)

	res, err := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {

		collection := s.getCollection(client)

		obj.SetVersion(1)
		insertRes, err := collection.InsertOne(sessCtx, obj)
		if err != nil {
			return nil, fmt.Errorf("failed to insert object: %v", err)
		}

		// set _isCurrentVersion to true
		_, err = collection.UpdateByID(sessCtx, insertRes.InsertedID, bson.M{"$set": bson.M{"_isCurrentVersion": true}})
		if err != nil {
			return nil, fmt.Errorf("failed to set current version: %v", err)
		}

		return obj, nil

	}, txnOpts)

	if err != nil {
		return err
	}

	if out != nil {
		reflect.ValueOf(out).Elem().Set(reflect.ValueOf(res).Elem())
	}

	return nil

}

func (s *genericStore) get(ctx context.Context, filter bson.M, out interface{}) error {

	client, err := s.getClient()
	if err != nil {
		return err
	}

	collection := s.getCollection(client)

	res := collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return fmt.Errorf("failed to find object: %v", res.Err())
	}

	if err := res.Decode(&out); err != nil {
		return fmt.Errorf("failed to decode object: %v", err)
	}

	return nil

}

func (s *genericStore) getCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(s.databaseName).Collection(s.collectionName)
}

func (s *genericStore) getClientAndCollection() (*mongo.Client, *mongo.Collection, error) {
	client, err := s.getClient()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get mongo client: %v", err)
	}
	collection := client.Database(s.databaseName).Collection(s.collectionName)
	return client, collection, nil
}

func (s *genericStore) getClient() (*mongo.Client, error) {
	return s.clientFn()
}
