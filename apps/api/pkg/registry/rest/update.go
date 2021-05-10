package rest

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/util/validation/field"
)

type RESTUpdateStrategy interface {
	runtime.ObjectTyper
	PrepareForUpdate(ctx context.Context, obj, old runtime.Object)
	ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList
	Canonicalize(obj runtime.Object)
	AllowCreateOnUpdate() bool
}
