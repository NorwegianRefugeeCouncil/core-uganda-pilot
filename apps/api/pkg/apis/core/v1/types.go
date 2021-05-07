package v1

import (
	"encoding/json"
	"fmt"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
)

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

type FormDefinitionList struct {
	metav1.ListMeta `json:"metadata,omitempty" bson:"metadata,omitempty"`
	metav1.TypeMeta `json:",inline" bson:",inline,omitempty"`
	Items           []FormDefinition `json:"items" bson:"items"`
}

var _ runtime.Object = &FormDefinitionList{}

type FormDefinitionSpec struct {
	Group    string                  `json:"group,omitempty"`
	Names    CustomResourceNames     `json:"names,omitempty"`
	Versions []FormDefinitionVersion `json:"versions"`
}

type FormDefinitionVersion struct {
	Name   string
	Schema FormSchema `json:"schema"`
}

type CustomResourceNames struct {
	Plural   string `json:"plural,omitempty"`
	Singular string `json:"singular,omitempty"`
	Kind     string `json:"kind,omitempty"`
}

type FormSchema struct {
	FormSchema FormSchemaDefinition `json:"formSchema,omitempty"`
}

type FormSchemaDefinition struct {
	Root FormElement `json:"root,omitempty"`
}

type FormElement struct {
	Key         string             `json:"key,omitempty"`
	ID          string             `json:"id,omitempty"`
	Name        []TranslatedString `json:"name,omitempty"`
	Type        string             `json:"type,omitempty"`
	Description []TranslatedString `json:"description,omitempty"`
	Children    []FormElement      `json:"children,omitempty"`
}

type TranslatedString struct {
	Locale string `json:"locale,omitempty"`
	Value  string `json:"value,omitempty"`
}
