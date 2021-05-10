package v1

// +k8s:deepcopy-gen=false
type Object interface {
	GetUID() string
	SetUID(uid string)
	GetResourceVersion() int
	SetResourceVersion(version int)
	GetCreationTimestamp() Time
	SetCreationTimestamp(timestamp Time)
	GetDeletionTimestamp() *Time
	SetDeletionTimestamp(timestamp *Time)
	GetLabels() map[string]string
	SetLabels(labels map[string]string)
	GetAnnotations() map[string]string
	SetAnnotations(annotations map[string]string)
}
