package apis

import (
	"encoding/json"
	"fmt"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"strings"
	"time"
)

type FormDefinition struct {
	ObjectMeta `json:"metadata,omitempty" bson:"metadata,omitempty"`
	TypeMeta   `json:",inline" bson:",inline,omitempty"`
	Spec       FormDefinitionSpec `json:"spec,omitempty" bson:"spec,omitempty"`
}

func (f *FormDefinition) String() string {
	b, err := json.MarshalIndent(f, "", "  ")
	if err == nil {
		return string(b)
	}
	return fmt.Sprintf("%#v", f)
}

type FormDefinitionList struct {
	ObjectMeta `json:"metadata,omitempty" bson:"metadata,omitempty"`
	TypeMeta   `json:",inline" bson:",inline,omitempty"`
	Items      []FormDefinition `json:"items" bson:"items"`
}

type TypeMeta struct {
	APIVersion string `json:"apiVersion,omitempty" bson:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty" bson:"kind,omitempty"`
}

func (t *TypeMeta) GetAPIVersion() string {
	return t.APIVersion
}

func (t *TypeMeta) SetAPIVersion(version string) {
	t.APIVersion = version
}

func (t *TypeMeta) GetAPIGroup() string {
	parts := strings.Split(t.APIVersion, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func (t *TypeMeta) SetAPIGroup(group string) {
	if len(t.APIVersion) == 0 {
		t.APIVersion = group + "/"
		return
	}
	parts := strings.Split(t.APIVersion, "/")
	t.APIVersion = group + "/" + parts[1]
}

func (t *TypeMeta) GetKind() string {
	return t.Kind
}

func (t *TypeMeta) SetKind(kind string) {
	t.Kind = kind
}

type ObjectMeta struct {
	UID               string     `json:"uid,omitempty" bson:"uid,omitempty"`
	ResourceVersion   int        `json:"resourceVersion,omitempty" bson:"resourceVersion,omitempty"`
	CreationTimestamp time.Time  `json:"creationTimestamp,omitempty" bson:"creationTimestamp,omitempty"`
	DeletionTimestamp *time.Time `json:"deletionTimestamp,omitempty" bson:"deletionTimestamp,omitempty"`
}

var _ runtime.Object = &FormDefinition{}

func (o *ObjectMeta) GetResourceVersion() int {
	return o.ResourceVersion
}
func (o *ObjectMeta) SetResourceVersion(version int) {
	o.ResourceVersion = version
}
func (o *ObjectMeta) GetUID() string {
	return o.UID
}

func (o *ObjectMeta) SetUID(uid string) {
	o.UID = uid
}

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
