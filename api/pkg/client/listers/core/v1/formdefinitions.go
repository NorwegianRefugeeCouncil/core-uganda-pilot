package v1

import (
	v1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// FormDefinitionLister helps list FormDefinitions.
// All objects returned here must be treated as read-only.
type FormDefinitionLister interface {
	// List lists all FormDefinitions in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.FormDefinition, err error)
	// Get retrieves the FormDefinition from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.FormDefinition, error)
}

// formDefinitionLister implements the FormDefinitionLister interface.
type formDefinitionLister struct {
	indexer cache.Indexer
}

// NewFormDefinitionLister returns a new FormDefinitionLister.
func NewFormDefinitionLister(indexer cache.Indexer) FormDefinitionLister {
	return &formDefinitionLister{indexer: indexer}
}

// List lists all FormDefinitions in the indexer.
func (s *formDefinitionLister) List(selector labels.Selector) (ret []*v1.FormDefinition, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.FormDefinition))
	})
	return ret, err
}

// Get retrieves the FormDefinition from the index for a given name.
func (s *formDefinitionLister) Get(name string) (*v1.FormDefinition, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("formdefinition"), name)
	}
	return obj.(*v1.FormDefinition), nil
}
