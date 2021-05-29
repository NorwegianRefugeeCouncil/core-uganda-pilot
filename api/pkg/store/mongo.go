package store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	fields2 "github.com/nrc-no/core/api/pkg/fields"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	IdKey                = "_id"
	Key                  = "key"
	PreviousIDKey        = "previousId"
	TimestampKey         = "timestamp"
	PreviousTimestampKey = "previousTimestamp"
	NextTimestampKey     = "nextTimestamp"
	CreationTimestampKey = "creationTimestamp"
	DeletionTimestampKey = "deletionTimestamp"
	NextIDKey            = "nextId"
	CurrentValueKey      = "currentValue"
	PreviousValueKey     = "previousValue"
	IsDeletedKey         = "isDeleted"
	IsCreatedKey         = "isCreated"
	IsCurrentKey         = "isCurrent"
)

// MongoStore mongo store implementation
type MongoStore struct {
	databaseName string
	client       *mongo.Client
	codec        runtime.Codec
	watcher      *watcher
}

// Ensure MongoStore implements store Interface
var _ Interface = &MongoStore{}

// NewMongoStore Mongo Store constructor
func NewMongoStore(client *mongo.Client, codec runtime.Codec, databaseName string, newFunc func() runtime.Object) *MongoStore {
	st := &MongoStore{
		client:       client,
		codec:        codec,
		databaseName: databaseName,
		watcher:      newWatcher(client, codec, newFunc, databaseName),
	}
	return st
}

// Get finds a single entry in the store and populates the "into" parameter
func (m *MongoStore) Get(ctx context.Context, key string, options GetOptions, into runtime.Object) error {
	// retrieve mongo collection info
	info, err := getMongoObjectInfo(key)
	if err != nil {
		return err
	}

	collection := m.client.Database(m.databaseName).Collection(info.collection)

	_, err = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			// Create index for the "key"
			Options: mongooptions.Index().SetName("key"),
			Keys: bson.D{
				{Key, 1},
				{IsCurrentKey, 1},
				{IsDeletedKey, 1},
			},
		}, {
			// Create index for FIND (list) queries
			Options: mongooptions.Index().SetName("isCurrent_isDeleted"),
			Keys: bson.D{
				{IsCurrentKey, 1},
				{IsDeletedKey, 1},
			},
		}, {
			// Create index for labels
			// Allows for quick database-enabled ABAC filtering on labels
			Options: mongooptions.Index().SetName("metadata"),
			Keys: bson.D{
				{"metadata.labels.$**", 1},
			},
		},
	})
	if err != nil {
		return err
	}

	// find mongo record
	result := collection.FindOne(ctx, bson.M{
		IsDeletedKey: false,
		IsCurrentKey: true,
		Key:          info.key,
	})
	if result.Err() != nil {
		return interpretMongoError(result.Err())
	}

	// map record
	objInfo, err := transformMongoSingleResult(result)
	if err != nil {
		return fmt.Errorf("unable to transform mongo SingleResult: %v", err)
	}

	// decode and return
	return decode(m.codec, objInfo.CurrentBytes, into, objInfo.Timestamp)

}

// Create inserts a new object in the store. If the "out" parameter is specified, the method will also attempt
// to return the created object
func (m *MongoStore) Create(ctx context.Context, key string, options CreateOptions, obj, out runtime.Object) error {

	info, err := getMongoObjectInfo(key)
	if err != nil {
		return err
	}

	session, err := m.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// It is possible that a previous record with the same key was deleted.
	// We should not prevent the user from creating a new record with the same key
	// This method will
	// 1. Try to find an existing record with the same key
	// 2. If there was no record, insert the new one we want to create
	// 3. If the previous record with the same key was not marked as deleted, throw an exception for key conflict
	// 4. If the previous record with the same key was deleted
	//    1. Insert the new record we want to create, and link the previousTimestamp, previousId
	//   	 and creationTimestamp of the old record
	//    2. Update the old record by setting the isCurrent, nextTimestamp, nextId
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {

		// Find the latest record with the same key
		existingResult := sessCtx.
			Client().
			Database(m.databaseName).
			Collection(info.collection).
			FindOne(ctx, bson.M{
				Key:          info.key,
				IsCurrentKey: true,
			})

		// We have an error which is not a "no documents in result"
		if existingResult.Err() != nil && existingResult.Err().Error() != "mongo: no documents in result" {
			return nil, interpretMongoError(existingResult.Err())
		}

		var oldID *primitive.ObjectID
		var oldObj runtime.Object
		var oldTimestamp *int64
		var initialCreationTimestamp int64
		var hasPreviousRecord = false
		var oldRecord *MongoRecord

		if existingResult.Err() != nil {
			oldID = nil
			oldObj = nil
			oldTimestamp = nil
			initialCreationTimestamp = time.Now().UTC().UnixNano()
		} else {
			// We don't have a "no documents in result" error. Meaning there was
			// a record with the same key in the database

			hasPreviousRecord = true

			// Retrieve the previous record info
			oldRecord, err = transformMongoSingleResult(existingResult)
			if err != nil {
				return nil, err
			}

			// Decode the previous record
			var oldObject = out.DeepCopyObject()
			if err := decode(m.codec, oldRecord.CurrentBytes, oldObject, oldRecord.Timestamp); err != nil {
				return nil, err
			}

			// Make sure that the record was marked as deleted, otherwise
			// throw a key conflict error
			if !oldRecord.IsDeleted {
				return nil, fmt.Errorf("key already exists")
			}

			oldID = &oldRecord.ID
			oldObj = oldObject
			oldTimestamp = &oldRecord.Timestamp
			initialCreationTimestamp = oldRecord.CreationTimestamp
		}

		// This is the timestamp for the current version
		timestamp := time.Now().UTC().UnixNano()

		// Transform the new record to mongo storage
		storedObj, err := transformRuntimeObjectToMongoStorage(
			m.codec,
			initialCreationTimestamp,
			oldID,
			oldObj,
			oldTimestamp,
			&timestamp,
			obj,
			true,
			false,
			true)

		if err != nil {
			return nil, err
		}

		collection := m.client.Database(m.databaseName).Collection(info.collection)

		// Insert the new record
		insertResult, err := collection.InsertOne(ctx, storedObj)
		if err != nil {
			return nil, interpretMongoError(err)
		}

		if hasPreviousRecord {

			// Update the old record, by making sure that it was not
			// updated between the time we pulled it and now.
			updateResult, err := collection.UpdateOne(ctx, bson.M{
				Key:          oldRecord.Key,
				IsCurrentKey: true,
				TimestampKey: oldRecord.Timestamp,
				IsDeletedKey: true,
			}, bson.M{
				"$set": bson.M{
					NextIDKey:        insertResult.InsertedID,
					NextTimestampKey: &timestamp,
					IsCurrentKey:     false,
				},
			})
			if err != nil {
				return nil, err
			}

			// Make sure we actually updated it. Otherwise that means that
			// the record was updated by another process between the time we
			// got it from the db and now. Abort the transaction then.
			if updateResult.MatchedCount == 0 {
				return nil, fmt.Errorf("couldn't match previously deleted record for update")
			}
			if updateResult.ModifiedCount == 0 {
				return nil, fmt.Errorf("couldn't update previously deleted record for update")
			}
		}

		if out != nil {

			findResult := collection.FindOne(ctx, bson.M{"_id": insertResult.InsertedID})
			if findResult.Err() != nil {
				return nil, interpretMongoError(findResult.Err())
			}

			resultRecord, err := transformMongoSingleResult(findResult)
			if err != nil {
				return nil, err
			}

			if err := decode(m.codec, resultRecord.CurrentBytes, out, resultRecord.Timestamp); err != nil {
				return nil, err
			}

		}

		return nil, nil

	},
		mongooptions.
			Transaction().
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
	)
	if err != nil {
		return interpretMongoError(err)
	}

	return nil
}

// List queries the database and returns a list of records matching the query
func (m *MongoStore) List(ctx context.Context, key string, options ListOptions, listObj runtime.Object) error {

	// Gets the pointer to the Items property
	itemsPtr, err := meta.GetItemsPtr(listObj)
	if err != nil {
		return err
	}

	// Get the value of the Items property
	itemsVal, err := conversion.EnforcePtr(itemsPtr)
	if err != nil {
		return fmt.Errorf("unable to convert list object to pointer: %v", err)
	}

	// Retrieve which collection we're gonna list from
	info, err := getMongoCollectionInfo(key)
	if err != nil {
		return err
	}

	// Build query options
	var findOptions []*mongooptions.FindOptions
	if options.Limit != nil && *options.Limit != 0 {
		// ListOptions has a limit set
		findOptions = append(findOptions, mongooptions.Find().SetLimit(*options.Limit))
	}

	// Build query filter
	var filter interface{}
	filter = bson.M{
		IsCurrentKey:     true,
		IsDeletedKey:     false,
		NextIDKey:        nil,
		NextTimestampKey: nil,
	}

	// Merge the query filter with the provided selector using an "$and" expression
	if options.Selector != nil {
		query, err := convertFieldSelectorToMongoFilter(options.Selector, CurrentValueKey+".")
		if err != nil {
			return err
		}
		// query might be nil, for example if it is an $and with no items in it
		if query != nil {
			filter = bson.M{
				"$and": bson.A{
					filter,
					query,
				},
			}
		}
	}

	// Query the results
	result, err := m.client.Database(m.databaseName).Collection(info.collectionName).Find(ctx, filter, findOptions...)
	if err != nil {
		return interpretMongoError(err)
	}

	if result.Err() != nil {
		return interpretMongoError(result.Err())
	}

	// a function used to instantiate new items of the given list
	newListItemFunc := getNewItemFunc(listObj, itemsVal)

	// iterate through all items
	for {

		// stop if there are no more results
		if !result.Next(ctx) {
			break
		}

		// the result.Err() might be non nil at that point
		if result.Err() != nil {
			return interpretMongoError(result.Err())
		}

		// instantiate new list item
		newItem := newListItemFunc()

		// convert mongo result to runtime.Object
		if err := transformBsonRawToRuntimeObject(m.codec, result.Current, newItem); err != nil {
			return err
		}

		// append item to list
		itemsVal.Set(reflect.Append(itemsVal, reflect.ValueOf(newItem).Elem()))

	}

	return nil
}

func (m *MongoStore) Update(ctx context.Context, key string, options UpdateOptions, out runtime.Object, tryUpdate UpdateFunc) error {

	mongoObjectInfo, err := getMongoObjectInfo(key)
	if err != nil {
		return err
	}

	_, err = conversion.EnforcePtr(out)
	if err != nil {
		return fmt.Errorf("unable to convert output object to pointer: %v", err)
	}

	// Find the current value for the object to be updated
	objToUpdateResult := m.client.Database(m.databaseName).Collection(mongoObjectInfo.collection).FindOne(ctx, bson.M{
		Key:              mongoObjectInfo.key,
		IsDeletedKey:     false,
		IsCurrentKey:     true,
		NextTimestampKey: nil,
		NextIDKey:        nil,
	})
	if objToUpdateResult.Err() != nil {
		return interpretMongoError(objToUpdateResult.Err())
	}

	// Convert to bytes
	objToUpdateInfo, err := transformMongoSingleResult(objToUpdateResult)
	if err != nil {
		return err
	}

	// Decode mongo response bytes
	objToUpdate := out.DeepCopyObject()
	if err := decode(m.codec, objToUpdateInfo.CurrentBytes, objToUpdate, objToUpdateInfo.Timestamp); err != nil {
		return err
	}

	// Keep a copy of the current object
	// We copy it because the tryUpdate function might mutate
	// the object
	objToUpdateCopy := objToUpdate.DeepCopyObject()

	// Update the object. Will perhaps be mutated.
	updatedObject, err := tryUpdate(objToUpdate)
	if err != nil {
		return err
	}

	updateTimestamp := time.Now().UTC().UnixNano()

	// Convert the updated object to mongo storage
	// for persistence
	storageValue, err := transformRuntimeObjectToMongoStorage(
		m.codec,
		objToUpdateInfo.CreationTimestamp,
		&objToUpdateInfo.ID,
		objToUpdateCopy,
		&objToUpdateInfo.Timestamp,
		&updateTimestamp,
		updatedObject,
		true,
		false,
		false)

	if err != nil {
		return err
	}

	session, err := m.client.StartSession()
	if err != nil {
		return interpretMongoError(err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(
		ctx,
		func(sessCtx mongo.SessionContext) (interface{}, error) {

			collection := sessCtx.Client().Database(m.databaseName).Collection(mongoObjectInfo.collection)

			insertResult, err := collection.InsertOne(ctx, storageValue)
			if err != nil {
				return nil, err
			}

			result, err := collection.UpdateOne(ctx, bson.M{
				IdKey:            objToUpdateInfo.ID,
				Key:              mongoObjectInfo.key,
				IsCurrentKey:     true,
				TimestampKey:     objToUpdateInfo.Timestamp,
				NextTimestampKey: nil,
				NextIDKey:        nil,
			}, bson.M{
				"$set": bson.M{
					IsCurrentKey:     false,
					NextTimestampKey: updateTimestamp,
					NextIDKey:        insertResult.InsertedID,
				},
			})
			if err != nil {
				return nil, err
			}

			if result.MatchedCount == 0 {
				return nil, fmt.Errorf("could not find object to update: no records matched")
			}
			if result.ModifiedCount == 0 {
				return nil, fmt.Errorf("could not find object to update: no records modified")
			}

			return nil, nil

		}, mongooptions.Transaction().SetWriteConcern(writeconcern.New(writeconcern.WMajority())))
	if err != nil {
		return interpretMongoError(err)
	}

	return m.Get(ctx, key, GetOptions{}, out)

}

func (m *MongoStore) Delete(ctx context.Context, key string, options DeleteOptions, out runtime.Object) error {

	_, err := conversion.EnforcePtr(out)
	if err != nil {
		return fmt.Errorf("unable to convert output object to pointer: %v", err)
	}

	mongoObjectInfo, err := getMongoObjectInfo(key)
	if err != nil {
		return err
	}

	session, err := m.client.StartSession()
	if err != nil {
		return interpretMongoError(err)
	}

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {

		collection := sessCtx.Client().Database(m.databaseName).Collection(mongoObjectInfo.collection)

		// Find the object to be deleted
		findResult := collection.FindOne(ctx, bson.M{
			IsCurrentKey:     true,
			Key:              mongoObjectInfo.key,
			IsDeletedKey:     false,
			NextIDKey:        nil,
			NextTimestampKey: nil,
		})
		if findResult.Err() != nil {
			return nil, interpretMongoError(err)
		}

		// Convert the mongo record
		mongoRecord, err := transformMongoSingleResult(findResult)
		if err != nil {
			return nil, err
		}

		deletionTimestamp := time.Now().UTC().UnixNano()

		// This is the record representing the deletion
		newMongoRecord := MongoRecord{
			Key:               mongoObjectInfo.key,
			CurrentValue:      nil,
			PreviousValue:     mongoRecord.CurrentValue,
			IsDeleted:         true,
			IsCreated:         false,
			IsCurrent:         true,
			Timestamp:         deletionTimestamp,
			CreationTimestamp: mongoRecord.CreationTimestamp,
			DeletionTimestamp: &deletionTimestamp,
			PreviousTimestamp: &mongoRecord.Timestamp,
			NextTimestamp:     nil,
			PreviousID:        &mongoRecord.ID,
		}

		// Insert the record representing the object deletion
		deleteResult, err := collection.InsertOne(ctx, newMongoRecord)
		if err != nil {
			return nil, interpretMongoError(err)
		}

		// Update the current record to point to the new version (deleted)
		updateResult, err := collection.UpdateOne(ctx, bson.M{
			IsCurrentKey:     true,
			IsDeletedKey:     false,
			Key:              mongoObjectInfo.key,
			IdKey:            mongoRecord.ID,
			NextTimestampKey: nil,
			NextIDKey:        nil,
		}, bson.M{
			"$set": bson.M{
				IsCurrentKey:     false,
				NextTimestampKey: deletionTimestamp,
				NextIDKey:        deleteResult.InsertedID,
			},
		})
		if err != nil {
			return nil, interpretMongoError(err)
		}

		// Ensure we actually modified something
		if updateResult.MatchedCount == 0 || updateResult.ModifiedCount == 0 {
			return nil, fmt.Errorf("could not delete record: could not find previous record")
		}

		// Retrieve the bytes corresponding to the deleted object
		deletedObj, err := json.Marshal(mongoRecord.CurrentValue)
		if err != nil {
			return nil, err
		}

		// Decode the deleted obj into the out param
		if err := decode(m.codec, deletedObj, out, deletionTimestamp); err != nil {
			return nil, err
		}

		return nil, nil

	},
		mongooptions.Transaction().SetReadConcern(readconcern.Linearizable()),
		mongooptions.Transaction().SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
	)

	// transaction failed
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoStore) Watch(ctx context.Context, key string, options ListOptions) (watch.Interface, error) {
	rev, err := parseResourceVersion(options.ResourceVersion)
	if err != nil {
		return nil, err
	}

	watchChan, err := m.watcher.Watch(
		ctx,
		key,
		rev,
		true,
		true,
		options.SyncOnly,
		options.Selector,
		options.Limit)

	if err != nil {
		return nil, err
	}
	return watchChan, nil
}

func parseResourceVersion(v string) (int64, error) {
	if v == "" {
		return 0, nil
	}
	rev, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}
	return rev, nil
}

// transformBsonRawToRuntimeObject uses the bson.Raw value to populate the given runtime.Object
func transformBsonRawToRuntimeObject(codec runtime.Codec, bsonRaw bson.Raw, into runtime.Object) error {

	_, err := conversion.EnforcePtr(into)
	if err != nil {
		return fmt.Errorf("unable to convert obj to pointer: %v", err)
	}

	var tempMap = map[string]interface{}{}

	if err := bson.Unmarshal(bsonRaw, tempMap); err != nil {
		return interpretMongoError(err)
	}

	currentValueIntf, ok := tempMap[CurrentValueKey]
	if !ok {
		return fmt.Errorf("could not get current value")
	}

	currentValue, ok := currentValueIntf.(map[string]interface{})
	if !ok {
		return fmt.Errorf("could not get current value")
	}

	us := &unstructured.Unstructured{}
	us.Object = currentValue

	buf := bytes.NewBuffer(nil)
	if err := codec.Encode(us, buf); err != nil {
		return err
	}

	_, _, err = codec.Decode(buf.Bytes(), nil, into)
	if err != nil {
		return err
	}

	return nil
}

type mongoStoredObjInfo struct {
	bytes                    []byte
	initialCreationTimestamp *time.Time
	timestamp                *time.Time
	deletionTimestamp        *time.Time
	isDeleted                bool
	objID                    *primitive.ObjectID
}

type MongoRecord struct {
	ID                primitive.ObjectID     `bson:"_id,omitempty"`
	Key               string                 `bson:"key"`
	CurrentValue      map[string]interface{} `bson:"currentValue"`
	CurrentBytes      []byte                 `bson:"-"`
	PreviousValue     map[string]interface{} `bson:"previousValue"`
	PreviousBytes     []byte                 `bson:"-"`
	IsDeleted         bool                   `bson:"isDeleted"`
	IsCreated         bool                   `bson:"isCreated"`
	IsCurrent         bool                   `bson:"isCurrent"`
	Timestamp         int64                  `bson:"timestamp"`
	CreationTimestamp int64                  `bson:"creationTimestamp"`
	DeletionTimestamp *int64                 `bson:"deletionTimestamp"`
	PreviousTimestamp *int64                 `bson:"previousTimestamp"`
	NextTimestamp     *int64                 `bson:"nextTimestamp"`
	PreviousID        *primitive.ObjectID    `bson:"previousId"`
	NextID            *primitive.ObjectID    `bson:"nextId"`
}

// transformMongoSingleResult will transform a mongo SingleResult to a byte array that can be
// used to decode a runtime.Object
func transformMongoSingleResult(result *mongo.SingleResult) (*MongoRecord, error) {

	// Expect non-nil SingleResult
	if result == nil {
		return nil, fmt.Errorf("unexpected nil mongo SingleResult")
	}

	record := MongoRecord{}
	if err := result.Decode(&record); err != nil {
		return nil, err
	}

	currentBytes, err := json.Marshal(record.CurrentValue)
	if err != nil {
		return nil, err
	}
	record.CurrentBytes = currentBytes

	if record.PreviousValue != nil && len(record.PreviousValue) != 0 {
		previousBytes, err := json.Marshal(record.PreviousValue)
		if err != nil {
			return nil, err
		}
		record.PreviousBytes = previousBytes
	}

	return &record, nil

}

// transformRuntimeObjectToMongoStorage takes a runtime.Object and converts it to a format that
// can be persisted to mongo
func transformRuntimeObjectToMongoStorage(
	codec runtime.Codec,
	initialCreationTimestamp int64,
	oldID *primitive.ObjectID,
	oldObj runtime.Object,
	oldTimestamp *int64,
	currentTimestamp *int64,
	obj runtime.Object,
	isCurrent,
	isDeleted,
	isCreated bool,
) (interface{}, error) {

	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, fmt.Errorf("could not get meta accessor for object: %v", err)
	}

	objMap, err := transformRuntimeObjectToMap(codec, obj)
	if err != nil {
		return nil, err
	}

	var oldMap map[string]interface{}
	if oldObj != nil {
		if oldID == nil {
			return nil, fmt.Errorf("oldObj object id must be supplied when oldObj is not nil")
		}
		oldMap, err = transformRuntimeObjectToMap(codec, oldObj)
		if err != nil {
			return nil, err
		}
	}

	namespace := accessor.GetNamespace()
	name := accessor.GetName()

	var nullDateTime *time.Time

	var out = map[string]interface{}{}

	// The current value of the object
	out[CurrentValueKey] = objMap

	// The key of the object
	out[Key] = path.Join(namespace, name)

	// Is the object the current version
	out[IsCurrentKey] = isCurrent

	// Is the object deleted
	out[IsDeletedKey] = isDeleted

	// Is the object created
	out[IsCreatedKey] = isCreated

	// The insertion timestamp of that version
	out[TimestampKey] = currentTimestamp

	// The deletion timestamp of the object
	out[DeletionTimestampKey] = nullDateTime

	// The initial creation timestamp of that object
	// An object might store multiple versions, this version
	// stores the initial timestamp of the initial object
	out[CreationTimestampKey] = initialCreationTimestamp

	// The timestamp of the next version, if any
	out[NextTimestampKey] = nullDateTime

	if oldObj != nil {

		// The value of the previous version
		out[PreviousValueKey] = oldMap

		// The previous object _id
		out[PreviousIDKey] = *oldID

		// The previous object timestamp
		out[PreviousTimestampKey] = oldTimestamp
	}
	return out, nil

}

func transformRuntimeObjectToMap(codec runtime.Codec, obj runtime.Object) (map[string]interface{}, error) {
	buf := bytes.NewBuffer(nil)
	if err := codec.Encode(obj, buf); err != nil {
		return nil, err
	}

	us := unstructured.Unstructured{}
	_, _, err := codec.Decode(buf.Bytes(), nil, &us)
	if err != nil {
		return nil, err
	}

	return us.Object, nil
}

// interpretMongoError will convert a mongo exception to an helpful exception for the api
// for example, a mongo "not found" error will be converted to the proper StatusError with
// exception NotFound
func interpretMongoError(err error) error {
	// TODO: map mongo errors
	return err
}

// mongoObjectInfo is a helper struct that contains the information to locate
// a single record in a mongo database
type mongoObjectInfo struct {
	collection string
	key        string
}

// getMongoObjectInfo maps a given object key to its mongo storage position
// the key should be in the form {group}/{resource (plural)}/id
// the collection name will be in the form of {group}__{resource (plural)}
func getMongoObjectInfo(key string) (mongoObjectInfo, error) {
	// key in the form of {group}/{resource}/{id}
	parts := strings.Split(key, "/")
	if len(parts) != 3 {
		return mongoObjectInfo{}, fmt.Errorf("invalid key format. expecting {group}/{resource}/{id} format")
	}
	collectionName := parts[0] + "__" + parts[1]
	return mongoObjectInfo{
		collection: collectionName,
		key:        parts[2],
	}, nil
}

type mongoCollectionInfo struct {
	collectionName string
}

// getMongoCollectionInfo maps a given key to the mongo storage collection
// the key should be in the form {group}/{resource (plural)}
// the collection name will be in the form of {group}__{resource (plural)}
func getMongoCollectionInfo(key string) (mongoCollectionInfo, error) {
	parts := strings.Split(key, "/")
	if len(parts) < 2 {
		return mongoCollectionInfo{}, fmt.Errorf("invalid key format. expecting {group}/{resource} format")
	}
	collectionName := parts[0] + "__" + parts[1]
	return mongoCollectionInfo{
		collectionName: collectionName,
	}, nil
}

// decode takes a byte array and tries to convert it to the given runtime.Object specified in the "into"
// parameter
func decode(codec runtime.Codec, value []byte, into runtime.Object, rev int64) error {
	if _, err := conversion.EnforcePtr(into); err != nil {
		return fmt.Errorf("unable to convert object to pointer: %v", err)
	}
	if _, _, err := codec.Decode(value, nil, into); err != nil {
		return err
	}
	accessor, err := meta.Accessor(into)
	if err != nil {
		return err
	}
	accessor.SetResourceVersion(strconv.FormatInt(rev, 10))
	return nil
}

// getNewItemFunc will inspect the given listObj, determine what type of item it contains and return a
// function that can be used to instantiate new instances of that object. Used by the MongoStore.List function
// to populate the response list
func getNewItemFunc(listObj runtime.Object, v reflect.Value) func() runtime.Object {
	// For unstructured lists with a target group/version, preserve the group/version in the instantiated list items
	if unstructuredList, isUnstructured := listObj.(*unstructured.UnstructuredList); isUnstructured {
		if apiVersion := unstructuredList.GetAPIVersion(); len(apiVersion) > 0 {
			return func() runtime.Object {
				return &unstructured.Unstructured{Object: map[string]interface{}{"apiVersion": apiVersion}}
			}
		}
	}

	// Otherwise just instantiate an empty item
	elem := v.Type().Elem()
	return func() runtime.Object {
		return reflect.New(elem).Interface().(runtime.Object)
	}
}

func convertFieldSelectorToMongoFilter(selector fields2.Selector, propertyPrefix string) (bson.M, error) {
	if selector == nil {
		return nil, nil
	}
	switch s := selector.(type) {
	case fields2.AndSelector:
		var filterItems = bson.A{}
		for _, item := range s {
			filterItem, err := convertFieldSelectorToMongoFilter(item, propertyPrefix)
			if err != nil {
				return nil, err
			}
			filterItems = append(filterItems, filterItem)
		}
		if len(filterItems) == 0 {
			return nil, nil
		}
		return bson.M{
			"$and": filterItems,
		}, nil
	case fields2.OrSelector:
		var filterItems = bson.A{}
		for _, item := range s {
			filterItem, err := convertFieldSelectorToMongoFilter(item, propertyPrefix)
			if err != nil {
				return nil, err
			}
			filterItems = append(filterItems, filterItem)
		}
		if len(filterItems) == 0 {
			return nil, nil
		}
		return bson.M{
			"$or": filterItems,
		}, nil
	case fields2.EqualSelector:
		if len(s.Key) == 0 {
			return nil, fmt.Errorf("unexpected property key")
		}
		return bson.M{
			propertyPrefix + s.Key: bson.M{"$eq": s.Value},
		}, nil
	default:
		return nil, fmt.Errorf("unexpected selector type: %#v", s)
	}
}
