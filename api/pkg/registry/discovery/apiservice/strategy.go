package apiservice

import (
	"context"
	"github.com/nrc-no/core/api/pkg/apis/discovery"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage/names"
)

func NewStrategy(typer runtime.ObjectTyper) apiServiceStrategy {
	return apiServiceStrategy{typer, names.SimpleNameGenerator}
}

type apiServiceStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var _ rest.RESTCreateUpdateStrategy = apiServiceStrategy{}

func (a apiServiceStrategy) NamespaceScoped() bool {
	return false
}

func (a apiServiceStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	return
}

func (a apiServiceStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	// TODO
	return field.ErrorList{}
}

func (a apiServiceStrategy) Canonicalize(obj runtime.Object) {
}

func (a apiServiceStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (a apiServiceStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newAPIService := obj.(*discovery.APIService)
	oldAPIService := obj.(*discovery.APIService)
	newAPIService.Spec = oldAPIService.Spec
	newAPIService.Labels = oldAPIService.Labels
	newAPIService.Annotations = oldAPIService.Annotations
	newAPIService.Finalizers = oldAPIService.Finalizers
	newAPIService.OwnerReferences = oldAPIService.OwnerReferences
}

func (a apiServiceStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	// TODO
	return field.ErrorList{}
}

func (a apiServiceStrategy) AllowUnconditionalUpdate() bool {
	return false
}
