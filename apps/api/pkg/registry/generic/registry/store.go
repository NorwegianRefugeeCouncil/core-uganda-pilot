package registry

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/apps/api/pkg/apis/meta"
	"github.com/nrc-no/core/apps/api/pkg/apis/meta/internalversion"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/request"
	"github.com/nrc-no/core/apps/api/pkg/registry/generic"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/storage"
	"github.com/nrc-no/core/apps/api/pkg/storage/backend"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
	"github.com/nrc-no/core/apps/api/pkg/util/validation/field"
	"github.com/nrc-no/core/apps/api/pkg/watch"
)

const (
	OptimisticLockErrorMsg = "the object has been modified; please apply your changes to the latest version and try again"
)

// FinishFunc is a function returned by Begin hooks to complete an operation.
type FinishFunc func(ctx context.Context, success bool)

// AfterDeleteFunc is the type used for the Store.AfterDelete hook.
type AfterDeleteFunc func(obj runtime.Object, options *metav1.DeleteOptions)

// BeginCreateFunc is the type used for the Store.BeginCreate hook.
type BeginCreateFunc func(ctx context.Context, obj runtime.Object, options *metav1.CreateOptions) (FinishFunc, error)

// AfterCreateFunc is the type used for the Store.AfterCreate hook.
type AfterCreateFunc func(obj runtime.Object, options *metav1.CreateOptions)

// BeginUpdateFunc is the type used for the Store.BeginUpdate hook.
type BeginUpdateFunc func(ctx context.Context, obj, old runtime.Object, options *metav1.UpdateOptions) (FinishFunc, error)

// AfterUpdateFunc is the type used for the Store.AfterUpdate hook.
type AfterUpdateFunc func(obj runtime.Object, options *metav1.UpdateOptions)

type Store struct {
	KeyFunc                  func(ctx context.Context, name string) (string, error)
	NameFunc                 func(ctx context.Context, obj runtime.Object) (string, error)
	NewFunc                  func() runtime.Object
	NewListFunc              func() runtime.Object
	DefaultQualifiedResource schema.GroupResource
	CreateStrategy           rest.RESTCreateStrategy
	BeginCreate              BeginCreateFunc
	AfterCreate              AfterCreateFunc
	UpdateStrategy           rest.RESTUpdateStrategy
	BeginUpdate              BeginUpdateFunc
	AfterUpdate              AfterUpdateFunc
	DeleteStrategy           rest.RESTDeleteStrategy
	AfterDelete              AfterDeleteFunc
	ReturnDeletedObject      bool
	Storage                  storage.Interface
	StorageVersioner         runtime.GroupVersioner
	DestroyFunc              func()
	Decorator                func(runtime.Object)
	StorageDestroyFunc       backend.DestroyFunc
}

func (e *Store) New() runtime.Object {
	return e.NewFunc()
}

func (e *Store) NewList() runtime.Object {
	return e.NewListFunc()
}

// finishNothing is a do-nothing FinishFunc.
func finishNothing(context.Context, bool) {}

func (e *Store) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {

	var finishCreate FinishFunc = finishNothing

	if e.BeginCreate != nil {
		fn, err := e.BeginCreate(ctx, obj, options)
		if err != nil {
			return nil, err
		}
		finishCreate = fn
		defer func() {
			finishCreate(ctx, false)
		}()
	}

	if err := rest.BeforeCreate(e.CreateStrategy, ctx, obj); err != nil {
		return nil, err
	}

	if createValidation != nil {
		if err := createValidation(ctx, obj.DeepCopyObject()); err != nil {
			return nil, err
		}
	}

	name, err := e.NameFunc(ctx, obj)
	if err != nil {
		return nil, err
	}

	key, err := e.KeyFunc(ctx, name)
	if err != nil {
		return nil, err
	}

	// qualifiedResource := e.qualifiedResourceFromContext(ctx)

	out := e.NewFunc()

	if err := e.Storage.Create(ctx, obj); err != nil {
		if errGet := e.Storage.Get(ctx, key, storage.GetOptions{}, out); errGet != nil {
			return nil, err
		}
		accessor, errGetAcc := meta.Accessor(out)
		if errGetAcc != nil {
			return nil, err
		}
		if accessor.GetDeletionTimestamp() != nil {
			msg := &err.(*exceptions.StatusError).ErrStatus.Message
			*msg = fmt.Sprintf("object is being deleted: %s", *msg)
		}
		return nil, err
	}

	fn := finishCreate
	finishCreate = finishNothing
	fn(ctx, true)

	if e.AfterCreate != nil {
		e.AfterCreate(out, options)
	}
	if e.Decorator != nil {
		e.Decorator(out)
	}
	return out, nil

}

func (e *Store) qualifiedResourceFromContext(ctx context.Context) schema.GroupResource {
	if info, ok := request.RequestInfoFrom(ctx); ok {
		return schema.GroupResource{Group: info.APIGroup, Resource: info.Resource}
	}
	return e.DefaultQualifiedResource
}

func (e *Store) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {

	key, err := e.KeyFunc(ctx, name)
	if err != nil {
		return nil, false, err
	}

	var (
		// creatingObj runtime.Object
		creating = false
	)

	qualifiedResource := e.qualifiedResourceFromContext(ctx)

	out := e.NewFunc()

	err = e.Storage.Update(ctx, key, out, func(input runtime.Object, res storage.ResponseMeta) (existing runtime.Object, ttl *uint64, err error) {

		existingResourceVersion, err := e.Storage.Versioner().ObjectResourceVersion(existing)
		if err != nil {
			return nil, nil, err
		}
		if existingResourceVersion == 0 {
			if !e.UpdateStrategy.AllowCreateOnUpdate() && !forceAllowCreate {
				return nil, nil, exceptions.NewNotFound(qualifiedResource, name)
			}
		}
		obj, err := objInfo.UpdatedObject(ctx, existing)
		if err != nil {
			return nil, nil, err
		}
		newResourceVersion, err := e.Storage.Versioner().ObjectResourceVersion(obj)
		if err != nil {
			return nil, nil, err
		}

		if existingResourceVersion == 0 {
			var finishCreate FinishFunc = finishNothing

			if e.BeginCreate != nil {
				fn, err := e.BeginCreate(ctx, obj, newCreateOptionsFromUpdateOptions(options))
				if err != nil {
					return nil, nil, err
				}
				finishCreate = fn
				defer func() {
					finishCreate(ctx, false)
				}()
			}

			creating = true
			// creatingObj = obj

			if err := rest.BeforeCreate(e.CreateStrategy, ctx, obj); err != nil {
				return nil, nil, err
			}
			if createValidation != nil {
				if err := createValidation(ctx, obj.DeepCopyObject()); err != nil {
					return nil, nil, err
				}
			}

			fn := finishCreate
			finishCreate = finishNothing
			fn(ctx, true)
			return obj, nil, nil

		}

		creating = false
		// creatingObj = nil

		if newResourceVersion == 0 {
			qualifiedKind := schema.GroupKind{Group: qualifiedResource.Group, Kind: qualifiedResource.Resource}
			fieldErrList := field.ErrorList{field.Invalid(field.NewPath("metadata").Child("resourceVersion"), newResourceVersion, "must be specified for an update")}
			return nil, nil, exceptions.NewInvalid(qualifiedKind, name, fieldErrList)
		}

		if newResourceVersion != existingResourceVersion {
			return nil, nil, exceptions.NewConflict(qualifiedResource, name, fmt.Errorf(OptimisticLockErrorMsg))
		}

		var finishUpdate FinishFunc = finishNothing

		if e.BeginUpdate != nil {
			fn, err := e.BeginUpdate(ctx, obj, existing, options)
			if err != nil {
				return nil, nil, err
			}
			finishUpdate = fn
			defer func() {
				finishUpdate(ctx, false)
			}()
		}

		if err := rest.BeforeUpdate(e.UpdateStrategy, ctx, obj, existing); err != nil {
			return nil, nil, err
		}

		if updateValidation != nil {
			if err := updateValidation(ctx, obj.DeepCopyObject(), existing.DeepCopyObject()); err != nil {
				return nil, nil, err
			}
		}

		fn := finishUpdate
		finishUpdate = finishNothing
		fn(ctx, true)

		return obj, nil, nil

	})

	if err != nil {
		return nil, false, err
	}

	if creating {
		if e.AfterCreate != nil {
			e.AfterCreate(out, newCreateOptionsFromUpdateOptions(options))
		}
	} else {
		if e.AfterUpdate != nil {
			e.AfterUpdate(out, options)
		}
	}

	if e.Decorator != nil {
		e.Decorator(out)
	}

	return out, creating, nil

}

func newCreateOptionsFromUpdateOptions(in *metav1.UpdateOptions) *metav1.CreateOptions {
	co := &metav1.CreateOptions{
		//
	}
	co.TypeMeta.SetGroupVersionKind(metav1.SchemeGroupVersion.WithKind("CreateOptions"))
	return co
}

func (e *Store) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj := e.NewFunc()
	key, err := e.KeyFunc(ctx, name)
	if err != nil {
		return nil, err
	}
	if err := e.Storage.Get(ctx, key, storage.GetOptions{ResourceVersion: options.ResourceVersion}, obj); err != nil {
		return nil, err
	}
	if e.Decorator != nil {
		e.Decorator(obj)
	}
	return obj, nil
}

func (e *Store) CompleteWithOptions(options *generic.StoreOptions) error {
	if e.DefaultQualifiedResource.Empty() {
		return fmt.Errorf("store %#v must have a non-empty qualified resource", e)
	}
	if e.NewFunc == nil {
		return fmt.Errorf("store for %s must have NewFunc set", e.DefaultQualifiedResource.String())
	}
	if e.NewListFunc == nil {
		return fmt.Errorf("store for %s must have NewListFunc set", e.DefaultQualifiedResource.String())
	}
	if e.DeleteStrategy == nil {
		return fmt.Errorf("store for %s must have DeleteStrategy set", e.DefaultQualifiedResource.String())
	}
	if options.RESTOptions == nil {
		return fmt.Errorf("options for %s must ahve RESTOptions set", e.DefaultQualifiedResource.String())
	}
	opts, err := options.RESTOptions.GetRESTOptions(e.DefaultQualifiedResource)
	if err != nil {
		return err
	}

	if e.NameFunc == nil {
		e.NameFunc = func(ctx context.Context, obj runtime.Object) (string, error) {
			accessor, err := meta.Accessor(obj)
			if err != nil {
				return "", err
			}
			return accessor.GetUID(), nil
		}
	}

	keyFunc := func(obj runtime.Object) (string, error) {
		accessor, err := meta.Accessor(obj)
		if err != nil {
			return "", err
		}
		gvk := obj.GetObjectKind().GroupVersionKind()
		return gvk.Group + "/" + gvk.Kind + "/" + accessor.GetUID(), nil
	}
	e.KeyFunc = func(ctx context.Context, name string) (string, error) {
		return name, nil
	}

	if e.Storage == nil {
		s, destroyFunc, err := opts.Decorator(
			opts.StorageConfig,
			opts.ResourcePrefix,
			keyFunc,
			e.NewFunc,
			e.NewListFunc,
		)
		if err != nil {
			return err
		}
		e.Storage = s
		e.StorageDestroyFunc = destroyFunc
		e.StorageVersioner = opts.StorageConfig.EncodeVersioner
	}

	return nil

}

func (e *Store) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {

	if options == nil {
		options = &internalversion.ListOptions{ResourceVersion: ""}
	}
	list := e.NewListFunc()
	// qualifiedResource := e.qualifiedResourceFromContext(ctx)
	storageOpts := storage.ListOptions{
		ResourceVersion: options.ResourceVersion,
	}

	err := e.Storage.List(ctx, storageOpts, list)
	return list, err

}

func (e *Store) Watch(ctx context.Context, options *internalversion.ListOptions) (watch.Interface, error) {
	storageOpts := storage.ListOptions{ResourceVersion: options.ResourceVersion}
	w, err := e.Storage.Watch(ctx, "", storageOpts)
	if err != nil {
		return nil, err
	}
	return w, nil
}
