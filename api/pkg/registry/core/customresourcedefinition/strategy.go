package customresourcedefinition

import (
	"context"
	"github.com/nrc-no/core/api/pkg/apis/core"
	"github.com/nrc-no/core/api/pkg/apis/core/validation"
	request2 "github.com/nrc-no/core/api/pkg/endpoints/request"
	rest2 "github.com/nrc-no/core/api/pkg/registry/rest"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
)

type strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var _ rest2.RESTCreateStrategy = &strategy{}
var _ rest2.RESTUpdateStrategy = &strategy{}
var _ rest2.RESTDeleteStrategy = &strategy{}

func NewStrategy(typer runtime.ObjectTyper) strategy {
	return strategy{typer, names.SimpleNameGenerator}
}

func (s strategy) NamespaceScoped() bool {
	return false
}

func (s strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	crd := obj.(*core.CustomResourceDefinition)
	crd.Generation = 1
}

func (s strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	var groupVersion schema.GroupVersion
	if requestInfo, found := request2.RequestInfoFrom(ctx); found {
		groupVersion = schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}
	}
	return validation.ValidateCustomResourceDefinition(obj.(*core.CustomResourceDefinition), groupVersion)
}

func (s strategy) Canonicalize(obj runtime.Object) {
}

func (s strategy) AllowCreateOnUpdate() bool {
	return false
}

func (s strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newCRD := obj.(*core.CustomResourceDefinition)
	oldCRD := obj.(*core.CustomResourceDefinition)
	if !apiequality.Semantic.DeepEqual(oldCRD.Spec, newCRD.Spec) {
		newCRD.Generation = oldCRD.Generation + 1
	}
}

func (s strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	var groupVersion schema.GroupVersion
	if requestInfo, found := request2.RequestInfoFrom(ctx); found {
		groupVersion = schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}
	}
	return validation.ValidateCustomResourceDefinitionUpdate(
		obj.(*core.CustomResourceDefinition),
		old.(*core.CustomResourceDefinition),
		groupVersion)
}

func (s strategy) AllowUnconditionalUpdate() bool {
	return false
}
