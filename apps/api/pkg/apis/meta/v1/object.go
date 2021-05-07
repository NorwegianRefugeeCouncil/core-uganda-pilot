package v1

import "time"

type Object interface {
	GetUID() string
	SetUID(uid string)
	GetResourceVersion() int
	SetResourceVersion(version int)
	GetCreationTimestamp() time.Time
	SetCreationTimestamp(timestamp time.Time)
	GetDeletionTimestamp() *time.Time
	SetDeletionTimestamp(timestamp *time.Time)
	GetLabels() map[string]string
	SetLabels(labels map[string]string)
	GetAnnotations() map[string]string
	SetAnnotations(annotations map[string]string)
}
