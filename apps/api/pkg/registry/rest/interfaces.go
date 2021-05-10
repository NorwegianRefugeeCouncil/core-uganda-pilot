package rest

import (
	"context"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
)

type ValidateObjectFunc func(ctx context.Context, obj runtime.Object) error
type ValidateObjectUpdateFunc func(ctx context.Context, obj, old runtime.Object) error

type Storage interface {
	New() runtime.Object
}

type KindProvider interface {
	Kind() string
}

type GroupVersionKindProvider interface {
	GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind
}

type GroupVersionAcceptor interface {
	AcceptsGroupVersion(gv schema.GroupVersion) bool
}

type Lister interface {
	NewList() runtime.Object
	List(ctx context.Context) (runtime.Object, error)
}

type Getter interface {
	Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error)
}

type Deleter interface {
	Delete(ctx context.Context, uid string, validation ValidateObjectFunc) (runtime.Object, bool, error)
}

type Creater interface {
	New() runtime.Object
	Create(ctx context.Context, obj runtime.Object, createValidation ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error)
}

type Updater interface {
	New() runtime.Object
	Update(ctx context.Context, name string, objInfo UpdatedObjectInfo, createValidation ValidateObjectFunc, updateValidation ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error)
}

type UpdatedObjectInfo interface {
	UpdatedObject(ctx context.Context, oldObj runtime.Object) (newObj runtime.Object, err error)
}
