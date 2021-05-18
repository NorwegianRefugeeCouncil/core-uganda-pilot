package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FormDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              FormDefinitionSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FormDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []FormDefinition `json:"items,omitempty" protobuf:"bytes,2,opt,name=items"`
}

type FormDefinitionSpec struct {
	Group    string                  `json:"group,omitempty" protobuf:"bytes,1,opt,name=group"`
	Names    FormDefinitionNames     `json:"names,omitempty" protobuf:"bytes,2,opt,name=names"`
	Versions []FormDefinitionVersion `json:"versions,omitempty" protobuf:"bytes,3,opt,name=versions"`
}

type FormDefinitionNames struct {
	Plural   string `json:"plural,omitempty" protobuf:"bytes,1,opt,name=plural"`
	Singular string `json:"singular,omitempty" protobuf:"bytes,2,opt,name=singular"`
	Kind     string `json:"kind,omitempty" protobuf:"bytes,3,opt,name=kind"`
}

type FormDefinitionVersion struct {
	Name   string                   `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Schema FormDefinitionValidation `json:"schema,omitempty" protobuf:"bytes,2,opt,name=schema"`
}

type FormDefinitionValidation struct {
	FormSchema FormDefinitionSchema `json:"formSchema,omitempty" protobuf:"bytes,1,opt,name=formSchema"`
}

type FormElementType string

const (
	SectionType   FormElementType = "section"
	ShortTextType FormElementType = "shortText"
	LongTextType  FormElementType = "longText"
	IntegerType   FormElementType = "integer"
	SelectType    FormElementType = "select"
	DateType      FormElementType = "date"
	DateTimeType  FormElementType = "dateTime"
	TimeType      FormElementType = "time"
)

type FormDefinitionSchema struct {
	Root FormElementDefinition `json:"root,omitempty" protobuf:"bytes,1,opt,name=root"`
}

type TranslatedString struct {
	Locale string `json:"locale,omitempty" protobuf:"bytes,1,opt,name=locale"`
	Value  string `json:"value,omitempty" protobuf:"bytes,1,opt,name=value"`
}

type TranslatedStrings []TranslatedString

type FormElementDefinition struct {
	Key         string                  `json:"key,omitempty" protobuf:"bytes,1,opt,name=key"`
	Label       TranslatedStrings       `json:"label,omitempty" protobuf:"bytes,2,opt,name=label"`
	Description TranslatedStrings       `json:"description,omitempty" protobuf:"bytes,3,opt,name=description"`
	Help        TranslatedStrings       `json:"help,omitempty" protobuf:"bytes,4,opt,name=help"`
	Type        FormElementType         `json:"type,omitempty" protobuf:"bytes,5,opt,name=type"`
	Required    bool                    `json:"required,omitempty" protobuf:"bytes,6,opt,name=required"`
	Children    []FormElementDefinition `json:"children,omitempty" protobuf:"bytes,7,opt,name=children"`
	Min         string                  `json:"min,omitempty" protobuf:"bytes,8,opt,name=min"`
	Max         string                  `json:"max,omitempty" protobuf:"bytes,9,opt,name=max"`
	Pattern     string                  `json:"pattern,omitempty" protobuf:"bytes,10,opt,name=pattern"`
	MinLength   int64                   `json:"minLength,omitempty" protobuf:"bytes,11,opt,name=minLength"`
	MaxLength   *int64                  `json:"maxLength,omitempty" protobuf:"bytes,12,opt,name=maxLength"`
}
