package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Tenant represents an organization that operates core
type Tenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

// User represents an OIDC User
type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Attributes        map[string][]string `json:"attributes,omitempty" protobuf:"bytes,1,opt,name=attributes"`
}

// Claim represent an OIDC Claim
type Claim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              ClaimSpec `json:"spec"`
}

type ClaimSpec struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

// Scope represents an OIDC Scope
type Scope struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              ScopeSpec `json:"spec"`
}

type ScopeSpec struct {
	Mappers []ScopeMapper
}

type ScopeMapper struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              ScopeMapperSpec `json:"spec"`
}

type ScopeMapperSpec struct {
	Scope                  string `json:"scope"`
	Protocol               string `json:"protocol"`
	ProtocolMapper         string `json:"protocolMapper"`
	IncludeInAccessToken   bool   `json:"includeInAccessToken"`
	IncludeInIDToken       bool   `json:"includeInIdToken"`
	IncludeInUserInfoToken bool   `json:"includeInUserInfoToken"`
	ClaimName              string `json:"claimName"`
	JsonType               string `json:"jsonType"`
	UserAttribute          string `json:"userAttribute"`
}
