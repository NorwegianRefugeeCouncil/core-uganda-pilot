package mongo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nrc-no/core/apps/api/pkg/apis/meta"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/storage"
	"github.com/nrc-no/core/apps/api/pkg/util/conversion"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"strings"
)

type Store struct {
	mongoClient *mongo.Client
	database    string
	collection  string
	create      func() runtime.Object
	codec       runtime.Codec
	versioner   storage.Versioner
}

type objState struct {
	obj   runtime.Object
	meta  *storage.ResponseMeta
	rev   int64
	data  []byte
	stale bool
}

func NewStore(
	mongoClient *mongo.Client,
	codec runtime.Codec,
	create func() runtime.Object,
	prefix string,
) (*Store, error) {

	if strings.HasPrefix(prefix, "/") {
		prefix = prefix[1:]
	}
	if strings.HasSuffix(prefix, "/") {
		prefix = prefix[:1]
	}
	parts := strings.Split(prefix, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expecting format of <database>/<collection> from 'prefix' argument")
	}
	database := parts[0]
	collection := parts[1]

	if len(database) == 0 {
		return nil, fmt.Errorf("database is required")
	}
	if len(collection) == 0 {
		return nil, fmt.Errorf("collection is required")
	}

	return &Store{
		mongoClient: mongoClient,
		database:    database,
		collection:  collection,
		create:      create,
		codec:       codec,
	}, nil
}

var _ storage.Interface = &Store{}

func (s *Store) Get(ctx context.Context, key string, getOptions storage.GetOptions, out runtime.Object) error {

	objectID, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return err
	}

	collection := s.mongoClient.Database(s.database).Collection(s.collection)
	result := collection.FindOne(ctx, bson.M{"_id": objectID})
	if result.Err() != nil {
		return result.Err()
	}

	bytes, err := result.DecodeBytes()
	if err != nil {
		return err
	}

	return convertDocument(bytes, out)

}

func (s *Store) List(ctx context.Context, listOptions storage.ListOptions, out runtime.Object) error {

	listPtr, err := GetItemsPtr(out)
	if err != nil {
		return err
	}

	v, err := conversion.EnforcePtr(listPtr)
	if err != nil || v.Kind() != reflect.Slice {
		return fmt.Errorf("need ptr to slice: %v", err)
	}

	newItemFunc := getNewItemFunc(out, v)

	result, err := s.mongoClient.Database(s.database).Collection(s.collection).Find(ctx, bson.D{})
	if err != nil {
		return err
	}

	var objects []runtime.Object

	for {
		if !result.Next(ctx) {
			break
		}
		current := result.Current
		item := newItemFunc()
		if err := convertDocument(current, item); err != nil {
			return err
		}
		objects = append(objects, item)
	}

	for _, object := range objects {
		if err := appendListItem(v, object); err != nil {
			return err
		}
	}

	return nil

}

func convertDocument(data bson.Raw, into runtime.Object) error {

	var resourceVersion int
	if err := data.Lookup("resourceVersion").Unmarshal(&resourceVersion); err != nil {
		logrus.Errorf("unable to get resourceVersion from document: %v", err)
		return err
	}

	raw := data.Lookup("currentRevision")
	if err := raw.Unmarshal(into); err != nil {
		logrus.Errorf("unable to get currentRevision from document: %v", err)
		return err
	}

	var objectID primitive.ObjectID
	if err := objectID.UnmarshalJSON(data.Lookup("_id").Value); err != nil {
		logrus.Errorf("unable to unmarshal objectID from document: %v", err)
		return err
	}

	accessor, err := meta.Accessor(into)
	if err != nil {
		return err
	}

	accessor.SetUID(objectID.Hex())
	accessor.SetResourceVersion(resourceVersion)

	return nil

}

func (s *Store) Create(ctx context.Context, obj runtime.Object) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	accessor.SetResourceVersion(1)

	response, err := s.mongoClient.Database(s.database).Collection(s.collection).InsertOne(ctx, obj)
	if err != nil {
		return err
	}

	objectID := response.InsertedID.(primitive.ObjectID)
	id := objectID.Hex()
	accessor.SetUID(id)

	return nil
}

func (s *Store) Versioner() storage.Versioner {
	return s.versioner
}

func (s *Store) Update(ctx context.Context, key string, out runtime.Object, updateFunc storage.UpdateFunc) error {

	objectID, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return err
	}

	v, err := conversion.EnforcePtr(out)
	if err != nil {
		return err
	}

	getCurrentState := func() (*objState, error) {
		getResp := s.mongoClient.Database(s.database).Collection(s.collection).FindOne(ctx, bson.M{"_id": key})
		if getResp.Err() != nil {
			return nil, getResp.Err()
		}
		return s.getState(getResp, key, v, false)
	}

	var origState *objState
	var origStateIsCurrent bool

	origState, err = getCurrentState()
	origStateIsCurrent = true

	if err != nil {
		return err
	}

	for {

		ret, _, err := s.updateState(origState, updateFunc)
		if err != nil {
			// if data is already up to date, return the error
			if origStateIsCurrent {
				return err
			}

			// refresh id data is stale
			origState, err = getCurrentState()
			if err != nil {
				return err
			}
			origStateIsCurrent = true
			//retry
			continue
		}

		data, err := runtime.Encode(s.codec, ret)
		if err != nil {
			return err
		}

		if !origState.stale && bytes.Equal(data, origState.data) {
			// if we skipped the original get in the loop, we must refresh
			// in order to be sure the data in the store is equivalent
			// to our desired serialization
			if !origStateIsCurrent {
				origState, err := getCurrentState()
				if err != nil {
					return err
				}
				origStateIsCurrent = true
				if !bytes.Equal(data, origState.data) {
					// original data changed, restart loop
					continue
				}
			}
			// recheck that the data is not stale before short-circuiting a write
			if !origState.stale {
				return decode(s.codec, s.versioner, origState.data, out, origState.rev)
			}
		}

		var temp = map[string]interface{}{}
		if err := json.Unmarshal(data, &temp); err != nil {
			return err
		}
		temp["__revision"] = origState.rev

		mongoData, err := bson.Marshal(temp)
		if err != nil {
			return err
		}

		_, err = s.mongoClient.Database(s.database).Collection(s.collection).UpdateOne(
			ctx,
			bson.M{
				"_id":        objectID,
				"__revision": origState.rev,
			},
			mongoData,
		)
		if err != nil {
			return err
		}

		return decode(s.codec, s.versioner, mongoData, out, origState.rev)

	}

	//
	//_, err = s.mongoClient.Database(s.database).Collection(s.collection).UpdateByID(ctx, objectID, bson.A{
	//  bson.M{"$set": bson.M{"previousRevisions": bson.M{"$concatArrays": bson.A{"$previousRevisions", bson.A{"$currentRevision"}}}}},
	//  bson.M{"$set": bson.M{"currentRevision": in}},
	//  bson.M{"$unset": "currentRevision.metadata.uid"},
	//  bson.M{"$set": bson.M{"resourceVersion": bson.M{"$add": bson.A{1, bson.M{"$size": "$previousRevisions"}}}}},
	//  bson.M{"$set": bson.M{"updatedAt": "$$NOW"}},
	//  bson.M{"$set": bson.M{"apiVersion": gvk.GroupVersion().String()}},
	//  bson.M{"$set": bson.M{"apiGroup": gvk.GroupVersion().Group}},
	//  bson.M{"$set": bson.M{"kind": in.GetObjectKind()}},
	//})
	//if err != nil {
	//  return err
	//}
	//
	//return s.Get(ctx, key, obj)

}

func (s *Store) Watch(ctx context.Context, objPtr runtime.Object, watchFunc func(eventType string, obj runtime.Object)) error {

	stream, err := s.mongoClient.Database(s.database).Collection(s.collection).Watch(ctx, mongo.Pipeline{})
	if err != nil {
		return err
	}
	defer stream.Close(ctx)

	go func() {
		for {
			if stream.Next(ctx) {
				operationType := stream.Current.Lookup("operationType")
				operationTypeStr := ""
				if err := operationType.Unmarshal(&operationTypeStr); err != nil {
					logrus.Errorf("unable to unmarshal operationType: %v", err)
					continue
				}

				document := stream.Current.Lookup("fullDocument")
				out := reflect.New(reflect.TypeOf(objPtr).Elem()).Interface().(runtime.Object)
				if err := convertDocument(document.Value, out); err != nil {
					logrus.Errorf("error converting document: %v", err)
					continue
				}

				watchFunc(operationTypeStr, out)
			}
		}
	}()

	<-ctx.Done()

	return nil

}

var objectSliceType = reflect.TypeOf([]runtime.Object{})

var (
	errExpectFieldItems = errors.New("no Items field in this object")
	errExpectSliceItems = errors.New("Items field must be a slice of objects")
)

// SetList sets the given list object's Items member have the elements given in
// objects.
// Returns an error if list is not a List type (does not have an Items member),
// or if any of the objects are not of the right type.
func SetList(list runtime.Object, objects []runtime.Object) error {
	itemsPtr, err := GetItemsPtr(list)
	if err != nil {
		return err
	}
	items, err := conversion.EnforcePtr(itemsPtr)
	if err != nil {
		return err
	}
	if items.Type() == objectSliceType {
		items.Set(reflect.ValueOf(objects))
		return nil
	}
	slice := reflect.MakeSlice(items.Type(), len(objects), len(objects))
	for i := range objects {
		dest := slice.Index(i)
		//if dest.Type() == reflect.TypeOf(runtime.RawExtension{}) {
		//  dest = dest.FieldByName("Object")
		//}

		// check to see if you're directly assignable
		if reflect.TypeOf(objects[i]).AssignableTo(dest.Type()) {
			dest.Set(reflect.ValueOf(objects[i]))
			continue
		}

		src, err := conversion.EnforcePtr(objects[i])
		if err != nil {
			return err
		}
		if src.Type().AssignableTo(dest.Type()) {
			dest.Set(src)
		} else if src.Type().ConvertibleTo(dest.Type()) {
			dest.Set(src.Convert(dest.Type()))
		} else {
			return fmt.Errorf("item[%d]: can't assign or convert %v into %v", i, src.Type(), dest.Type())
		}
	}
	items.Set(slice)
	return nil
}

func GetItemsPtr(list runtime.Object) (interface{}, error) {
	obj, err := getItemsPtr(list)
	if err != nil {
		return nil, fmt.Errorf("%T is not a list: %v", list, err)
	}
	return obj, nil
}

// getItemsPtr returns a pointer to the list object's Items member or an error.
func getItemsPtr(list runtime.Object) (interface{}, error) {
	v, err := conversion.EnforcePtr(list)
	if err != nil {
		return nil, err
	}

	items := v.FieldByName("Items")
	if !items.IsValid() {
		return nil, errExpectFieldItems
	}
	switch items.Kind() {
	case reflect.Interface, reflect.Ptr:
		target := reflect.TypeOf(items.Interface()).Elem()
		if target.Kind() != reflect.Slice {
			return nil, errExpectSliceItems
		}
		return items.Interface(), nil
	case reflect.Slice:
		return items.Addr().Interface(), nil
	default:
		return nil, errExpectSliceItems
	}
}

func (s *Store) getState(getResp *mongo.SingleResult, key string, v reflect.Value, ignoreNotFound bool) (*objState, error) {
	state := &objState{
		meta: &storage.ResponseMeta{},
	}

	if u, ok := v.Addr().Interface().(runtime.Unstructured); ok {
		state.obj = u.NewEmptyInstance()
	} else {
		state.obj = reflect.New(v.Type()).Interface().(runtime.Object)
	}

	if getResp.Err() != nil {
		return nil, getResp.Err()
	}

	//if len(getResp.Kvs) == 0 {
	//  if !ignoreNotFound {
	//    return nil, storage.NewKeyNotFoundError(key, 0)
	//  }
	//  if err := runtime.SetZeroValue(state.obj); err != nil {
	//    return nil, err
	//  }
	//} else {

	temp := map[string]interface{}{}
	if err := getResp.Decode(&temp); err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(temp)
	if err != nil {
		return nil, err
	}

	state.rev = temp["__revision"].(int64)
	state.meta.ResourceVersion = uint64(state.rev)
	state.data = bytes
	if err := decode(s.codec, s.versioner, bytes, state.obj, state.rev); err != nil {
		return nil, err
	}
	return state, nil
}

// decode decodes value of bytes into object. It will also set the object resource version to rev.
// On success, objPtr would be set to the object.
func decode(codec runtime.Codec, versioner storage.Versioner, value []byte, objPtr runtime.Object, rev int64) error {
	if _, err := conversion.EnforcePtr(objPtr); err != nil {
		return fmt.Errorf("unable to convert output object to pointer: %v", err)
	}
	_, _, err := codec.Decode(value, nil, objPtr)
	if err != nil {
		return err
	}
	// being unable to set the version does not prevent the object from being extracted
	if err := versioner.UpdateObject(objPtr, uint64(rev)); err != nil {
		logrus.Errorf("failed to update object version: %v", err)
	}
	return nil
}

func (s *Store) updateState(st *objState, userUpdate storage.UpdateFunc) (runtime.Object, uint64, error) {
	ret, ttlPtr, err := userUpdate(st.obj, *st.meta)
	if err != nil {
		return nil, 0, err
	}

	if err := s.versioner.PrepareObjectForStorage(ret); err != nil {
		return nil, 0, fmt.Errorf("PrepareObjectForStorage failed: %v", err)
	}

	var ttl uint64
	if ttlPtr != nil {
		ttl = *ttlPtr
	}
	return ret, ttl, nil
}
