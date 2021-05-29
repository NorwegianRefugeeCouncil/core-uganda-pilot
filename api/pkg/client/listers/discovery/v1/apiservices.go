package v1

import (
	discoveryv1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// APIServiceLister helps list FormDefinitions.
// All objects returned here must be treated as read-only.
type APIServiceLister interface {
	// List lists all FormDefinitions in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*discoveryv1.APIService, err error)
	// Get retrieves the FormDefinition from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*discoveryv1.APIService, error)
}

// apiServiceLister implements the APIServiceLister interface.
type apiServiceLister struct {
	indexer cache.Indexer
}

// NewAPIServiceLister returns a new APIServiceLister.
func NewAPIServiceLister(indexer cache.Indexer) APIServiceLister {
	return &apiServiceLister{indexer: indexer}
}

// List lists all FormDefinitions in the indexer.
func (s *apiServiceLister) List(selector labels.Selector) (ret []*discoveryv1.APIService, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*discoveryv1.APIService))
	})
	return ret, err
}

// Get retrieves the FormDefinition from the index for a given name.
func (s *apiServiceLister) Get(name string) (*discoveryv1.APIService, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(discoveryv1.Resource("apiservices"), name)
	}
	return obj.(*discoveryv1.APIService), nil
}
