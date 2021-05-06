package runtime

type Object interface {
	GetResourceVersion() int
	SetResourceVersion(version int)
	GetUID() string
	SetUID(uid string)
	GetAPIVersion() string
	SetAPIVersion(version string)
	GetAPIGroup() string
	SetAPIGroup(apiGroup string)
	GetKind() string
	SetKind(kind string)
}
