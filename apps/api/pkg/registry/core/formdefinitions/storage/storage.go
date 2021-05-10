package storage

import (
	"github.com/nrc-no/core/apps/api/pkg/apis/core"
	"github.com/nrc-no/core/apps/api/pkg/registry/core/formdefinitions"
	"github.com/nrc-no/core/apps/api/pkg/registry/generic"
	genericregistry "github.com/nrc-no/core/apps/api/pkg/registry/generic/registry"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
)

type REST struct {
	*genericregistry.Store
}

func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {

	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &core.FormDefinition{} },
		NewListFunc:              func() runtime.Object { return &core.FormDefinitionList{} },
		DefaultQualifiedResource: core.Resource("formdefinitions"),
		CreateStrategy:           formdefinitions.Strategy,
		UpdateStrategy:           formdefinitions.Strategy,
		DeleteStrategy:           formdefinitions.Strategy,
	}

	options := &generic.StoreOptions{
		RESTOptions: optsGetter,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return &REST{store}, nil

}
