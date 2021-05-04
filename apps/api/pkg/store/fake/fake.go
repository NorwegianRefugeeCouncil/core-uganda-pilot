package fake

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/apps/api/apis/exceptions"
	"github.com/nrc-no/core/apps/api/pkg/store"
	"k8s.io/apimachinery/pkg/runtime"
)

type FakeStore struct {
	objMap map[string]runtime.Object
}

var _ store.Interface = &FakeStore{}

func NewFakeStore() *FakeStore {
	return &FakeStore{
		objMap: map[string]runtime.Object{},
	}
}

func (f *FakeStore) Create(ctx context.Context, key string, obj, out runtime.Object) error {
	panic("implement me")
}

func (f *FakeStore) Get(ctx context.Context, key string, getOptions store.GetOptions, objPtr runtime.Object) error {
	obj, ok := f.objMap[key]
	if !ok {
		if getOptions.IgnoreNotFound {
			return nil
		} else {
			return exceptions.ErrNotFound.WithError(fmt.Errorf("resource with key '%s' not found", key))
		}
	}
	objPtr = obj
	return nil
}
