package rest

import (
	"context"
	v1 "github.com/nrc-no/core/api/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type Storage interface {
	New() runtime.Object
}

type KindProvider interface {
	Kind() string
}

type Lister interface {
	NewList() runtime.Object
	List(ctx context.Context) (runtime.Object, error)
}

type Getter interface {
	Get(ctx context.Context, name string) (runtime.Object, error)
}

type Deleter interface {
	Delete(ctx context.Context, name string) (runtime.Object, bool, error)
}

type Creater interface {
	New() runtime.Object
	Create(ctx context.Context, name string, obj runtime.Object) (runtime.Object, error)
}

type UpdatedObjectInfo interface {
	UpdatedObject(ctx context.Context, oldObj runtime.Object) (newObj runtime.Object, err error)
}

type Updater interface {
	New() runtime.Object
	Update(ctx context.Context, name string, objInfo UpdatedObjectInfo) (runtime.Object, error)
}

type Watcher interface {
	Watch(ctx context.Context, options v1.ListResourcesOptions) (watch.Interface, error)
}

type StandardStorage interface {
	Getter
	Lister
	Updater
	Deleter
	Creater
}

type TransformFunc func(ctx context.Context, newObj runtime.Object, oldObj runtime.Object) (transformedNewObj runtime.Object, err error)

type defaultUpdatedObjectInfo struct {
	obj          runtime.Object
	transformers []TransformFunc
}

func (d *defaultUpdatedObjectInfo) UpdatedObject(ctx context.Context, oldObj runtime.Object) (runtime.Object, error) {
	var err error
	newObj := d.obj
	if newObj != nil {
		newObj = newObj.DeepCopyObject()
	}
	for _, transformer := range d.transformers {
		newObj, err = transformer(ctx, newObj, oldObj)
		if err != nil {
			return nil, err
		}
	}
	return newObj, nil
}

func DefaultUpdatedObjectInfo(obj runtime.Object, transformers ...TransformFunc) UpdatedObjectInfo {
	return &defaultUpdatedObjectInfo{obj, transformers}
}
