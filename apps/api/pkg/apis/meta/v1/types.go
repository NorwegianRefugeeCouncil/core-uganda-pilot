package v1

import (
	"fmt"
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

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ListOptions struct {
	TypeMeta `json:",inline"`
}

// resourceVersionMatch specifies how the resourceVersion parameter is applied. resourceVersionMatch
// may only be set if resourceVersion is also set.
//
// "NotOlderThan" matches data at least as new as the provided resourceVersion.
// "Exact" matches data at the exact resourceVersion provided.
//
// See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for
// details.
type ResourceVersionMatch string

const (
	// ResourceVersionMatchNotOlderThan matches data at least as new as the provided
	// resourceVersion.
	ResourceVersionMatchNotOlderThan ResourceVersionMatch = "NotOlderThan"
	// ResourceVersionMatchExact matches data at the exact resourceVersion
	// provided.
	ResourceVersionMatchExact ResourceVersionMatch = "Exact"
)

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type GetOptions struct {
	TypeMeta        `json:",inline"`
	ResourceVersion string
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DeleteOptions struct {
	TypeMeta `json:",inline"`
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CreateOptions struct {
	TypeMeta `json:",inline"`
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type UpdateOptions struct {
	TypeMeta `json:",inline"`
}

// APIResource specifies the name of a resource and whether it is namespaced.
type APIResource struct {
	// name is the plural name of the resource.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// singularName is the singular name of the resource.  This allows clients to handle plural and singular opaquely.
	// The singularName is more correct for reporting status on a single item and both singular and plural are allowed
	// from the kubectl CLI interface.
	SingularName string `json:"singularName" protobuf:"bytes,6,opt,name=singularName"`
	// namespaced indicates if a resource is namespaced or not.
	Namespaced bool `json:"namespaced" protobuf:"varint,2,opt,name=namespaced"`
	// group is the preferred group of the resource.  Empty implies the group of the containing resource list.
	// For subresources, this may have a different value, for example: Scale".
	Group string `json:"group,omitempty" protobuf:"bytes,8,opt,name=group"`
	// version is the preferred version of the resource.  Empty implies the version of the containing resource list
	// For subresources, this may have a different value, for example: v1 (while inside a v1beta1 version of the core resource's group)".
	Version string `json:"version,omitempty" protobuf:"bytes,9,opt,name=version"`
	// kind is the kind for the resource (e.g. 'Foo' is the kind for a resource 'foo')
	Kind string `json:"kind" protobuf:"bytes,3,opt,name=kind"`
	// verbs is a list of supported kube verbs (this includes get, list, watch, create,
	// update, patch, delete, deletecollection, and proxy)
	Verbs Verbs `json:"verbs" protobuf:"bytes,4,opt,name=verbs"`
	// shortNames is a list of suggested short names of the resource.
	ShortNames []string `json:"shortNames,omitempty" protobuf:"bytes,5,rep,name=shortNames"`
	// categories is a list of the grouped resources this resource belongs to (e.g. 'all')
	Categories []string `json:"categories,omitempty" protobuf:"bytes,7,rep,name=categories"`
	// The hash value of the storage version, the version this resource is
	// converted to when written to the data store. Value must be treated
	// as opaque by clients. Only equality comparison on the value is valid.
	// This is an alpha feature and may change or be removed in the future.
	// The field is populated by the apiserver only if the
	// StorageVersionHash feature gate is enabled.
	// This field will remain optional even if it graduates.
	// +optional
	StorageVersionHash string `json:"storageVersionHash,omitempty" protobuf:"bytes,10,opt,name=storageVersionHash"`
}

type Verbs []string

func (vs Verbs) String() string {
	return fmt.Sprintf("%v", []string(vs))
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIResourceList is a list of APIResource, it is used to expose the name of the
// resources supported in a specific group and version, and if the resource
// is namespaced.
type APIResourceList struct {
	TypeMeta `json:",inline"`
	// groupVersion is the group and version this APIResourceList is for.
	GroupVersion string `json:"groupVersion" protobuf:"bytes,1,opt,name=groupVersion"`
	// resources contains the name of the resources and if they are namespaced.
	APIResources []APIResource `json:"resources" protobuf:"bytes,2,rep,name=resources"`
}

// APIVersions lists the versions that are available, to allow clients to
// discover the API at /api, which is the root path of the legacy v1 API.
//
// +protobuf.options.(gogoproto.goproto_stringer)=false
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type APIVersions struct {
	TypeMeta `json:",inline"`
	// versions are the api versions that are available.
	Versions []string `json:"versions" protobuf:"bytes,1,rep,name=versions"`
	// a map of client CIDR to server address that is serving this group.
	// This is to help clients reach servers in the most network-efficient way possible.
	// Clients can use the appropriate server address as per the CIDR that they match.
	// In case of multiple matches, clients should use the longest matching CIDR.
	// The server returns only those CIDRs that it thinks that the client can match.
	// For example: the master will return an internal IP CIDR only, if the client reaches the server using an internal IP.
	// Server looks at X-Forwarded-For header or X-Real-Ip header or request.RemoteAddr (in that order) to get the client IP.
	ServerAddressByClientCIDRs []ServerAddressByClientCIDR `json:"serverAddressByClientCIDRs" protobuf:"bytes,2,rep,name=serverAddressByClientCIDRs"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIGroupList is a list of APIGroup, to allow clients to discover the API at
// /apis.
type APIGroupList struct {
	TypeMeta `json:",inline"`
	// groups is a list of APIGroup.
	Groups []APIGroup `json:"groups" protobuf:"bytes,1,rep,name=groups"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIGroup contains the name, the supported versions, and the preferred version
// of a group.
type APIGroup struct {
	TypeMeta `json:",inline"`
	// name is the name of the group.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// versions are the versions supported in this group.
	Versions []GroupVersionForDiscovery `json:"versions" protobuf:"bytes,2,rep,name=versions"`
	// preferredVersion is the version preferred by the API server, which
	// probably is the storage version.
	// +optional
	PreferredVersion GroupVersionForDiscovery `json:"preferredVersion,omitempty" protobuf:"bytes,3,opt,name=preferredVersion"`
	// a map of client CIDR to server address that is serving this group.
	// This is to help clients reach servers in the most network-efficient way possible.
	// Clients can use the appropriate server address as per the CIDR that they match.
	// In case of multiple matches, clients should use the longest matching CIDR.
	// The server returns only those CIDRs that it thinks that the client can match.
	// For example: the master will return an internal IP CIDR only, if the client reaches the server using an internal IP.
	// Server looks at X-Forwarded-For header or X-Real-Ip header or request.RemoteAddr (in that order) to get the client IP.
	// +optional
	ServerAddressByClientCIDRs []ServerAddressByClientCIDR `json:"serverAddressByClientCIDRs,omitempty" protobuf:"bytes,4,rep,name=serverAddressByClientCIDRs"`
}

// ServerAddressByClientCIDR helps the client to determine the server address that they should use, depending on the clientCIDR that they match.
type ServerAddressByClientCIDR struct {
	// The CIDR with which clients can match their IP to figure out the server address that they should use.
	ClientCIDR string `json:"clientCIDR" protobuf:"bytes,1,opt,name=clientCIDR"`
	// Address of this server, suitable for a client that matches the above CIDR.
	// This can be a hostname, hostname:port, IP or IP:port.
	ServerAddress string `json:"serverAddress" protobuf:"bytes,2,opt,name=serverAddress"`
}

// GroupVersion contains the "group/version" and "version" string of a version.
// It is made a struct to keep extensibility.
type GroupVersionForDiscovery struct {
	// groupVersion specifies the API group and version in the form "group/version"
	GroupVersion string `json:"groupVersion" protobuf:"bytes,1,opt,name=groupVersion"`
	// version specifies the version in the form of "version". This is to save
	// the clients the trouble of splitting the GroupVersion.
	Version string `json:"version" protobuf:"bytes,2,opt,name=version"`
}
