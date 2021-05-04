package fake

import (
  "context"
  "encoding/json"
  "fmt"
  "github.com/nrc-no/core/apps/api/apis/exceptions"
  "github.com/nrc-no/core/apps/api/pkg/store"
  "k8s.io/apimachinery/pkg/conversion"
  "k8s.io/apimachinery/pkg/runtime"
  "reflect"
)

type FakeStore struct {
  objMap map[string][][]byte
}

var _ store.Interface = &FakeStore{}

func NewFakeStore() *FakeStore {
  return &FakeStore{
    objMap: map[string][][]byte{},
  }
}

func (f *FakeStore) Create(ctx context.Context, key string, obj, out runtime.Object) error {
  bytes, err := json.Marshal(obj)
  if err != nil {
    return err
  }
  f.objMap[key] = append(f.objMap[key], bytes)
  if err := json.Unmarshal(bytes, &out); err != nil {
    return err
  }
  return nil
}

func (f *FakeStore) Get(ctx context.Context, key string, getOptions store.GetOptions, objPtr runtime.Object) error {
  obj, ok := f.objMap[key]
  if !ok || len(obj) == 0 {
    if getOptions.IgnoreNotFound {
      return nil
    } else {
      return exceptions.ErrNotFound.WithError(fmt.Errorf("resource with key '%s' not found", key))
    }
  }
  if err := json.Unmarshal(obj[len(obj)-1], &objPtr); err != nil {
    return err
  }
  return nil
}

func (f *FakeStore) Update(ctx context.Context, key string, objType runtime.Object, updateFunc store.UpdateFunc, updateOptions store.UpdateOptions) error {

  v, err := conversion.EnforcePtr(objType)
  if err != nil {
    return fmt.Errorf("unable to convert output object to pointer: %v", err)
  }

  data, ok := f.objMap[key]
  if !ok || len(data) == 0 {
    if updateOptions.IgnoreNotFound {
      return nil
    } else {
      return exceptions.ErrNotFound.WithError(fmt.Errorf("resource with id '%s' could not be found", key))
    }
  }
  lastData := data[len(data)-1]
  var newObj runtime.Object
  if u, ok := v.Addr().Interface().(runtime.Unstructured); ok {
    newObj = u.NewEmptyInstance()
  } else {
    newObj = reflect.New(v.Type()).Interface().(runtime.Object)
  }

  if err := json.Unmarshal(lastData, &newObj); err != nil {
    return err
  }

  updated, err := updateFunc(newObj)
  if err != nil {
    return err
  }

  newBytes, err := json.Marshal(updated)
  if err != nil {
    return err
  }

  f.objMap[key] = append(f.objMap[key], newBytes)
  return nil
}
