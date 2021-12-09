package seeder

import (
	"github.com/nrc-no/core/pkg/api/types"
)

func newFieldDefinition(name, description string, key, required bool, fieldType types.FieldType) *types.FieldDefinition {
	return &types.FieldDefinition{
		Name:        name,
		Description: description,
		Key:         key,
		Required:    required,
		FieldType:   fieldType,
	}
}

func newFormDefinition(databaseId, folderId, name string, fields []*types.FieldDefinition) *types.FormDefinition {
	return &types.FormDefinition{
		DatabaseID: databaseId,
		FolderID:   folderId,
		Name:       name,
		Fields:     fields,
	}
}
