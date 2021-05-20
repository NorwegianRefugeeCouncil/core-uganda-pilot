package customresource

import (
	"github.com/nrc-no/core/api/pkg/server2/registry/generic"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type CustomResourceStore struct {
	CustomResource *REST
}

type REST struct {
	*generic.Store
}

func NewStorage(resource schema.GroupResource, kind, listKind schema.GroupVersionKind, strategy customResourceStrategy, optsGetter generic.RESTOptionsGetter) (*CustomResourceStore, error) {
	rest, err := newRest(resource, kind, listKind, strategy, optsGetter)
	if err != nil {
		return nil, err
	}
	return &CustomResourceStore{
		rest,
	}, nil
}

func newRest(resource schema.GroupResource, kind, listKind schema.GroupVersionKind, strategy customResourceStrategy, optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &generic.Store{
		NewFunc: func() runtime.Object {
			ret := &unstructured.Unstructured{}
			ret.SetGroupVersionKind(kind)
			return ret
		},
		NewListFunc: func() runtime.Object {
			ret := &unstructured.UnstructuredList{}
			ret.SetGroupVersionKind(listKind)
			return ret
		},
		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
	}

	opts := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(opts); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}
