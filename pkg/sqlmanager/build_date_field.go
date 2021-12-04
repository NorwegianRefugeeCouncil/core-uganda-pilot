package sqlmanager

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
)

// dateFieldActions returns the SQL actions necessary to create a date field
// a date field is stored as a regular SQL date. Though, the record store will coerce the
// value of the stored date to be the first day of the week.
func dateFieldActions(formReference types.FormReference, fieldDefinition *types.FieldDefinition) sqlActions {
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

