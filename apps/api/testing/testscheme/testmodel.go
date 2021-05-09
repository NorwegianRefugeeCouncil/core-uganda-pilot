package testscheme

import (
	v1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/conversion"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
)

// +k8s:conversion-gen:explicit-from=net/url.Values
// +kubebuilder:object:root=true
type TestModelUrlValues struct {
	v1.TypeMeta `json:",inline"`
	Abc         string `json:"abc"`
}

// +kubebuilder:object:root=true
type TestModel struct {
	v1.TypeMeta   `json:",inline"`
	v1.ObjectMeta `json:"metadata,inline"`
	Spec          TestModelSpec `json:"spec"`
}

type TestModelSpec struct {
	SomeProperty string `json:"someProperty"`
}

// +kubebuilder:object:root=true
type TestModel2 struct {
	v1.TypeMeta   `json:",inline"`
	v1.ObjectMeta `json:"metadata,inline"`
	Spec          TestModelSpec `json:"spec"`
}

type TestModel2Spec struct {
	SomeProperty string `json:"someProperty"`
}

func AddConversionFunc(s *runtime.Scheme) error {
	return s.AddGeneratedConversionFunc(
		&TestModel{},
		&TestModel2{},
		func(a, b interface{}, scope conversion.Scope) error {
			var left = a.(*TestModel)
			var right = b.(*TestModel2)
			right.Spec.SomeProperty = left.Spec.SomeProperty
			return nil
		})
}
