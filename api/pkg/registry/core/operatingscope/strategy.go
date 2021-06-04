package operatingscope

import (
	"context"
	rest2 "github.com/nrc-no/core/api/pkg/registry/rest"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
)

func NewStrategy(typer runtime.ObjectTyper) operatingScopeStrategy {
	return operatingScopeStrategy{typer, names.SimpleNameGenerator}
}

type operatingScopeStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (o operatingScopeStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (o operatingScopeStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {

}

func (o operatingScopeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (o operatingScopeStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (o operatingScopeStrategy) NamespaceScoped() bool {
	return false
}

func (o operatingScopeStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {

}

func (o operatingScopeStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (o operatingScopeStrategy) Canonicalize(obj runtime.Object) {

}

var _ rest2.RESTCreateStrategy = &operatingScopeStrategy{}
var _ rest2.RESTUpdateStrategy = &operatingScopeStrategy{}
var _ rest2.RESTDeleteStrategy = &operatingScopeStrategy{}
