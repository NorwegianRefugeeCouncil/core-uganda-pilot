package customresource

import (
	"context"
	structuralschema "github.com/nrc-no/core/api/pkg/openapi"
	"github.com/nrc-no/core/api/pkg/openapi/listtype"
	"github.com/nrc-no/core/api/pkg/openapi/objectmeta"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kube-openapi/pkg/validation/validate"
)

type customResourceStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
	validator         customResourceValidator
	structuralSchemas map[string]*structuralschema.Structural
	kind              schema.GroupVersionKind
	namespacedScoped  bool
}

func NewStrategy(
	typer runtime.ObjectTyper,
	namespacedScope bool,
	kind schema.GroupVersionKind,
	schemaValidator *validate.SchemaValidator,
	structuralSchemas map[string]*structuralschema.Structural,
) customResourceStrategy {
	return customResourceStrategy{
		ObjectTyper:      typer,
		NameGenerator:    names.SimpleNameGenerator,
		namespacedScoped: namespacedScope,
		validator: customResourceValidator{
			namespaceScoped: namespacedScope,
			kind:            kind,
			schemaValidator: schemaValidator,
		},
		structuralSchemas: structuralSchemas,
		kind:              kind,
	}
}

func (c customResourceStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (c customResourceStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	return
}

func (c customResourceStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	var errs field.ErrorList
	errs = append(errs, c.validator.ValidateUpdate(ctx, obj, old)...)

	uNew, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return errs
	}
	uOld, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return errs
	}

	v := obj.GetObjectKind().GroupVersionKind().Version
	errs = append(errs, objectmeta.Validate(nil, uNew.Object, c.structuralSchemas[v], false)...)

	if oldErrs := listtype.ValidateListSetsAndMaps(nil, c.structuralSchemas[v], uOld.Object); len(oldErrs) == 0 {
		errs = append(errs, listtype.ValidateListSetsAndMaps(nil, c.structuralSchemas[v], uNew.Object)...)
	}

	return errs
}

func (c customResourceStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (c customResourceStrategy) NamespaceScoped() bool {
	return c.namespacedScoped
}

func (c customResourceStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	accessor, _ := meta.Accessor(obj)
	accessor.SetGeneration(1)
}

func (c customResourceStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	var errs field.ErrorList
	errs = append(errs, c.validator.Validate(ctx, obj)...)

	if u, ok := obj.(*unstructured.Unstructured); ok {
		v := obj.GetObjectKind().GroupVersionKind().Version
		errs = append(errs, objectmeta.Validate(nil, u.Object, c.structuralSchemas[v], false)...)
		errs = append(errs, listtype.ValidateListSetsAndMaps(nil, c.structuralSchemas[v], u.Object)...)
	}

	return errs

}

func (c customResourceStrategy) Canonicalize(obj runtime.Object) {
}
