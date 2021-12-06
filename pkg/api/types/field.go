package types

// FieldDefinition usually represents a question in a FormDefinition.
// A FieldDefinition defines the name, description, boundaries of data collection.
type FieldDefinition struct {
	// ID is the ID of the FieldDefinition
	ID string `json:"id" yaml:"id"`
	// Code is the unique Code of the FieldDefinition within the FormDefinition
	Code string `json:"code,omitempty" yaml:"code,omitempty"`
	// Name is the Name of the FieldDefinition
	Name string `json:"name" yaml:"name"`
	// Description is a helpful text helping the users to understand the question
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Options TODO: Remove this, put inside of FieldTypeMultiSelect / FieldTypeSelect
	Options []string `json:"options,omitempty" yaml:"options,omitempty"`
	// Key indicates that the FieldDefinition is part of the Unique Keys for the FormDefinition.
	// When a FormDefinition is created with Key fields, this means that there will be no
	// two records with the same combination of Key field values.
	//
	// For example, if a FormDefinition has 2 key fields, "Year" and "Month", then there could
	// be no two records with "2021" and "January".
	Key bool `json:"key" yaml:"key"`
	// Required indicates that the user must enter data for that FieldDefinition
	Required bool `json:"required" yaml:"required"`
	// FieldType contains the type of FieldDefinition
	FieldType FieldType `json:"fieldType" yaml:"fieldType"`
}

// FieldDefinitions represent a list of FieldDefinition
type FieldDefinitions []*FieldDefinition
