package models

type ObjectMeta struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type TypeMeta struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

type FormDefinition struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata"`
	Spec       FormDefinitionSpec `json:"spec"`
}

type FormDefinitionList struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata"`
	Items      []FormDefinition `json:"items"`
}

type FormDefinitionSpec struct {
	Group    string                  `json:"group"`
	Names    FormDefinitionNames     `json:"names"`
	Versions []FormDefinitionVersion `json:"versions"`
}

type FormDefinitionNames struct {
	Singular string `json:"singular"`
	Plural   string `json:"plural"`
	Kind     string `json:"kind"`
}

type FormDefinitionVersion struct {
	Name   string    `json:"name"`
	Schema FormModel `json:"schema"`
}

type FormModel struct {
	FormSchema FormSchema `json:"formSchema"`
}

type FormSchema struct {
	Root FormElement `json:"root"`
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
