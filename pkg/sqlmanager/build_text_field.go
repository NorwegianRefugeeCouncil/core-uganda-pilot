package sqlmanager

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
)

// textFieldActions returns the SQL actions necessary to create a text field
// a text field is stored as a VarChar field with a maximum length of 1024.
func textFieldActions(formReference types.FormReference, fieldDefinition *types.FieldDefinition) sqlActions {
	sqlField := getStandardSQLColumnForField(fieldDefinition)
	sqlField.DataType = schema.SQLDataType{
		VarChar: &schema.SQLDataTypeVarChar{
			Length: 1024,
		},
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
