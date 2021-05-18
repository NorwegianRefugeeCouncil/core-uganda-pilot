package formdefinition

import (
	"github.com/nrc-no/core/api/pkg/apis/core"
	coreregistry "github.com/nrc-no/core/api/pkg/registry/core"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
)

func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*coreregistry.REST, error) {
	strategy := NewStrategy(scheme)

	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &core.FormDefinition{} },
		NewListFunc:              func() runtime.Object { return &core.FormDefinitionList{} },
		PredicateFunc:            MatchFormDefinition,
		DefaultQualifiedResource: core.Resource("formdefinitions"),
		CreateStrategy:           strategy,
		UpdateStrategy:           strategy,
		DeleteStrategy:           strategy,
		TableConvertor:           rest.NewDefaultTableConvertor(core.Resource("formdefinitions")),
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &coreregistry.REST{store}, nil
}
