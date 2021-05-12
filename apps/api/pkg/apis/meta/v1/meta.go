package v1

// +k8s:deepcopy-gen=false
type ObjectMetaAccessor interface {
	GetObjectMeta() Object
}

// ListInterface lets you work with list metadata from any of the versioned or
// internal API objects. Attempting to set or retrieve a field on an object that does
// not support that field will be a no-op and return a default value.
// +k8s:deepcopy-gen=false
// TODO: move this, and TypeMeta and ListMeta, to a different package
type ListInterface interface {
	GetResourceVersion() int
	SetResourceVersion(version int)
	//GetSelfLink() string
	//SetSelfLink(selfLink string)
	GetContinue() string
	SetContinue(c string)
	//GetRemainingItemCount() *int64
	//SetRemainingItemCount(c *int64)
}

// Type exposes the type and APIVersion of versioned or internal API objects.
// TODO: move this, and TypeMeta and ListMeta, to a different package
type Type interface {
	GetAPIVersion() string
	SetAPIVersion(version string)
	GetKind() string
	SetKind(kind string)
}

// ListMetaAccessor retrieves the list interface from an object
type ListMetaAccessor interface {
	GetListMeta() ListInterface
}
