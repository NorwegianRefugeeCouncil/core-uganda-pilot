package v1

import (
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
)

type ObjectMeta struct {
	UID               string            `json:"uid,omitempty" bson:"uid,omitempty"`
	ResourceVersion   int               `json:"resourceVersion,omitempty" bson:"resourceVersion,omitempty"`
	CreationTimestamp Time              `json:"creationTimestamp,omitempty" bson:"creationTimestamp"`
	DeletionTimestamp *Time             `json:"deletionTimestamp,omitempty" bson:"deletionTimestamp"`
	Labels            map[string]string `json:"labels" bson:"labels"`
	Annotations       map[string]string `json:"annotations" bson:"annotations"`
}

var _ Object = &ObjectMeta{}

func (o *ObjectMeta) GetResourceVersion() int {
	return o.ResourceVersion
}

func (o *ObjectMeta) SetResourceVersion(version int) {
	o.ResourceVersion = version
}

func (o *ObjectMeta) GetCreationTimestamp() Time {
	return o.CreationTimestamp
}

func (o *ObjectMeta) SetCreationTimestamp(timestamp Time) {
	o.CreationTimestamp = timestamp
}

func (o *ObjectMeta) GetDeletionTimestamp() *Time {
	return o.DeletionTimestamp
}

func (o *ObjectMeta) SetDeletionTimestamp(timestamp *Time) {
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

type StatusType string

const (
	StatusSuccess StatusType = "Success"
	StatusFailure StatusType = "Failure"
)

// +kubebuilder:object:root=true
type Status struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	// Status is the status of the operation. One of "Failure" or "Success"
	Status StatusType `json:"status,omitempty"`

	// Message is a human-readable description of this operation
	Message string `json:"message,omitempty"`

	// Reason is a machine-readable description of why this operation is in the
	// "Failure" status.
	Reason StatusReason `json:"reason,omitempty"`

	// Details represents extended data associated with the reason.
	Details *StatusDetails `json:"details"`

	// Suggested HTTP status code.
	Code int32 `json:"code,omitempty"`
}

type StatusDetails struct {
	UID    string        `json:"uid,omitempty"`
	Group  string        `json:"group,omitempty"`
	Kind   string        `json:"kind,omitempty"`
	Causes []StatusCause `json:"causes,omitempty"`
}

type StatusReason string

const (
	StatusReasonUnknown          StatusReason = "Unknown"
	StatusReasonForbidden        StatusReason = "Forbidden"
	StatusReasonNotFound         StatusReason = "NotFound"
	StatusReasonAlreadyExists    StatusReason = "AlreadyExists"
	StatusReasonConflict         StatusReason = "Conflict"
	StatusReasonBadRequest       StatusReason = "BadRequest"
	StatusReasonMethodNotAllowed StatusReason = "MethodNotAllowed"
	StatusReasonNotAcceptable    StatusReason = "NotAcceptable"
	StatusReasonInternalError    StatusReason = "InternalError"
	StatusReasonInvalid          StatusReason = "Invalid"
)

type StatusCause struct {
	Type    CauseType `json:"reason,omitempty"`
	Message string    `json:"message,omitempty"`
	Field   string    `json:"field,omitempty"`
}

type CauseType string

const (
	CauseTypeFieldValueNotFound     CauseType = "FieldValueNotFound"
	CauseTypeFieldValueRequired     CauseType = "FieldValueRequired"
	CauseTypeFieldValueDuplicate    CauseType = "FieldValueDuplicate"
	CauseTypeFieldValueInvalid      CauseType = "FieldValueInvalid"
	CauseTypeFieldValueNotSupported CauseType = "FieldValueNotSupported"
)

// TypeMeta represents an individual object in an API response or request.
// It represents the API schema version and kind/type of object
//
// +k8s:deepcopy-gen=false
type TypeMeta struct {

	// APIVersion defines the versioned schema of this representation
	// of an object.
	APIVersion string `json:"apiVersion,omitempty" bson:"apiVersion,omitempty"`

	// Kind is a string value representing the REST resource this
	// object represents
	Kind string `json:"kind,omitempty" bson:"kind,omitempty"`
}

func (obj *TypeMeta) GetObjectKind() schema.ObjectKind { return obj }

func (obj *TypeMeta) SetGroupVersionKind(gvk schema.GroupVersionKind) {
	obj.APIVersion, obj.Kind = gvk.ToAPIVersionAndKind()
}
func (obj *TypeMeta) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(obj.APIVersion, obj.Kind)
}
