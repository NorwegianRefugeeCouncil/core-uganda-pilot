package store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	"path"
	"reflect"
	"strings"
)

// Interface represents the generic store interface
type Interface interface {
	Get(ctx context.Context, key string, options GetOptions, into runtime.Object) error
	Create(ctx context.Context, key string, options CreateOptions, obj, out runtime.Object) error
	List(ctx context.Context, key string, options ListOptions, listObj runtime.Object) error
	Update(ctx context.Context, key string, options UpdateOptions, out runtime.Object, tryUpdate UpdateFunc) error
	Delete(ctx context.Context, key string, options DeleteOptions, out runtime.Object) error
}

// GetOptions placeholder to put store options for GET requests that return a single result
type GetOptions struct{}

// CreateOptions placeholder to put store options for PUT requests
type CreateOptions struct{}

// ListOptions placeholder to put store options for GET requests that return a list result
type ListOptions struct{}

// UpdateOptions placeholder to put store options for PUT request
type UpdateOptions struct{}

// UpdateOptions placeholder to put store options for PUT request
type DeleteOptions struct{}

type UpdateFunc func(input runtime.Object) (output runtime.Object, err error)

// MongoStore mongo store implementation
type MongoStore struct {
	databaseName string
	client       *mongo.Client
	codec        runtime.Codec
}

// Ensure MongoStore implements store Interface
var _ Interface = &MongoStore{}

// NewMongoStore Mongo Store constructor
func NewMongoStore(client *mongo.Client, codec runtime.Codec, databaseName string) *MongoStore {
	st := &MongoStore{
		client:       client,
		codec:        codec,
		databaseName: databaseName,
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

	// find mongo record
	result := m.client.Database(m.databaseName).Collection(info.collection).FindOne(ctx, bson.M{
		"__key": info.key,
	})
	if result.Err() != nil {
		return interpretMongoError(result.Err())
	}

	// map record
	objBytes, err := transformMongoSingleResult(result)
	if err != nil {
		return fmt.Errorf("unable to transform mongo SingleResult: %v", err)
	}

	// decode and return
	return decode(m.codec, objBytes, into)

}

// Create inserts a new object in the store. If the "out" parameter is specified, the method will also attempt
// to return the created object
func (m *MongoStore) Create(ctx context.Context, key string, options CreateOptions, obj, out runtime.Object) error {

	info, err := getMongoObjectInfo(key)
	if err != nil {
		return err
	}

	storedObj, err := transformRuntimeObjectToMongoStorage(m.codec, obj)
	if err != nil {
		return err
	}

	_, err = m.client.Database(m.databaseName).Collection(info.collection).InsertOne(ctx, storedObj)
	if err != nil {
		return interpretMongoError(err)
	}

	if out != nil {
		return m.Get(ctx, key, GetOptions{}, out)
	}

	return nil
}

// List queries the database and returns a list of records matching the query
func (m *MongoStore) List(ctx context.Context, key string, options ListOptions, listObj runtime.Object) error {

	listPtr, err := meta.GetItemsPtr(listObj)
	if err != nil {
		return err
	}

	listVal, err := conversion.EnforcePtr(listPtr)
	if err != nil {
		return fmt.Errorf("unable to convert list object to pointer: %v", err)
	}

	info, err := getMongoCollectionInfo(key)
	if err != nil {
		return err
	}

	// TODO: filter out results based on some query
	result, err := m.client.Database(m.databaseName).Collection(info.collectionName).Find(ctx, bson.M{})
	if err != nil {
		return interpretMongoError(err)
	}

	if result.Err() != nil {
		return interpretMongoError(result.Err())
	}

	// a function used to instantiate new items of the given list
	newListItemFunc := getNewItemFunc(listObj, listVal)

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
		listVal.Set(reflect.Append(listVal, reflect.ValueOf(newItem).Elem()))

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

	objCopy := out.DeepCopyObject()
	if err := m.Get(ctx, key, GetOptions{}, objCopy); err != nil {
		return err
	}

	updated, err := tryUpdate(objCopy)
	if err != nil {
		return err
	}

	storageValue, err := transformRuntimeObjectToMongoStorage(m.codec, updated)
	if err != nil {
		return err
	}

	_, err = m.client.Database(m.databaseName).Collection(mongoObjectInfo.collection).ReplaceOne(ctx, bson.M{
		"__key": mongoObjectInfo.key,
	}, storageValue)
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

	if err := m.Get(ctx, key, GetOptions{}, out); err != nil {
		return err
	}

	_, err = m.client.Database(m.databaseName).Collection(mongoObjectInfo.collection).DeleteOne(ctx, bson.M{
		"__key": mongoObjectInfo.key,
	})
	if err != nil {
		return interpretMongoError(err)
	}

	return nil
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

	us := &unstructured.Unstructured{}
	us.Object = tempMap

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

// transformMongoSingleResult will transform a mongo SingleResult to a byte array that can be
// used to decode a runtime.Object
func transformMongoSingleResult(result *mongo.SingleResult) ([]byte, error) {

	if result == nil {
		return nil, fmt.Errorf("unexpected nil mongo SingleResult")
	}

	var tempMap = map[string]interface{}{}

	if err := result.Decode(&tempMap); err != nil {
		return nil, fmt.Errorf("unable to decode mongo result into map[string]interface{}: %v", err)
	}

	jsonBytes, err := json.Marshal(tempMap)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal map[string]interface{}: %v", err)
	}

	return jsonBytes, nil

}

// transformRuntimeObjectToMongoStorage takes a runtime.Object and converts it to a format that
// can be persisted to mongo
func transformRuntimeObjectToMongoStorage(codec runtime.Codec, obj runtime.Object) (interface{}, error) {

	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, fmt.Errorf("could not get meta accessor for object: %v", err)
	}

	buf := bytes.NewBuffer(nil)
	if err := codec.Encode(obj, buf); err != nil {
		return nil, err
	}

	us := unstructured.Unstructured{}
	_, _, err = codec.Decode(buf.Bytes(), nil, &us)
	if err != nil {
		return nil, err
	}

	out := us.Object
	namespace := accessor.GetNamespace()
	name := accessor.GetName()
	out["__key"] = path.Join(namespace, name)
	return out, nil

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
	if len(parts) != 2 {
		return mongoCollectionInfo{}, fmt.Errorf("invalid key format. expecting {group}/{resource} format")
	}
	collectionName := parts[0] + "__" + parts[1]
	return mongoCollectionInfo{
		collectionName: collectionName,
	}, nil
}

// decode takes a byte array and tries to convert it to the given runtime.Object specified in the "into"
// parameter
func decode(codec runtime.Codec, value []byte, into runtime.Object) error {
	if _, err := conversion.EnforcePtr(into); err != nil {
		return fmt.Errorf("unable to convert object to pointer: %v", err)
	}
	if _, _, err := codec.Decode(value, nil, into); err != nil {
		return err
	}
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
