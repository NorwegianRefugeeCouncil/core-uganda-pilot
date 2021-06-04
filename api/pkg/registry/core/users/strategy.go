package users

import (
	"context"
	rest2 "github.com/nrc-no/core/api/pkg/registry/rest"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
)

func NewStrategy(typer runtime.ObjectTyper) userStrategy {
	return userStrategy{typer, names.SimpleNameGenerator}
}

type userStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (o userStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (o userStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {

}

func (o userStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (o userStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (o userStrategy) NamespaceScoped() bool {
	return false
}

func (o userStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {

}

func (o userStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (o userStrategy) Canonicalize(obj runtime.Object) {

}

var _ rest2.RESTCreateStrategy = &userStrategy{}
var _ rest2.RESTUpdateStrategy = &userStrategy{}
var _ rest2.RESTDeleteStrategy = &userStrategy{}
