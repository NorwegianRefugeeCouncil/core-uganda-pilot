package formdefinitions

import (
	"github.com/nrc-no/core/api/pkg/apis/core"
	genericregistry "github.com/nrc-no/core/api/pkg/server2/registry/generic"
	"k8s.io/apimachinery/pkg/runtime"
)

type REST struct {
	*genericregistry.Store
}

func NewREST(scheme *runtime.Scheme, restOptionsGetter genericregistry.RESTOptionsGetter) (*REST, error) {

	strategy := NewStrategy(scheme)

	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &core.FormDefinition{} },
		NewListFunc:              func() runtime.Object { return &core.FormDefinitionList{} },
		DefaultQualifiedResource: core.Resource("formdefinitions"),
		DeleteStrategy:           strategy,
		UpdateStrategy:           strategy,
		CreateStrategy:           strategy,
	}
	options := &genericregistry.StoreOptions{RESTOptions: restOptionsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}
