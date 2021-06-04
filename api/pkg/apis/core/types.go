package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FormDefinition struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec FormDefinitionSpec
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FormDefinitionList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []FormDefinition
}

type FormDefinitionSpec struct {
	Group    string
	Names    FormDefinitionNames
	Versions []FormDefinitionVersion
}

type FormDefinitionNames struct {
	Plural   string
	Singular string
	Kind     string
}

type FormDefinitionVersion struct {
	Name    string
	Served  bool
	Storage bool
	Schema  FormDefinitionValidation
}

type FormDefinitionValidation struct {
	FormSchema FormDefinitionSchema
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
	Root FormElementDefinition
}

type TranslatedString struct {
	Locale string
	Value  string
}

type TranslatedStrings []TranslatedString

type FormElementDefinition struct {
	Key         string
	Label       TranslatedStrings
	Description TranslatedStrings
	Help        TranslatedStrings
	Type        FormElementType
	Required    bool
	Children    []FormElementDefinition
	Min         string
	Max         string
	Pattern     string
	MinLength   int64
	MaxLength   *int64
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomResourceDefinition represents a custom model registered by the api server.
// When registered, a CustomResourceDefinition will expose a dedicated api with the
// usual GET, PUT, POST, DELETE, ...
type CustomResourceDefinition struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	// Spec represents the specification of this FormDefinition
	Spec CustomResourceDefinitionSpec
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomResourceDefinitionList represents a list of CustomResourceDefinition
type CustomResourceDefinitionList struct {
	metav1.TypeMeta
	metav1.ListMeta

	// Items represents the CustomResourceDefinition items contained in this list
	Items []CustomResourceDefinition
}

// CustomResourceDefinitionSpec represents the specification of a CustomResourceDefinition
type CustomResourceDefinitionSpec struct {
	// Group represents the API group that will be exposed by the API server for this
	// CustomResourceDefinition
	Group string

	// Names represent the identifiers used to build the API
	Names CustomResourceDefinitionNames

	// Versions represent the api versions of that CustomResourceDefinition
	Versions []CustomResourceDefinitionVersion
}

// CustomResourceDefinitionNames represent the different names used to identify the resources
// of that CustomResourceDefinition API
type CustomResourceDefinitionNames struct {
	// Plural represents the lowercase, pluralized name of the resource
	// eg. formintakes, generalintakes, vulnerabilityassessments
	Plural string

	// Singular represents the lowercase, singular name of the resource
	// eg. formintake, generalintake, vulnerabilityassessment
	Singular string

	// Kind represents the CamelCased, singular name of the resource
	// eg. FormIntake, GeneralIntake, VulnerabilityAssessment
	Kind string
}

// CustomResourceDefinitionVersion represent a single version of a CustomResourceDefinition
type CustomResourceDefinitionVersion struct {
	// Name represents the name of the version
	// The form definition will then present an API available at
	// /apis/{group}/{version}/{plural}
	Name string

	// Served represents wheter or not this version is actually served
	// by the API
	Served bool

	// Storage represents whether or not this version of the FormDefinition
	// is used to persist object in permanent storage
	Storage bool

	// Schema contains the openAPI schema used to validate the payloads
	Schema CustomResourceDefinitionValidation
}

type CustomResourceDefinitionValidation struct {
	OpenAPIV3Schema JSONSchemaProps
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OperatingScope represent the operating scope for
// a country, region or area
type OperatingScope struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec OperatingScopeSpec
}

type OperatingScopeSpec struct {
	AdditionalBeneficiaryInformation []AdditionalBeneficiaryInformation
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OperatingScopeList represents a list of OrganizationOperatingScope
type OperatingScopeList struct {
	metav1.TypeMeta
	metav1.ListMeta

	// Items represents the OrganizationOperatingScope items contained in this list
	Items []OperatingScope
}

// AdditionalBeneficiaryInformation represent an OIDC claim
type AdditionalBeneficiaryInformation struct {
	Key  string
	Type string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type User struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Attributes map[string][]string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type UserList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// Items represents the User items contained in this list
	Items []User
}
