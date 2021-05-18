package v1

import (
	v1 "github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/labels"
	"github.com/nrc-no/core/apps/api/pkg/tools/cache"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
)

type FormDefinitionLister interface {
	List(selector labels.Selector) (ret []*v1.FormDefinition, err error)
	FormDefinitions(namespace string) FormDefinitionNamespaceLister
}

type formDefinitionsLister struct {
	indexer cache.Indexer
}

func NewFormDefinitionLister(indexer cache.Indexer) FormDefinitionLister {
	return &formDefinitionsLister{indexer: indexer}
}

func (s *formDefinitionsLister) List(selector labels.Selector) (ret []*v1.FormDefinition, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.FormDefinition))
	})
	return ret, err
}

func (s *formDefinitionsLister) FormDefinitions(namespace string) FormDefinitionNamespaceLister {
	return formDefinitionNamespaceLister{indexer: s.indexer, namespace: namespace}
}

type FormDefinitionNamespaceLister interface {
	List(selector labels.Selector) (ret []*v1.FormDefinition, err error)
	Get(name string) (*v1.FormDefinition, error)
}

type formDefinitionNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

func (s formDefinitionNamespaceLister) List(selector labels.Selector) (ret []*v1.FormDefinition, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.FormDefinition))
	})
	return ret, err
}

func (s formDefinitionNamespaceLister) Get(name string) (*v1.FormDefinition, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, exceptions.NewNotFound(v1.Resource("formdefinition"), name)
	}
	return obj.(*v1.FormDefinition), nil
}
