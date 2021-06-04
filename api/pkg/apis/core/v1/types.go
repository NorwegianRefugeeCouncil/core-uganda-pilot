package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FormDefinition represents a user-defined form structure. Users can specify the logical and visual
// configuration of their forms through this data structure
type FormDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec represents the specification of this FormDefinition
	Spec FormDefinitionSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FormDefinitionList Represents a list of FormDefinition
type FormDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items represents the FormDefinition items contained in this list
	Items []FormDefinition `json:"items,omitempty" protobuf:"bytes,2,opt,name=items"`
}

// FormDefinitionSpec represents the specification of a FormDefinition
type FormDefinitionSpec struct {
	// Group represents the name of the group that the API will serve data for
	// this FormDefinition. The API will be available at the path
	// /apis/{version}/{group}/...
	Group string `json:"group,omitempty" protobuf:"bytes,1,opt,name=group"`

	// Names represents the different required names to build the API
	Names FormDefinitionNames `json:"names,omitempty" protobuf:"bytes,2,opt,name=names"`

	// Versions represents the list of versions that this FormDefinition contains
	// a FormDefinition might evolve over time, hence the need to maintain a list
	// of versions
	Versions []FormDefinitionVersion `json:"versions,omitempty" protobuf:"bytes,3,opt,name=versions"`
}

// FormDefinitionNames represent the names that the API will then serve.
// When creating a FormDefinition, an API will be available at the endpoint
// /apis/{version}/{group}/{plural}
type FormDefinitionNames struct {

	// Plural represents the lowercase, pluralized name of the resource
	// eg. formintakes, generalintakes, vulnerabilityassessments
	Plural string `json:"plural,omitempty" protobuf:"bytes,1,opt,name=plural"`

	// Singular represents the lowercase, singular name of the resource
	// eg. formintake, generalintake, vulnerabilityassessment
	Singular string `json:"singular,omitempty" protobuf:"bytes,2,opt,name=singular"`

	// Kind represents the CamelCased, singular name of the resource
	// eg. FormIntake, GeneralIntake, VulnerabilityAssessment
	Kind string `json:"kind,omitempty" protobuf:"bytes,3,opt,name=kind"`
}

// FormDefinitionVersion represents a single version of a FormDefinition
type FormDefinitionVersion struct {

	// Name represents the name of the version
	// The form definition will then present an API available at
	// /apis/{group}/{version}/{plural}
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// Served represents wheter or not this version is actually served
	// by the API
	Served bool `json:"served,omitempty" protobuf:"bytes,2,opt,name=served"`

	// Storage represents whether or not this version of the FormDefinition
	// is used to persist object in permanent storage
	Storage bool `json:"storage,omitempty" protobuf:"bytes,3,opt,name=storage"`

	// Schema represents the validation schema of the FormDefinition
	Schema FormDefinitionValidation `json:"schema,omitempty" protobuf:"bytes,4,opt,name=schema"`
}

// FormDefinitionValidation contains the validation schema of the FormDefinition
type FormDefinitionValidation struct {

	// FormScheme represents the logical and visual schema of the form
	FormSchema FormDefinitionSchema `json:"formSchema,omitempty" protobuf:"bytes,1,opt,name=formSchema"`
}

// Represents a FormDefinition element type
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

// FormDefinitionSchema represents the root form element of a form
// This is the root of the AST that contains the presents form elements
type FormDefinitionSchema struct {

	// Root represents the root of the AST that constitutes the form definition
	// A Form definition is a tree of "Containers" and "Nodes" that represent
	// the different form controls
	Root FormElementDefinition `json:"root,omitempty" protobuf:"bytes,1,opt,name=root"`
}

// TranslatedString represents the localized value of a string
type TranslatedString struct {
	// The locale name
	Locale string `json:"locale,omitempty" protobuf:"bytes,1,opt,name=locale"`

	// The translated string in the Locale locale
	Value string `json:"value,omitempty" protobuf:"bytes,1,opt,name=value"`
}

// TranslatedStrings represents a collection of TranslatedString
type TranslatedStrings []TranslatedString

// FormElementDefinition represents the configuration for a single field type
type FormElementDefinition struct {

	// Key represents the name of the property that will be filled with the value entered by the user
	// For example, if the key is "firstName", then the payload for the form will contain a property
	// "firstName"
	Key string `json:"key,omitempty" protobuf:"bytes,1,opt,name=key"`

	// Label represents the textual label presented above the field control
	Label TranslatedStrings `json:"label,omitempty" protobuf:"bytes,2,opt,name=label"`

	// Description represents the text shown under the field control
	Description TranslatedStrings `json:"description,omitempty" protobuf:"bytes,3,opt,name=description"`

	// Help represents the tooltip that the user will be able to activate to get contextual help
	Help TranslatedStrings `json:"help,omitempty" protobuf:"bytes,4,opt,name=help"`

	// Type represents the form element type
	Type FormElementType `json:"type,omitempty" protobuf:"bytes,5,opt,name=type"`

	// Required represents whether the value for this property is required or not
	Required bool `json:"required,omitempty" protobuf:"bytes,6,opt,name=required"`

	// Children contains the sub-elements of that form control. Only available for container-type elements
	Children []FormElementDefinition `json:"children,omitempty" protobuf:"bytes,7,opt,name=children"`

	// Min represents the minimum numerical value for the user input
	Min string `json:"min,omitempty" protobuf:"bytes,8,opt,name=min"`

	// Max represents the maximum numerical value for the user input
	Max string `json:"max,omitempty" protobuf:"bytes,9,opt,name=max"`

	// Pattern is a regex that will be used to validate the textual user input
	Pattern string `json:"pattern,omitempty" protobuf:"bytes,10,opt,name=pattern"`

	// MinLength represents the minimum text length for the user input
	MinLength int64 `json:"minLength,omitempty" protobuf:"bytes,11,opt,name=minLength"`

	// MaxLength represents the maximum text length for the user input
	MaxLength *int64 `json:"maxLength,omitempty" protobuf:"bytes,12,opt,name=maxLength"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomResourceDefinition represents a custom model registered by the api server.
// When registered, a CustomResourceDefinition will expose a dedicated api with the
// usual GET, PUT, POST, DELETE, ...
type CustomResourceDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec represents the specification of this FormDefinition
	Spec CustomResourceDefinitionSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomResourceDefinitionList represents a list of CustomResourceDefinition
type CustomResourceDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items represents the CustomResourceDefinition items contained in this list
	Items []CustomResourceDefinition `json:"items,omitempty" protobuf:"bytes,2,opt,name=items"`
}

// CustomResourceDefinitionSpec represents the specification of a CustomResourceDefinition
type CustomResourceDefinitionSpec struct {
	// Group represents the API group that will be exposed by the API server for this
	// CustomResourceDefinition
	Group string `json:"group" protobuf:"bytes,1,opt,name=group"`

	// Names represent the identifiers used to build the API
	Names CustomResourceDefinitionNames `json:"names" protobuf:"bytes,2,opt,name=names"`

	// Versions represents the different available api versions for that custom resource
	Versions []CustomResourceDefinitionVersion `json:"versions" protobuf:"bytes,3,opt,name=versions"`
}

// CustomResourceDefinitionNames represent the different names used to identify the resources
// of that CustomResourceDefinition API
type CustomResourceDefinitionNames struct {
	// Plural represents the lowercase, pluralized name of the resource
	// eg. formintakes, generalintakes, vulnerabilityassessments
	Plural string `json:"plural,omitempty" protobuf:"bytes,1,opt,name=plural"`

	// Singular represents the lowercase, singular name of the resource
	// eg. formintake, generalintake, vulnerabilityassessment
	Singular string `json:"singular,omitempty" protobuf:"bytes,2,opt,name=singular"`

	// Kind represents the CamelCased, singular name of the resource
	// eg. FormIntake, GeneralIntake, VulnerabilityAssessment
	Kind string `json:"kind,omitempty" protobuf:"bytes,3,opt,name=kind"`
}

// CustomResourceDefinitionVersion represent a single version of a CustomResourceDefinition
type CustomResourceDefinitionVersion struct {
	// Name represents the name of the version
	// The form definition will then present an API available at
	// /apis/{group}/{version}/{plural}
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// Served represents wheter or not this version is actually served
	// by the API
	Served bool `json:"served,omitempty" protobuf:"bytes,2,opt,name=served"`

	// Storage represents whether or not this version of the FormDefinition
	// is used to persist object in permanent storage
	Storage bool `json:"storage,omitempty" protobuf:"bytes,3,opt,name=storage"`

	// Schema contains the openAPI schema used to validate the payloads
	Schema CustomResourceDefinitionValidation `json:"schema,omitempty" protobuf:"bytes,4,opt,name=schema"`
}

type CustomResourceDefinitionValidation struct {
	OpenAPIV3Schema JSONSchemaProps `json:"openAPIV3Schema" protobuf:"bytes,1,opt,name=openAPIV3Schema"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OperatingScope represent the operating scope for
// a country, region or area
type OperatingScope struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              OperatingScopeSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type OperatingScopeSpec struct {
	AdditionalBeneficiaryInformation []AdditionalBeneficiaryInformation `json:"additionalBeneficiaryInformation,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OperatingScopeList represents a list of OrganizationOperatingScope
type OperatingScopeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items represents the OperatingScope items contained in this list
	Items []OperatingScope `json:"items,omitempty" protobuf:"bytes,2,opt,name=items"`
}

// AdditionalBeneficiaryInformation represent an OIDC claim
type AdditionalBeneficiaryInformation struct {
	Key  string `json:"key,omitempty" protobuf:"bytes,1,opt,name=key"`
	Type string `json:"type,omitempty" protobuf:"bytes,2,opt,name=type"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Attributes        map[string][]string `json:"attributes,omitempty" protobuf:"bytes,2,opt,name=attributes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items represents the User items contained in this list
	Items []User `json:"items,omitempty" protobuf:"bytes,2,opt,name=items"`
}
