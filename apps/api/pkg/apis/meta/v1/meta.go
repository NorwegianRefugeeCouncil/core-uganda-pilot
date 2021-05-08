package v1

// +k8s:deepcopy-gen=false
type ObjectMetaAccessor interface {
	GetObjectMeta() Object
}

// ListInterface lets you work with list metadata from any of the versioned or
// internal API objects. Attempting to set or retrieve a field on an object that does
// not support that field will be a no-op and return a default value.
// TODO: move this, and TypeMeta and ListMeta, to a different package
type ListInterface interface {
	GetResourceVersion() int
	SetResourceVersion(version int)
	//GetSelfLink() string
	//SetSelfLink(selfLink string)
	//GetContinue() string
	//SetContinue(c string)
	//GetRemainingItemCount() *int64
	//SetRemainingItemCount(c *int64)
}
