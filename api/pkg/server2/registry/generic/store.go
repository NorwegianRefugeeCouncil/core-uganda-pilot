package generic

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/api/pkg/server2/registry/rest"
	"github.com/nrc-no/core/api/pkg/server2/store"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type Store struct {
	DefaultQualifiedResource schema.GroupResource
	NewFunc                  func() runtime.Object
	NewListFunc              func() runtime.Object
	KeyFunc                  func(ctx context.Context, name string) (string, error)
	KeyRootFunc              func(ctx context.Context) (string, error)
	ObjectNameFunc           func(obj runtime.Object) (string, error)
	Storage                  store.Interface
	CreateStrategy           rest.RESTCreateStrategy
	UpdateStrategy           rest.RESTUpdateStrategy
	DeleteStrategy           rest.RESTDeleteStrategy
	DestroyFunc              store.DestroyFunc
}

var _ rest.StandardStorage = &Store{}

func (s *Store) Get(ctx context.Context, name string) (runtime.Object, error) {
	obj := s.NewFunc()
	key, err := s.KeyFunc(ctx, name)
	if err != nil {
		return nil, err
	}
	if err := s.Storage.Get(ctx, key, store.GetOptions{}, obj); err != nil {
		return nil, store.InterpretGetError(err, s.DefaultQualifiedResource, name)
	}
	return obj, nil
}

func (s *Store) NewList() runtime.Object {
	return s.NewListFunc()
}

func (s *Store) List(ctx context.Context) (runtime.Object, error) {
	list := s.NewListFunc()
	key, err := s.KeyRootFunc(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.Storage.List(ctx, key, store.ListOptions{}, list); err != nil {
		return nil, store.InterpretListError(err, s.DefaultQualifiedResource)
	}
	return list, nil
}

func (s *Store) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, error) {
	key, err := s.KeyFunc(ctx, name)
	if err != nil {
		return nil, err
	}
	out := s.NewFunc()
	if err := s.Storage.Update(ctx, key, store.UpdateOptions{}, out, func(input runtime.Object) (output runtime.Object, err error) {
		obj, err := objInfo.UpdatedObject(ctx, input)
		if err != nil {
			return nil, store.InterpretUpdateError(err, s.DefaultQualifiedResource, name)
		}
		return obj, nil
	}); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Store) Delete(ctx context.Context, name string) (runtime.Object, bool, error) {
	key, err := s.KeyFunc(ctx, name)
	if err != nil {
		return nil, false, err
	}

	obj := s.NewFunc()
	if err := s.Storage.Get(ctx, key, store.GetOptions{}, obj); err != nil {
		return nil, false, store.InterpretGetError(err, s.DefaultQualifiedResource, name)
	}

	out := s.NewFunc()
	if err := s.Storage.Delete(ctx, key, store.DeleteOptions{}, out); err != nil {
		return nil, false, store.InterpretDeleteError(err, s.DefaultQualifiedResource, name)
	}
	return out, true, nil
}

func (s *Store) Create(ctx context.Context, name string, obj runtime.Object) (runtime.Object, error) {

	if err := rest.BeforeCreate(s.CreateStrategy, ctx, obj); err != nil {
		return nil, err
	}

	name, err := s.ObjectNameFunc(obj)
	if err != nil {
		return nil, err
	}
	key, err := s.KeyFunc(ctx, name)
	if err != nil {
		return nil, err
	}
	out := s.NewFunc()
	if err := s.Storage.Create(ctx, key, store.CreateOptions{}, obj, out); err != nil {
		return nil, store.InterpretCreateError(err, s.DefaultQualifiedResource, name)
	}

	return out, nil
}

func (s *Store) New() runtime.Object {
	return s.NewFunc()
}

func (s *Store) CompleteWithOptions(options *StoreOptions) error {
	if s.DefaultQualifiedResource.Empty() {
		return fmt.Errorf("store %#v must have a non-empty qualified resource", s)
	}
	if s.NewFunc == nil {
		return fmt.Errorf("store for %s must have NewFunc set", s.DefaultQualifiedResource.String())
	}
	if s.NewListFunc == nil {
		return fmt.Errorf("store for %s must have NewListFunc set", s.DefaultQualifiedResource.String())
	}
	if s.CreateStrategy == nil {
		return fmt.Errorf("store for %s must have CreateStrategy set", s.DefaultQualifiedResource.String())
	}
	if s.UpdateStrategy == nil {
		return fmt.Errorf("store for %s must have UpdateStrategy set", s.DefaultQualifiedResource.String())
	}
	if s.DeleteStrategy == nil {
		return fmt.Errorf("store for %s must have DeleteStrategy set", s.DefaultQualifiedResource.String())
	}
	if options.RESTOptions == nil {
		return fmt.Errorf("options for %s must have RESTOptions set", s.DefaultQualifiedResource.String())
	}
	opts, err := options.RESTOptions.GetRESTOptions(s.DefaultQualifiedResource)
	if err != nil {
		return err
	}

	s.KeyFunc = func(ctx context.Context, name string) (string, error) {
		return strings.Join([]string{
			s.DefaultQualifiedResource.Group,
			s.DefaultQualifiedResource.Resource,
			name,
		}, "/"), nil
	}

	s.KeyRootFunc = func(ctx context.Context) (string, error) {
		return strings.Join([]string{
			s.DefaultQualifiedResource.Group,
			s.DefaultQualifiedResource.Resource,
		}, "/"), nil
	}

	if s.ObjectNameFunc == nil {
		s.ObjectNameFunc = func(obj runtime.Object) (string, error) {
			accessor, err := meta.Accessor(obj)
			if err != nil {
				return "", err
			}
			return accessor.GetName(), nil
		}
	}

	if s.Storage == nil {
		storage, destroy, err := store.Create(*opts.StorageConfig)
		if err != nil {
			return err
		}
		s.DestroyFunc = destroy
		s.Storage = storage
	}

	return nil
}
