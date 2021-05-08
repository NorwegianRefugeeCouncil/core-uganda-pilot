package v1

import (
	"encoding/json"
	"fmt"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
)

//go:generate controller-gen object paths=$GOFILE

// +kubebuilder:object:root=true
type FormDefinition struct {
	metav1.TypeMeta   `json:",inline" bson:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Spec              FormDefinitionSpec `json:"spec,omitempty" bson:"spec,omitempty"`
}

var _ runtime.Object = &FormDefinition{}

func (f *FormDefinition) String() string {
	b, err := json.MarshalIndent(f, "", "  ")
	if err == nil {
		return string(b)
	}
	return fmt.Sprintf("%#v", f)
}

// +kubebuilder:object:root=true
type FormDefinitionList struct {
	metav1.ListMeta `json:"metadata,omitempty" bson:"metadata,omitempty"`
	metav1.TypeMeta `json:",inline" bson:",inline,omitempty"`
	Items           []FormDefinition `json:"items" bson:"items"`
}

var _ runtime.Object = &FormDefinitionList{}

// +kubebuilder:object:generate=true
type FormDefinitionSpec struct {
	Group    string                  `json:"group,omitempty" bson:"group,omitempty"`
	Names    CustomResourceNames     `json:"names,omitempty" bson:"names,omitempty"`
	Versions []FormDefinitionVersion `json:"versions" bson:"versions,omitempty"`
}

// +kubebuilder:object:generate=true
type FormDefinitionVersion struct {
	Name   string
	Schema FormSchema `json:"schema" bson:"schema,omitempty"`
}

// +kubebuilder:object:generate=true
type CustomResourceNames struct {
	Plural   string `json:"plural,omitempty" bson:"plural,omitempty"`
	Singular string `json:"singular,omitempty" bson:"singular,omitempty"`
	Kind     string `json:"kind,omitempty" bson:"kind,omitempty"`
}

// +kubebuilder:object:generate=true
type FormSchema struct {
	FormSchema FormSchemaDefinition `json:"formSchema,omitempty" bson:"formSchema,omitempty"`
}

// +kubebuilder:object:generate=true
type FormSchemaDefinition struct {
	Root FormElement `json:"root,omitempty" bson:"root"`
}

// +kubebuilder:object:generate=true
type FormElement struct {
	Key         string             `json:"key,omitempty" bson:"key,omitempty"`
	ID          string             `json:"id,omitempty" bson:"id,omitempty"`
	Name        []TranslatedString `json:"name,omitempty" bson:"name,omitempty"`
	Type        string             `json:"type,omitempty" bson:"type,omitempty"`
	Description []TranslatedString `json:"description,omitempty" bson:"description,omitempty"`
	Children    []FormElement      `json:"children,omitempty" bson:"children,omitempty"`
}

// +kubebuilder:object:generate=true
type TranslatedString struct {
	Locale string `json:"locale,omitempty" bson:"locale,omitempty"`
	Value  string `json:"value,omitempty" bson:"value,omitempty"`
}
