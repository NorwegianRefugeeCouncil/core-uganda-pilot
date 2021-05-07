package v1

// ListMeta describes metadata that synthetic resources must have
type ListMeta struct {
	ResourceVersion int `json:"resourceVersion,omitempty" bson:"resourceVersion,omitempty"`
}

func (l *ListMeta) GetResourceVersion() int {
	return l.ResourceVersion
}

func (l *ListMeta) SetResourceVersion(version int) {
	l.ResourceVersion = version
}
