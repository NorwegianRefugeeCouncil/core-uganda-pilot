package rest

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/api/validation/path"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
)

type RESTCreateStrategy interface {
	runtime.ObjectTyper
	names.NameGenerator
	NamespaceScoped() bool
	PrepareForCreate(ctx context.Context, obj runtime.Object)
	Validate(ctx context.Context, obj runtime.Object) field.ErrorList
	Canonicalize(obj runtime.Object)
}

type RESTDeleteStrategy interface {
	runtime.ObjectTyper
}

type RESTUpdateStrategy interface {
	runtime.ObjectTyper
	NamespaceScoped() bool
	AllowCreateOnUpdate() bool
	PrepareForUpdate(ctx context.Context, obj, old runtime.Object)
	ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList
	Canonicalize(obj runtime.Object)
	AllowUnconditionalUpdate() bool
}

func BeforeCreate(strategy RESTCreateStrategy, ctx context.Context, obj runtime.Object) error {
	objectMeta, kind, kerr := objectMetaAndKind(strategy, obj)
	if kerr != nil {
		return kerr
	}
	objectMeta.SetDeletionTimestamp(nil)
	objectMeta.SetDeletionGracePeriodSeconds(nil)
	strategy.PrepareForCreate(ctx, obj)
	FillObjectMetaSystemFields(objectMeta)
	if len(objectMeta.GetGenerateName()) > 0 && len(objectMeta.GetName()) == 0 {
		objectMeta.SetName(strategy.GenerateName(objectMeta.GetGenerateName()))
	}
	if errs := strategy.Validate(ctx, obj); len(errs) > 0 {
		return errors.NewInvalid(kind.GroupKind(), objectMeta.GetName(), errs)
	}
	if errs := validation.ValidateObjectMetaAccessor(objectMeta, strategy.NamespaceScoped(), path.ValidatePathSegmentName, field.NewPath("metadata")); len(errs) > 0 {
		return errors.NewInvalid(kind.GroupKind(), objectMeta.GetName(), errs)
	}
	strategy.Canonicalize(obj)
	return nil
}

func objectMetaAndKind(typer runtime.ObjectTyper, obj runtime.Object) (metav1.Object, schema.GroupVersionKind, error) {
	objectMeta, err := meta.Accessor(obj)
	if err != nil {
		return nil, schema.GroupVersionKind{}, errors.NewInternalError(err)
	}
	kinds, _, err := typer.ObjectKinds(obj)
	if err != nil {
		return nil, schema.GroupVersionKind{}, errors.NewInternalError(err)
	}
	return objectMeta, kinds[0], nil
}

func FillObjectMetaSystemFields(meta metav1.Object) {
	meta.SetCreationTimestamp(metav1.Now())
	meta.SetUID(uuid.NewUUID())
	meta.SetSelfLink("")
}
