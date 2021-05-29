package customresource

import (
	"github.com/nrc-no/core/api/pkg/endpoints/discovery"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
	"strings"
	"sync"
)

// CRDGroupDiscoveryHandler serves HTTP requests for
// discovery of the CustomResourceDefinition groups
// it serves the discoveryv1.APIGroup representing
// the different versions available for that group
type CRDGroupDiscoveryHandler struct {
	discoveryLock sync.RWMutex
	discovery     map[string]*discovery.APIGroupHandler
	delegate      http.Handler
}

func NewCRDGroupDiscoveryHandler(delegate http.Handler) *CRDGroupDiscoveryHandler {
	return &CRDGroupDiscoveryHandler{
		discovery:     map[string]*discovery.APIGroupHandler{},
		discoveryLock: sync.RWMutex{},
		delegate:      delegate,
	}
}

func (h *CRDGroupDiscoveryHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pathParams := splitPath(req.URL.Path)
	if len(pathParams) != 2 || pathParams[0] != "apis" {
		h.delegate.ServeHTTP(w, req)
		return
	}
	handler, ok := h.getDiscovery(pathParams[1])
	if !ok {
		h.delegate.ServeHTTP(w, req)
		return
	}
	handler.ServeHTTP(w, req)
}

// getDiscovery safely returns the APIGroupHandler for a group
func (h *CRDGroupDiscoveryHandler) getDiscovery(group string) (*discovery.APIGroupHandler, bool) {
	h.discoveryLock.RLock()
	defer h.discoveryLock.RUnlock()
	ret, ok := h.discovery[group]
	return ret, ok
}

// setDiscovery safely sets the APIGroupHandler for a group
func (h *CRDGroupDiscoveryHandler) setDiscovery(group string, handler *discovery.APIGroupHandler) {
	h.discoveryLock.Lock()
	defer h.discoveryLock.Unlock()
	h.discovery[group] = handler
}

// unsetDiscovery safely removes the APIGroupHandler for a group
func (h *CRDGroupDiscoveryHandler) unsetDiscovery(group string) {
	h.discoveryLock.Lock()
	defer h.discoveryLock.Unlock()
	delete(h.discovery, group)
}

// CRDVersionDiscoveryHandler serves HTTP requests for
// the discovery of CustomResourceDefinition resources
// It serves discoveryv1.APIResourceList
// Representing the different types of resources
// available under this Group/Version
type CRDVersionDiscoveryHandler struct {
	discoveryLock sync.RWMutex
	discovery     map[schema.GroupVersion]*discovery.APIVersionHandler
	delegate      http.Handler
}

func NewCRDVersionDiscoveryHandler(delegate http.Handler) *CRDVersionDiscoveryHandler {
	return &CRDVersionDiscoveryHandler{
		discovery:     map[schema.GroupVersion]*discovery.APIVersionHandler{},
		discoveryLock: sync.RWMutex{},
		delegate:      delegate,
	}
}

func (h *CRDVersionDiscoveryHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pathParams := splitPath(req.URL.Path)
	if len(pathParams) != 3 || pathParams[0] != "apis" {
		h.delegate.ServeHTTP(w, req)
		return
	}
	handler, ok := h.getDiscovery(schema.GroupVersion{Group: pathParams[1], Version: pathParams[2]})
	if !ok {
		h.delegate.ServeHTTP(w, req)
		return
	}
	handler.ServeHTTP(w, req)
}

// getDiscovery safely returns the APIVersionHandler for a schema.GroupVersion
func (h *CRDVersionDiscoveryHandler) getDiscovery(gv schema.GroupVersion) (*discovery.APIVersionHandler, bool) {
	h.discoveryLock.RLock()
	defer h.discoveryLock.RUnlock()
	ret, ok := h.discovery[gv]
	return ret, ok
}

// setDiscovery safely sets the APIVersionHandler for a schema.GroupVersion
func (h *CRDVersionDiscoveryHandler) setDiscovery(gv schema.GroupVersion, handler *discovery.APIVersionHandler) {
	h.discoveryLock.Lock()
	defer h.discoveryLock.Unlock()
	h.discovery[gv] = handler
}

// unsetDiscovery safely removes the APIVersionHandler for a schema.GroupVersion
func (h *CRDVersionDiscoveryHandler) unsetDiscovery(gv schema.GroupVersion) {
	h.discoveryLock.Lock()
	defer h.discoveryLock.Unlock()
	delete(h.discovery, gv)
}

// splitPath returns the segments for a URL path.
func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}
