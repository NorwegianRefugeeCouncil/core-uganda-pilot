package sqlmanager

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
)

// quantityFieldActions returns the SQL actions necessary to create a quantity field
func quantityFieldActions(formReference types.FormReference, fieldDefinition *types.FieldDefinition) sqlActions {
	sqlField := getStandardSQLColumnForField(fieldDefinition)
	sqlField.DataType = schema.SQLDataType{
		Int: &schema.SQLDataTypeInt{},
	}
	return sqlActions{
		{
			createColumn: &sqlActionCreateColumn{
				tableName:  formReference.GetFormID(),
				schemaName: formReference.GetDatabaseID(),
				sqlColumn:  sqlField,
			},
		},
	}
}

