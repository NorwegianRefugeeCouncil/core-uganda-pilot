package example

import (
	v1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
)

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TestModelUrlValues struct {
	v1.TypeMeta `json:",inline"`
	Abc         string `json:"abc"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TestModel struct {
	v1.TypeMeta   `json:",inline"`
	v1.ObjectMeta `json:"metadata,inline"`
	Spec          TestModelSpec `json:"spec"`
}

type TestModelSpec struct {
	SomeProperty string `json:"someProperty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TestModel2 struct {
	v1.TypeMeta   `json:",inline"`
	v1.ObjectMeta `json:"metadata,inline"`
	Spec          TestModelSpec `json:"spec"`
}

type TestModel2Spec struct {
	SomeProperty string `json:"someProperty"`
}
