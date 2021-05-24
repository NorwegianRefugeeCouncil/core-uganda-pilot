package v1

import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ListResourcesOptions represent options for list operation
type ListResourcesOptions struct {
	v1.TypeMeta         `json:",inline"`
	Watch               bool   `json:"watch,omitempty" protobuf:"bytes,1,opt,name=watch"`
	AllowWatchBookmarks bool   `json:"allowWatchBookmarks,omitempty" protobuf:"bytes,1,opt,name=allowWatchBookmarks"`
	ResourceVersion     string `json:"resourceVersion,omitempty" protobuf:"bytes,1,opt,name=resourceVersion"`
	TimeoutSeconds      *int64 `json:"timeoutSeconds" protobuf:"bytes,1,opt,name=timeoutSeconds"`
	Limit               *int64 `json:"limit,omitempty" protobuf:"bytes,1,opt,name=limit"`
	Continue            string `json:"continue" protobuf:"bytes,1,opt,name=continue"`
}
