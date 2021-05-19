package helpers

import (
	"github.com/nrc-no/core/api/pkg/apis/core/v1"
	v12 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

func ConvertToFormDefinition(formDefinition *v1.FormDefinition) *v12.CustomResourceDefinition {

	crd := &v12.CustomResourceDefinition{
		ObjectMeta: v13.ObjectMeta{
			Name: formDefinition.Spec.Names.Plural + "." + formDefinition.Spec.Group,
		},
		Spec: v12.CustomResourceDefinitionSpec{
			Group: formDefinition.Spec.Group,
			Scope: v12.ClusterScoped,
			Conversion: &v12.CustomResourceConversion{
				Strategy: v12.NoneConverter,
			},
			Names: v12.CustomResourceDefinitionNames{
				Plural:   formDefinition.Spec.Names.Plural,
				Singular: formDefinition.Spec.Names.Singular,
				Kind:     formDefinition.Spec.Names.Kind,
				ListKind: formDefinition.Spec.Names.Kind + "List",
			},
		},
	}

	for _, version := range formDefinition.Spec.Versions {
		crdVersion := v12.CustomResourceDefinitionVersion{
			Name:    version.Name,
			Storage: version.Storage,
			Served:  version.Served,
			Schema:  &v12.CustomResourceValidation{},
		}
		validation := ConvertFormDefinitionVersion(formDefinition, version)
		crdVersion.Schema = validation
		crd.Spec.Versions = append(crd.Spec.Versions, crdVersion)
	}

	return crd
}

func ConvertFormDefinitionVersion(fs *v1.FormDefinition, version v1.FormDefinitionVersion) *v12.CustomResourceValidation {

	specSchema := &v12.JSONSchemaProps{
		Description: `Defines the desired state fo ` + fs.Spec.Names.Kind,
		Type:        "object",
	}
	formSchema := version.Schema.FormSchema.Root
	WalkFormSchema(formSchema, specSchema)

	return &v12.CustomResourceValidation{
		OpenAPIV3Schema: &v12.JSONSchemaProps{
			Description: "Schema for the " + fs.Spec.Names.Kind + " api",
			Type:        "object",
			Properties: map[string]v12.JSONSchemaProps{
				"apiVersion": {
					Description: `APIVersion defines the versioned schema of this representation
of an object. Servers should convert recognized schemas to the latest internal value, and may
reject unrecognized values.`,
					Type: "string",
				},
				"kind": {
					Description: `Kind is a string value representing the REST resource this 
object represents. Servers may infer this from the endpoint the client submits requests to.
Cannot be updated. In CamelCase.`,
					Type: "string",
				},
				"metadata": {
					Type: "object",
				},
				"spec": *specSchema,
			},
		},
	}
}

func WalkFormSchema(element v1.FormElementDefinition, jsonProps *v12.JSONSchemaProps) {

	var intMultipleOf float64 = 1
	switch element.Type {
	case v1.IntegerType:
		jsonProps.Type = "number"
		jsonProps.MultipleOf = &intMultipleOf
	case v1.ShortTextType:
		jsonProps.Type = "string"
	case v1.LongTextType:
		jsonProps.Type = "string"
	case v1.SectionType:
		jsonProps.Type = "object"
	case v1.DateTimeType:
		jsonProps.Type = "string"
		jsonProps.Format = "datetime"
	case v1.DateType:
		jsonProps.Type = "string"
		jsonProps.Format = "date"
	case v1.SelectType:
	case v1.TimeType:
	}

	if element.MinLength != 0 {
		jsonProps.MinLength = &element.MinLength
	}
	if element.MaxLength != nil {
		jsonProps.MaxLength = element.MaxLength
	}
	if element.Max != "" {
		max, err := strconv.ParseFloat(element.Max, 64)
		if err == nil {
			jsonProps.Maximum = &max
		}
	}
	if element.Min != "" {
		min, err := strconv.ParseFloat(element.Min, 64)
		if err == nil {
			jsonProps.Minimum = &min
		}
	}
	if element.Pattern != "" {
		jsonProps.Pattern = element.Pattern
	}

	if jsonProps.Description == "" {
		jsonProps.Description = findDescription(element.Description)
	}

	for _, child := range element.Children {
		childJsonProps := &v12.JSONSchemaProps{}
		WalkFormSchema(child, childJsonProps)

		if childJsonProps.Type == "object" {
			if childJsonProps.Properties != nil {
				if jsonProps.Properties == nil {
					jsonProps.Properties = map[string]v12.JSONSchemaProps{}
				}
				for key, props := range childJsonProps.Properties {
					jsonProps.Properties[key] = props
				}
				for _, propName := range childJsonProps.Required {
					jsonProps.Required = append(jsonProps.Required, propName)
				}
			}
		} else {
			if jsonProps.Properties == nil {
				jsonProps.Properties = map[string]v12.JSONSchemaProps{}
			}
			jsonProps.Properties[child.Key] = *childJsonProps
			if child.Required {
				jsonProps.Required = append(jsonProps.Required, child.Key)
			}
		}
	}
}

func findDescription(strs v1.TranslatedStrings) string {
	for _, str := range strs {
		if str.Locale == "en" {
			return str.Value
		}
	}
	if len(strs) > 0 {
		return strs[0].Value
	}
	return ""
}
