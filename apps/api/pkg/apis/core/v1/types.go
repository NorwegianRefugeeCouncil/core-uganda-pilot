package v1

import (
	"encoding/json"
	"fmt"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FormDefinitionList struct {
	metav1.ListMeta `json:"metadata,omitempty" bson:"metadata,omitempty"`
	metav1.TypeMeta `json:",inline" bson:",inline,omitempty"`
	Items           []FormDefinition `json:"items" bson:"items"`
}

var _ runtime.Object = &FormDefinitionList{}

type FormDefinitionSpec struct {
	Group    string                  `json:"group,omitempty" bson:"group,omitempty"`
	Names    CustomResourceNames     `json:"names,omitempty" bson:"names,omitempty"`
	Versions []FormDefinitionVersion `json:"versions" bson:"versions,omitempty"`
}

type FormDefinitionVersion struct {
	Name   string
	Schema FormSchema `json:"schema" bson:"schema,omitempty"`
}

type CustomResourceNames struct {
	Plural   string `json:"plural,omitempty" bson:"plural,omitempty"`
	Singular string `json:"singular,omitempty" bson:"singular,omitempty"`
	Kind     string `json:"kind,omitempty" bson:"kind,omitempty"`
}

type FormSchema struct {
	FormSchema FormSchemaDefinition `json:"formSchema,omitempty" bson:"formSchema,omitempty"`
}

type FormSchemaDefinition struct {
	Root FormElement `json:"root,omitempty" bson:"root"`
}

type FormElement struct {
	Key         string             `json:"key,omitempty" bson:"key,omitempty"`
	ID          string             `json:"id,omitempty" bson:"id,omitempty"`
	Name        []TranslatedString `json:"name,omitempty" bson:"name,omitempty"`
	Type        string             `json:"type,omitempty" bson:"type,omitempty"`
	Description []TranslatedString `json:"description,omitempty" bson:"description,omitempty"`
	Children    []FormElement      `json:"children,omitempty" bson:"children,omitempty"`
}

type TranslatedString struct {
	Locale string `json:"locale,omitempty" bson:"locale,omitempty"`
	Value  string `json:"value,omitempty" bson:"value,omitempty"`
}
