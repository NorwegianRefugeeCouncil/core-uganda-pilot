package formdefinition

import (
	"context"
	"fmt"
	"github.com/nrc-no/coreapi/pkg/apis/core"
	"github.com/nrc-no/coreapi/pkg/apis/core/validation"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
)

func NewStrategy(typer runtime.ObjectTyper) formDefinitionStrategy {
	return formDefinitionStrategy{typer, names.SimpleNameGenerator}
}

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	formDefinition, ok := obj.(*core.FormDefinition)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a FormDefinition")
	}
	return formDefinition.ObjectMeta.Labels, SelectableFields(formDefinition), nil
}

func SelectableFields(obj *core.FormDefinition) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, false)
}

func MatchFormDefinition(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

type formDefinitionStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var _ rest.RESTCreateStrategy = &formDefinitionStrategy{}

func (formDefinitionStrategy) NamespaceScoped() bool {
	return false
}

func (formDefinitionStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {

}

func (formDefinitionStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {

}

func (formDefinitionStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	formDefinition := obj.(*core.FormDefinition)
	return validation.ValidateFormDefinition(formDefinition)
}

func (formDefinitionStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (formDefinitionStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (formDefinitionStrategy) Canonicalize(obj runtime.Object) {
}

func (formDefinitionStrategy) ValidateUpdate(ctx context.Context, old, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}
