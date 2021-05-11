package formdefinitions

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg/api/defaultscheme"
	"github.com/nrc-no/core/apps/api/pkg/apis/core"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
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

func (s strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) exceptions.ErrorList {
	return exceptions.ErrorList{}
}

func (s strategy) Canonicalize(obj runtime.Object) {

}

func (s strategy) AllowCreateOnUpdate() bool {
	return false
}

func (s strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {

}

func (s strategy) Validate(ctx context.Context, obj runtime.Object) exceptions.ErrorList {

	var fieldErrors exceptions.ErrorList

	formDefinition, ok := obj.(*core.FormDefinition)
	if !ok {
		return fieldErrors
	}

	if formDefinition.TypeMeta.Kind == "" {
		fieldErrors = append(fieldErrors, exceptions.Required(field.NewPath("kind"), "kind is required"))
	}
	if formDefinition.TypeMeta.APIVersion == "" {
		fieldErrors = append(fieldErrors, exceptions.Required(field.NewPath("apiVersion"), "apiVersion is required"))
	}
	if formDefinition.Spec.Group == "" {
		fieldErrors = append(fieldErrors, exceptions.Required(field.NewPath("spec.group"), "group is required"))
	}
	if formDefinition.Spec.Names.Singular == "" {
		fieldErrors = append(fieldErrors, exceptions.Required(field.NewPath("spec.names.singular"), "singular name is required"))
	}
	if formDefinition.Spec.Names.Plural == "" {
		fieldErrors = append(fieldErrors, exceptions.Required(field.NewPath("spec.names.plural"), "plural name is required"))
	}
	if formDefinition.Spec.Names.Kind == "" {
		fieldErrors = append(fieldErrors, exceptions.Required(field.NewPath("spec.names.kind"), "kind name is required"))
	}
	if len(formDefinition.Spec.Versions) == 0 {
		fieldErrors = append(fieldErrors, exceptions.Required(field.NewPath("spec.versions"), "versions must not be empty"))
	}
	versionsPath := field.NewPath("spec.versions")
	for i, version := range formDefinition.Spec.Versions {
		versionPath := versionsPath.Index(i)
		if len(version.Name) == 0 {
			fieldErrors = append(fieldErrors, exceptions.Required(versionPath.Child("name"), "version name is required"))
		}
		root := version.Schema.FormSchema.Root
		rootPath := versionPath.Child("schema", "formSchema", "root")
		if root.Type != "section" {
			if root.Key == "" {
				fieldErrors = append(fieldErrors, exceptions.Required(rootPath.Child("key"), "root key is required when root is not section"))
			}
		}
	}

	return fieldErrors
}
