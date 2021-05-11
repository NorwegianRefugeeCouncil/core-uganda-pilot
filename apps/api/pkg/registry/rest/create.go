package rest

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg/api/validation"
	"github.com/nrc-no/core/apps/api/pkg/apis/meta"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
	"github.com/nrc-no/core/apps/api/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/api/errors"
)

type RESTCreateStrategy interface {
	runtime.ObjectTyper
	PrepareForCreate(ctx context.Context, obj runtime.Object)
	Validate(ctx context.Context, obj runtime.Object) exceptions.ErrorList
	Canonicalize(obj runtime.Object)
}

func BeforeCreate(strategy RESTCreateStrategy, ctx context.Context, obj runtime.Object) error {

	objectMeta, kind, kerr := objectMetaAndKind(strategy, obj)
	if kerr != nil {
		return kerr
	}

	objectMeta.SetDeletionTimestamp(nil)
	strategy.PrepareForCreate(ctx, obj)
	FillObjectMetaSystemFields(objectMeta)

	if errs := strategy.Validate(ctx, obj); len(errs) > 0 {
		return exceptions.NewInvalid(kind.GroupKind(), "", errs)
	}

	if errs := validation.ValidateObjectMetaAccessor(objectMeta, field.NewPath("metadata")); len(errs) > 0 {
		return exceptions.NewInvalid(kind.GroupKind(), "", errs)
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

func BeforeUpdate(strategy RESTUpdateStrategy, ctx context.Context, obj, old runtime.Object) error {

	objectMeta, kind, kerr := objectMetaAndKind(strategy, obj)
	if kerr != nil {
		return kerr
	}

	oldMeta, err := meta.Accessor(old)
	if err != nil {
		return err
	}

	strategy.PrepareForUpdate(ctx, obj, old)

	oldUid := oldMeta.GetUID()
	objectMeta.SetUID(oldUid)

	// Ignore changes to creationTimestamp
	if oldCreationTime := oldMeta.GetCreationTimestamp(); !oldCreationTime.IsZero() {
		objectMeta.SetCreationTimestamp(oldMeta.GetCreationTimestamp())
	}

	// Ignore deletionTimestamps changes
	if !oldMeta.GetDeletionTimestamp().IsZero() {
		objectMeta.SetDeletionTimestamp(oldMeta.GetDeletionTimestamp())
	}

	errs := exceptions.ErrorList{}

	errs = append(errs, strategy.ValidateUpdate(ctx, obj, old)...)
	if len(errs) > 0 {
		return exceptions.NewInvalid(kind.GroupKind(), objectMeta.GetUID(), errs)
	}

	strategy.Canonicalize(obj)
	return nil

}
