package formdefinitions

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg/api/defaultscheme"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/util/validation/field"
)

type strategy struct {
	runtime.ObjectTyper
}

var Strategy = strategy{defaultscheme.Scheme}

var _ rest.RESTCreateStrategy = Strategy
var _ rest.RESTUpdateStrategy = Strategy

func (s strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {

}

func (s strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (s strategy) Canonicalize(obj runtime.Object) {

}

func (s strategy) AllowCreateOnUpdate() bool {
	return false
}

func (s strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {

}

func (s strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}
