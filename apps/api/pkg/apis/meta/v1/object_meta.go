package v1

import (
	"time"
)

type ObjectMeta struct {
	UID               string            `json:"uid,omitempty" bson:"uid,omitempty"`
	ResourceVersion   int               `json:"resourceVersion,omitempty" bson:"resourceVersion,omitempty"`
	CreationTimestamp time.Time         `json:"creationTimestamp,omitempty" bson:"creationTimestamp,omitempty"`
	DeletionTimestamp *time.Time        `json:"deletionTimestamp,omitempty" bson:"deletionTimestamp,omitempty"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
}

var _ Object = &ObjectMeta{}

func (o *ObjectMeta) GetResourceVersion() int {
	return o.ResourceVersion
}

func (o *ObjectMeta) SetResourceVersion(version int) {
	o.ResourceVersion = version
}

func (o *ObjectMeta) GetCreationTimestamp() time.Time {
	return o.CreationTimestamp
}

func (o *ObjectMeta) SetCreationTimestamp(timestamp time.Time) {
	o.CreationTimestamp = timestamp
}

func (o *ObjectMeta) GetDeletionTimestamp() *time.Time {
	return o.DeletionTimestamp
}

func (o *ObjectMeta) SetDeletionTimestamp(timestamp *time.Time) {
	o.DeletionTimestamp = timestamp
}

func (o *ObjectMeta) GetLabels() map[string]string {
	return o.Labels
}

func (o *ObjectMeta) SetLabels(labels map[string]string) {
	o.Labels = labels
}

func (o *ObjectMeta) GetAnnotations() map[string]string {
	return o.Annotations
}

func (o *ObjectMeta) SetAnnotations(annotations map[string]string) {
	o.Annotations = annotations
}

func (o *ObjectMeta) GetUID() string {
	return o.UID
}

func (o *ObjectMeta) SetUID(uid string) {
	o.UID = uid
}

func (o *ObjectMeta) GetObjectMeta() Object {
	return o
}
