package customresourcedefinition

import (
	"github.com/nrc-no/core/api/pkg/apis/core"
	"github.com/nrc-no/core/api/pkg/server2/registry/generic"
	"k8s.io/apimachinery/pkg/runtime"
)

type REST struct {
	*generic.Store
}

func NewRest(schema *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*REST, error) {
	strategy := NewStrategy(schema)

	store := &generic.Store{
		NewFunc:                  func() runtime.Object { return &core.CustomResourceDefinition{} },
		NewListFunc:              func() runtime.Object { return &core.CustomResourceDefinitionList{} },
		DefaultQualifiedResource: core.Resource("customresourcedefinitions"),
		CreateStrategy:           strategy,
		UpdateStrategy:           strategy,
		DeleteStrategy:           strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}
