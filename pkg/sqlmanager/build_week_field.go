package sqlmanager

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
)

// weekFieldActions returns the SQL actions necessary to create a week field
// a week field is stored as a regular SQL Date.
func weekFieldActions(formReference types.FormReference, fieldDefinition *types.FieldDefinition) sqlActions {
	sqlField := getStandardSQLColumnForField(fieldDefinition)
	sqlField.DataType = schema.SQLDataType{
		Date: &schema.SQLDataTypeDate{},
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
