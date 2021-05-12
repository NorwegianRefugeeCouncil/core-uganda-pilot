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
	"github.com/nrc-no/core/apps/api/pkg/watch"
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
	collection  *mongo.Collection
	create      func() runtime.Object
	codec       runtime.Codec
	versioner   storage.Versioner
	watcher     *watcher
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

	parts := strings.Split(prefix, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("expected prefix in form <database>/<group>/<resource>")
	}
	collection := parts[1] + "__" + parts[2]
	collection = strings.Replace(collection, ".", "_", -1)

	mongoCollection := mongoClient.Database(parts[0]).Collection(collection)

	versioner := APIObjectVersioner{}
	return &Store{
		mongoClient: mongoClient,
		database:    parts[0],
		collection:  mongoCollection,
		create:      create,
		codec:       codec,
		watcher:     newWatcher(mongoCollection, codec, create, versioner),
		versioner:   versioner,
	}, nil
}

var _ storage.Interface = &Store{}

func (s *Store) Get(ctx context.Context, key string, getOptions storage.GetOptions, out runtime.Object) error {

	objectID, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return err
	}

	result := s.collection.FindOne(ctx, bson.M{"_id": objectID})
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

	result, err := s.collection.Find(ctx, bson.D{})
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

	if err := bson.Unmarshal(data, into); err != nil {
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

	return nil

}

func (s *Store) Create(ctx context.Context, obj, out runtime.Object) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	if err := s.versioner.UpdateObject(obj, 1); err != nil {
		return err
	}

	response, err := s.collection.InsertOne(ctx, bson.M{
		"__revision": int64(1),
		"current":    obj,
		"previous":   nil,
	})
	if err != nil {
		return err
	}

	objectID := response.InsertedID.(primitive.ObjectID)
	id := objectID.Hex()
	accessor.SetUID(id)

	if out != nil {
		bytes, err := runtime.Encode(s.codec, obj)
		if err != nil {
			return err
		}
		_, _, err = s.codec.Decode(bytes, nil, out)
		return err
	}

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
		getResp := s.collection.FindOne(ctx, bson.M{"_id": objectID})
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

		_, err = s.collection.UpdateOne(
			ctx,
			bson.M{
				"_id":        objectID,
				"__revision": origState.rev,
			},
			bson.A{
				bson.M{"$set": bson.M{
					"previous":   "$current",
					"__revision": origState.rev + 1,
				},
				},
				bson.M{"$set": bson.M{
					"current": temp,
				},
				},
			},
		)
		if err != nil {
			return err
		}

		return decode(s.codec, s.versioner, data, out, origState.rev+1)

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

func (s *Store) Delete(ctx context.Context, key string, out runtime.Object, validateDeletion storage.ValidateObjectFunc) error {

	v, err := conversion.EnforcePtr(out)
	if err != nil {
		return fmt.Errorf("unable to convert output object to pointer: %v", err)
	}

	objectId, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return fmt.Errorf("could not convert key to objectID: %v", err)
	}

	getCurrentState := func() (*objState, error) {
		getResp := s.collection.FindOne(ctx, bson.M{"_id": objectId})
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

		if err := validateDeletion(ctx, origState.obj); err != nil {
			if origStateIsCurrent {
				return err
			}
			origState, err = getCurrentState()
			if err != nil {
				return err
			}
			origStateIsCurrent = true
			continue
		}

		// First, we need to set the "previous" to the "current" value and
		// erase the "current" value. This will notify watchers that
		// a document was deleted, while they would still be able to retrieve
		// the deleted document in the "previous" key.
		_, err := s.collection.UpdateOne(ctx, bson.M{
			"_id": objectId,
		}, bson.A{
			bson.M{"$set": bson.M{"previous": "$current"}},
			bson.M{"$set": bson.M{"current": nil}},
		})
		if err != nil {
			return err
		}

		// Then, we delete the document altogether
		_, err = s.collection.DeleteOne(ctx, bson.M{"_id": objectId})
		if err != nil {
			return err
		}

		return decode(s.codec, s.versioner, origState.data, out, origState.rev)
	}

}

func (s *Store) Watch(ctx context.Context, key string, opts storage.ListOptions) (watch.Interface, error) {
	return s.watcher.Watch(ctx, key, 0, false)
}

var objectSliceType = reflect.TypeOf([]runtime.Object{})

var (
	errExpectFieldItems = errors.New("no Items field in this object")
	errExpectSliceItems = errors.New("Items field must be a slice of objects")
)

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

	accessor, err := meta.Accessor(state.obj)
	if err != nil {
		return nil, err
	}

	id := temp["_id"].(primitive.ObjectID)
	accessor.SetUID(id.Hex())

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
