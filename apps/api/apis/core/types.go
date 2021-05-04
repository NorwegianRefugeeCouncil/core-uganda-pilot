package core

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

//go:generate controller-gen object paths=$GOFILE

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Model struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ModelSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ModelList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

type ModelSpec struct {
	Group string     `json:"group"`
	Names ModelNames `json:"names"`
}

type ModelNames struct {
	Singular string `json:"singular"`
	Plural   string `json:"plural"`
	Kind     string `json:"kind"`
}

type ModelVersion struct {
	Name   string      `json:"name"`
	Schema ModelSchema `json:"schema"`
}

type ModelSchema struct {
	FormSchema FormSchema `json:"formSchema"`
}

type FormSchema struct {
	Root FormElement
}

type FormElementType string

const (
	SectionElement = "section"
	TextElement    = "text"
)

type FormElement struct {
	ID       string          `json:"id"`
	Children []FormElement   `json:"children"`
	Type     FormElementType `json:"type"`
}
