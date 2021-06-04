package operatingscope

import (
	"github.com/nrc-no/core/api/pkg/apis/core"
	"github.com/nrc-no/core/api/pkg/registry/generic"
	"k8s.io/apimachinery/pkg/runtime"
)

type REST struct {
	*generic.Store
}

func NewREST(scheme *runtime.Scheme, restOptionsGetter generic.RESTOptionsGetter) (*REST, error) {

	strategy := NewStrategy(scheme)

	store := &generic.Store{
		NewFunc:                  func() runtime.Object { return &core.OperatingScope{} },
		NewListFunc:              func() runtime.Object { return &core.OperatingScopeList{} },
		DefaultQualifiedResource: core.Resource("operatingscopes"),
		DeleteStrategy:           strategy,
		UpdateStrategy:           strategy,
		CreateStrategy:           strategy,
	}
	options := &generic.StoreOptions{RESTOptions: restOptionsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}
