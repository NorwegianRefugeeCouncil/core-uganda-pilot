package defaults

import "github.com/nrc-no/core/pkg/api/types"

// FormDefinitionDefaults sets the default values for a types.FormDefinition
func FormDefinitionDefaults(definition types.FormDefinition) types.FormDefinition {

	// sets the default FormDefinition.Type
	if len(definition.Type) == 0 {
		definition.Type = types.DefaultFormType
	}

	return definition
}
