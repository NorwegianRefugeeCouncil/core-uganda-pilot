package v1

import (
	"fmt"
	"github.com/nrc-no/core/apps/api/pkg/conversion"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/watch"
)

const (
	// NamespaceDefault means the object is in the default namespace which is applied when not specified by clients
	NamespaceDefault = "default"
	// NamespaceAll is the default argument to specify on a context when you want to list or filter resources across all namespaces
	NamespaceAll = ""
	// NamespaceNone is the argument for a context when there is no namespace.
	NamespaceNone = ""
	// NamespaceSystem is the system namespace where we place system components.
	NamespaceSystem = "kube-system"
	// NamespacePublic is the namespace where we place public info (ConfigMaps)
	NamespacePublic = "kube-public"
)

type ObjectMeta struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	UID               string            `json:"uid,omitempty" bson:"uid,omitempty"`
	ResourceVersion   string            `json:"resourceVersion,omitempty" bson:"resourceVersion,omitempty"`
	CreationTimestamp Time              `json:"creationTimestamp,omitempty" bson:"creationTimestamp"`
	DeletionTimestamp *Time             `json:"deletionTimestamp,omitempty" bson:"deletionTimestamp"`
	Labels            map[string]string `json:"labels" bson:"labels"`
	Annotations       map[string]string `json:"annotations" bson:"annotations"`
}

var _ Object = &ObjectMeta{}

func (o *ObjectMeta) GetName() string {
	return o.Name
}

func (o *ObjectMeta) SetName(name string) {
	o.Name = name
}

func (o *ObjectMeta) GetNamespace() string {
	return o.Namespace
}

func (o *ObjectMeta) SetNamespace(namespace string) {
	o.Namespace = namespace
}

func (o *ObjectMeta) GetResourceVersion() string {
	return o.ResourceVersion
}

func (o *ObjectMeta) SetResourceVersion(version string) {
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
	ResourceVersion string `json:"resourceVersion,omitempty" bson:"resourceVersion,omitempty"`
	Continue        string `json:"continue,omitempty" protobuf:"bytes,3,opt,name=continue"`
}

func (l *ListMeta) GetResourceVersion() string {
	return l.ResourceVersion
}

func (l *ListMeta) SetResourceVersion(version string) {
	l.ResourceVersion = version
}

func (l *ListMeta) GetContinue() string {
	return l.Continue
}

func (l *ListMeta) SetContinue(continueVal string) {
	l.Continue = continueVal
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
	StatusReasonUnknown               StatusReason = "Unknown"
	StatusReasonForbidden             StatusReason = "Forbidden"
	StatusReasonNotFound              StatusReason = "NotFound"
	StatusReasonAlreadyExists         StatusReason = "AlreadyExists"
	StatusReasonConflict              StatusReason = "Conflict"
	StatusReasonBadRequest            StatusReason = "BadRequest"
	StatusReasonMethodNotAllowed      StatusReason = "MethodNotAllowed"
	StatusReasonNotAcceptable         StatusReason = "NotAcceptable"
	StatusReasonInternalError         StatusReason = "InternalError"
	StatusReasonInvalid               StatusReason = "Invalid"
	StatusReasonUnsupportedMediaType  StatusReason = "UnsupportedMediaType"
	StatusReasonServiceUnavailable    StatusReason = "ServiceUnavailable"
	StatusReasonTimeout               StatusReason = "Timeout"
	StatusReasonTooManyRequests       StatusReason = "TooManyRequests"
	StatusReasonUnauthorized          StatusReason = "Unauthorized"
	StatusReasonGone                  StatusReason = "Gone"
	StatusReasonExpired               StatusReason = "Expired"
	StatusReasonServerTimeout         StatusReason = "ServerTimeout"
	StatusReasonRequestEntityTooLarge StatusReason = "RequestEntityTooLarge"
)

type StatusCause struct {
	Type    CauseType `json:"reason,omitempty"`
	Message string    `json:"message,omitempty"`
	Field   string    `json:"field,omitempty"`
}

type CauseType string

const (
	CauseTypeFieldValueNotFound       CauseType = "FieldValueNotFound"
	CauseTypeFieldValueRequired       CauseType = "FieldValueRequired"
	CauseTypeFieldValueDuplicate      CauseType = "FieldValueDuplicate"
	CauseTypeFieldValueInvalid        CauseType = "FieldValueInvalid"
	CauseTypeFieldValueNotSupported   CauseType = "FieldValueNotSupported"
	CauseTypeUnexpectedServerResponse CauseType = "UnexpectedServerResponse"
	CauseTypeResourceVersionTooLarge  CauseType = "ResourceVersionTooLarge"
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// List holds a list of objects, which may not be known by the server.
type List struct {
	TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of objects
	Items []runtime.RawExtension `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ListOptions is the query options to a standard REST list call.
type ListOptions struct {
	TypeMeta `json:",inline"`

	// A selector to restrict the list of returned objects by their labels.
	// Defaults to everything.
	// +optional
	LabelSelector string `json:"labelSelector,omitempty" protobuf:"bytes,1,opt,name=labelSelector"`
	// A selector to restrict the list of returned objects by their fields.
	// Defaults to everything.
	// +optional
	FieldSelector string `json:"fieldSelector,omitempty" protobuf:"bytes,2,opt,name=fieldSelector"`

	// +k8s:deprecated=includeUninitialized,protobuf=6

	// Watch for changes to the described resources and return them as a stream of
	// add, update, and remove notifications. Specify resourceVersion.
	// +optional
	Watch bool `json:"watch,omitempty" protobuf:"varint,3,opt,name=watch"`
	// allowWatchBookmarks requests watch events with type "BOOKMARK".
	// Servers that do not implement bookmarks may ignore this flag and
	// bookmarks are sent at the server's discretion. Clients should not
	// assume bookmarks are returned at any specific interval, nor may they
	// assume the server will send any BOOKMARK event during a session.
	// If this is not a watch, this field is ignored.
	// +optional
	AllowWatchBookmarks bool `json:"allowWatchBookmarks,omitempty" protobuf:"varint,9,opt,name=allowWatchBookmarks"`

	// resourceVersion sets a constraint on what resource versions a request may be served from.
	// See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for
	// details.
	//
	// Defaults to unset
	// +optional
	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,4,opt,name=resourceVersion"`

	// resourceVersionMatch determines how resourceVersion is applied to list calls.
	// It is highly recommended that resourceVersionMatch be set for list calls where
	// resourceVersion is set
	// See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for
	// details.
	//
	// Defaults to unset
	// +optional
	ResourceVersionMatch ResourceVersionMatch `json:"resourceVersionMatch,omitempty" protobuf:"bytes,10,opt,name=resourceVersionMatch,casttype=ResourceVersionMatch"`
	// Timeout for the list/watch call.
	// This limits the duration of the call, regardless of any activity or inactivity.
	// +optional
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty" protobuf:"varint,5,opt,name=timeoutSeconds"`

	// limit is a maximum number of responses to return for a list call. If more items exist, the
	// server will set the `continue` field on the list metadata to a value that can be used with the
	// same initial query to retrieve the next set of results. Setting a limit may return fewer than
	// the requested amount of items (up to zero items) in the event all requested objects are
	// filtered out and clients should only use the presence of the continue field to determine whether
	// more results are available. Servers may choose not to support the limit argument and will return
	// all of the available results. If limit is specified and the continue field is empty, clients may
	// assume that no more results are available. This field is not supported if watch is true.
	//
	// The server guarantees that the objects returned when using continue will be identical to issuing
	// a single list call without a limit - that is, no objects created, modified, or deleted after the
	// first request is issued will be included in any subsequent continued requests. This is sometimes
	// referred to as a consistent snapshot, and ensures that a client that is using limit to receive
	// smaller chunks of a very large result can ensure they see all possible objects. If objects are
	// updated during a chunked list the version of the object that was present at the time the first list
	// result was calculated is returned.
	Limit int64 `json:"limit,omitempty" protobuf:"varint,7,opt,name=limit"`
	// The continue option should be set when retrieving more results from the server. Since this value is
	// server defined, clients may only use the continue value from a previous query result with identical
	// query parameters (except for the value of continue) and the server may reject a continue value it
	// does not recognize. If the specified continue value is no longer valid whether due to expiration
	// (generally five to fifteen minutes) or a configuration change on the server, the server will
	// respond with a 410 ResourceExpired error together with a continue token. If the client needs a
	// consistent list, it must restart their list without the continue field. Otherwise, the client may
	// send another list request with the token received with the 410 error, the server will respond with
	// a list starting from the next key, but from the latest snapshot, which is inconsistent from the
	// previous list results - objects that are created, modified, or deleted after the first list request
	// will be included in the response, as long as their keys are after the "next key".
	//
	// This field is not supported when watch is true. Clients may start a watch from the last
	// resourceVersion value returned by the server and not miss any modifications.
	Continue string `json:"continue,omitempty" protobuf:"bytes,8,opt,name=continue"`
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
	ResourceVersion string `json:"resourceVersion"`
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

// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type WatchEvent struct {
	Type   string               `json:"type"`
	Object runtime.RawExtension `json:"object"`
}

func Convert_watch_Event_To_v1_WatchEvent(in *watch.Event, out *WatchEvent, s conversion.Scope) error {
	out.Type = string(in.Type)
	switch t := in.Object.(type) {
	case *runtime.Unknown:
		// TODO: handle other fields on Unknown and detect type
		out.Object.Raw = t.Raw
	case nil:
	default:
		out.Object.Object = in.Object
	}
	return nil
}

func Convert_v1_InternalEvent_To_v1_WatchEvent(in *InternalEvent, out *WatchEvent, s conversion.Scope) error {
	return Convert_watch_Event_To_v1_WatchEvent((*watch.Event)(in), out, s)
}

func Convert_v1_WatchEvent_To_watch_Event(in *WatchEvent, out *watch.Event, s conversion.Scope) error {
	out.Type = watch.EventType(in.Type)
	if in.Object.Object != nil {
		out.Object = in.Object.Object
	} else if in.Object.Raw != nil {
		// TODO: handle other fields on Unknown and detect type
		out.Object = &runtime.Unknown{
			Raw:         in.Object.Raw,
			ContentType: runtime.ContentTypeJSON,
		}
	}
	return nil
}

type InternalEvent watch.Event

func (e *InternalEvent) GetObjectKind() schema.ObjectKind { return schema.EmptyObjectKind }
func (e *WatchEvent) GetObjectKind() schema.ObjectKind    { return schema.EmptyObjectKind }
func (e *InternalEvent) DeepCopyObject() runtime.Object {
	if c := e.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// +k8s:deepcopy-gen=true
// +k8s:conversion-gen:explicit-from=map

// A label selector is a label query over a set of resources. The result of matchLabels and
// matchExpressions are ANDed. An empty label selector matches all objects. A null
// label selector matches no objects.
// +structType=atomic
type LabelSelector struct {
	// matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels
	// map is equivalent to an element of matchExpressions, whose key field is "key", the
	// operator is "In", and the values array contains only "value". The requirements are ANDed.
	// +optional
	MatchLabels map[string]string `json:"matchLabels,omitempty" protobuf:"bytes,1,rep,name=matchLabels"`
	// matchExpressions is a list of label selector requirements. The requirements are ANDed.
	// +optional
	MatchExpressions []LabelSelectorRequirement `json:"matchExpressions,omitempty" protobuf:"bytes,2,rep,name=matchExpressions"`
}

// +k8s:deepcopy-gen=true

// A label selector requirement is a selector that contains values, a key, and an operator that
// relates the key and values.
type LabelSelectorRequirement struct {
	// key is the label key that the selector applies to.
	// +patchMergeKey=key
	// +patchStrategy=merge
	Key string `json:"key" patchStrategy:"merge" patchMergeKey:"key" protobuf:"bytes,1,opt,name=key"`
	// operator represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists and DoesNotExist.
	Operator LabelSelectorOperator `json:"operator" protobuf:"bytes,2,opt,name=operator,casttype=LabelSelectorOperator"`
	// values is an array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. This array is replaced during a strategic
	// merge patch.
	// +optional
	Values []string `json:"values,omitempty" protobuf:"bytes,3,rep,name=values"`
}

// A label selector operator is the set of operators that can be used in a selector requirement.
type LabelSelectorOperator string

const (
	LabelSelectorOpIn           LabelSelectorOperator = "In"
	LabelSelectorOpNotIn        LabelSelectorOperator = "NotIn"
	LabelSelectorOpExists       LabelSelectorOperator = "Exists"
	LabelSelectorOpDoesNotExist LabelSelectorOperator = "DoesNotExist"
)
