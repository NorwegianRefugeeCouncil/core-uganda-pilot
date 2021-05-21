package customresourcedefinition

import (
	"github.com/nrc-no/core/api/pkg/apis/core"
	generic2 "github.com/nrc-no/core/api/pkg/registry/generic"
	"k8s.io/apimachinery/pkg/runtime"
)

type REST struct {
	*generic2.Store
}

func NewREST(schema *runtime.Scheme, optsGetter generic2.RESTOptionsGetter) (*REST, error) {
	strategy := NewStrategy(schema)

	store := &generic2.Store{
		NewFunc:                  func() runtime.Object { return &core.CustomResourceDefinition{} },
		NewListFunc:              func() runtime.Object { return &core.CustomResourceDefinitionList{} },
		DefaultQualifiedResource: core.Resource("customresourcedefinitions"),
		CreateStrategy:           strategy,
		UpdateStrategy:           strategy,
		DeleteStrategy:           strategy,
	}
	options := &generic2.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}
