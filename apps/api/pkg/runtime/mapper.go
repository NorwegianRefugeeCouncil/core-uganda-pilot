package runtime

import (
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"sync"
)

type equivalentResourceRegistry struct {
	// keyFunc computes a key for the specified resource (this allows honoring colocated resources across API groups).
	// if null, or if "" is returned, resource.String() is used as the key
	keyFunc func(resource schema.GroupResource) string
	// resources maps key -> subresource -> equivalent resources (subresource is not included in the returned resources).
	// main resources are stored with subresource="".
	resources map[string]map[string][]schema.GroupVersionResource
	// kinds maps resource -> subresource -> kind
	kinds map[schema.GroupVersionResource]map[string]schema.GroupVersionKind
	// keys caches the computed key for each GroupResource
	keys map[schema.GroupResource]string

	mutex sync.RWMutex
}

// NewEquivalentResourceRegistry creates a resource registry that considers all versions of a GroupResource to be equivalent.
func NewEquivalentResourceRegistry() EquivalentResourceRegistry {
	return &equivalentResourceRegistry{}
}

// NewEquivalentResourceRegistryWithIdentity creates a resource mapper with a custom identity function.
// If "" is returned by the function, GroupResource#String is used as the identity.
// GroupResources with the same identity string are considered equivalent.
func NewEquivalentResourceRegistryWithIdentity(keyFunc func(schema.GroupResource) string) EquivalentResourceRegistry {
	return &equivalentResourceRegistry{keyFunc: keyFunc}
}

func (r *equivalentResourceRegistry) EquivalentResourcesFor(resource schema.GroupVersionResource, subresource string) []schema.GroupVersionResource {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.resources[r.keys[resource.GroupResource()]][subresource]
}
func (r *equivalentResourceRegistry) KindFor(resource schema.GroupVersionResource, subresource string) schema.GroupVersionKind {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.kinds[resource][subresource]
}
func (r *equivalentResourceRegistry) RegisterKindFor(resource schema.GroupVersionResource, subresource string, kind schema.GroupVersionKind) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.kinds == nil {
		r.kinds = map[schema.GroupVersionResource]map[string]schema.GroupVersionKind{}
	}
	if r.kinds[resource] == nil {
		r.kinds[resource] = map[string]schema.GroupVersionKind{}
	}
	r.kinds[resource][subresource] = kind

	// get the shared key of the parent resource
	key := ""
	gr := resource.GroupResource()
	if r.keyFunc != nil {
		key = r.keyFunc(gr)
	}
	if key == "" {
		key = gr.String()
	}

	if r.keys == nil {
		r.keys = map[schema.GroupResource]string{}
	}
	r.keys[gr] = key

	if r.resources == nil {
		r.resources = map[string]map[string][]schema.GroupVersionResource{}
	}
	if r.resources[key] == nil {
		r.resources[key] = map[string][]schema.GroupVersionResource{}
	}
	r.resources[key][subresource] = append(r.resources[key][subresource], resource)
}
