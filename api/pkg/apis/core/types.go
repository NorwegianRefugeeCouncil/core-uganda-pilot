package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FormDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Spec              FormDefinitionSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FormDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []FormDefinition `json:"items" protobuf:"bytes,2,opt,name=items"`
}

type FormDefinitionSpec struct {
	Group    string                  `json:"group" protobuf:"bytes,1,opt,name=group"`
	Names    FormDefinitionNames     `json:"names" protobuf:"bytes,2,opt,name=names"`
	Versions []FormDefinitionVersion `json:"versions" protobuf:"bytes,3,opt,name=versions"`
}

type FormDefinitionNames struct {
	Plural   string `json:"plural" protobuf:"bytes,1,opt,name=plural"`
	Singular string `json:"singular" protobuf:"bytes,2,opt,name=singular"`
	Kind     string `json:"kind" protobuf:"bytes,3,opt,name=kind"`
}

type FormDefinitionVersion struct {
	Name   string                   `json:"name" protobuf:"bytes,1,opt,name=name"`
	Schema FormDefinitionValidation `json:"schema" protobuf:"bytes,2,opt,name=schema"`
}

type FormDefinitionValidation struct {
	FormSchema FormDefinitionSchema `json:"formSchema" protobuf:"bytes,1,opt,name=formSchema"`
}

type FormElementType string

const (
	SectionType = "section"
)

type FormDefinitionSchema struct {
	Root FormElementDefinition `json:"root" protobuf:"bytes,1,opt,name=root"`
}

type TranslatedString struct {
	Locale string `json:"locale" protobuf:"bytes,1,opt,name=locale"`
	Value  string `json:"value" protobuf:"bytes,1,opt,name=value"`
}

type TranslatedStrings []TranslatedString

type FormElementDefinition struct {
	Key         string                  `json:"key" protobuf:"bytes,1,opt,name=key"`
	Name        TranslatedStrings       `json:"name" protobuf:"bytes,2,opt,name=name"`
	Description TranslatedStrings       `json:"description" protobuf:"bytes,3,opt,name=description"`
	Help        TranslatedStrings       `json:"help" protobuf:"bytes,4,opt,name=help"`
	Type        FormElementType         `json:"type" protobuf:"bytes,5,opt,name=type"`
	Required    bool                    `json:"required" protobuf:"bytes,6,opt,name=required"`
	Children    []FormElementDefinition `json:"children" protobuf:"bytes,7,opt,name=children"`
}
