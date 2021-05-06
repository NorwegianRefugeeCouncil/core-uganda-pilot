package v1

import "time"

type Object interface {
	GetName() string
	SetName(name string)
	GetUID() string
	SetUID(uid string)
	GetResourceVersion() string
	SetResourceVersion(version string)
	GetCreationTimestamp() time.Time
	SetCreationTimestamp(timestamp time.Time)
	GetDeletionTimestamp() *time.Time
	SetDeletionTimestamp(timestamp *time.Time)
	GetLabels() map[string]string
	SetLabels(labels map[string]string)
	GetAnnotations() map[string]string
	SetAnnotations(annotations map[string]string)
}

type ListInterface interface {
	GetResourceVersion() string
	SetResourceVersion(version string)
}
