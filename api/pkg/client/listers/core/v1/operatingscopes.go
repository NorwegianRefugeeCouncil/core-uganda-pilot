package v1

import (
	v1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// OperatingScopeLister helps list OperatingScopes.
// All objects returned here must be treated as read-only.
type OperatingScopeLister interface {
	// List lists all OperatingScopes in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.OrganizationScope, err error)
	// Get retrieves the OrganizationScope from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.OrganizationScope, error)
}

// operatingScopeLister implements the OperatingScopeLister interface.
type operatingScopeLister struct {
	indexer cache.Indexer
}

// NewOperatingScopeLister returns a new OperatingScopeLister.
func NewOperatingScopeLister(indexer cache.Indexer) OperatingScopeLister {
	return &operatingScopeLister{indexer: indexer}
}

// List lists all OperatingScopes in the indexer.
func (s *operatingScopeLister) List(selector labels.Selector) (ret []*v1.OrganizationScope, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.OrganizationScope))
	})
	return ret, err
}

// Get retrieves the OrganizationScope from the index for a given name.
func (s *operatingScopeLister) Get(name string) (*v1.OrganizationScope, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("operatingscope"), name)
	}
	return obj.(*v1.OrganizationScope), nil
}
