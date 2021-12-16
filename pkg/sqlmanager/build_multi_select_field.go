package sqlmanager

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
	"strings"
)

// multiSelectFieldActions returns the SQL actions necessary to build a MultiSelect field
// When a user creates a form with a Multi Select Field, it will add a column to the form's table
// This column will be of varchar(36)[] type (array of varchar(36))
//
// Also, we add an SQL Check constraint to make sure that the values passed for the field
// will be values that correspond to the available options.
//
// "MyField" varchar(36)[] CHECK "MyField" <@ array['option1','option2']::varchar[]
//
// For non-required (nullable) fields,
// "MyField" varchar(36)[] CHECK "MyField" is null or "MyField" <@ array['option1','option2']::varchar[]
//
//
func multiSelectFieldActions(formInterface types.FormInterface, fieldDefinition *types.FieldDefinition) (sqlActions, error) {
	result := sqlActions{}
	sqlField := getStandardSQLColumnForField(fieldDefinition)
	sqlField.DataType = schema.SQLDataType{
		Array: &schema.SQLDataTypeArray{
			DataType: schema.SQLDataType{
				VarChar: &schema.SQLDataTypeVarChar{
					Length: uuidFieldLength,
				},
			},
		},
	}

	var isNullConstraint = ""
	if !fieldDefinition.Required {
		isNullConstraint = fmt.Sprintf("%s is null or ",
			pq.QuoteIdentifier(sqlField.Name))
	}

	var availableOptions []string
	for _, option := range fieldDefinition.FieldType.MultiSelect.Options {
		availableOptions = append(availableOptions, pq.QuoteLiteral(option.ID))
	}

	checkConstraint := fmt.Sprintf("%s%s <@ array[%s]::varchar[]",
		isNullConstraint,
		pq.QuoteIdentifier(sqlField.Name),
		strings.Join(availableOptions, ","),
	)

	sqlField.Constraints = append(sqlField.Constraints,
		schema.SQLColumnConstraint{
			Check: &schema.CheckSQLColumnConstraint{
				Expression: checkConstraint,
			},
		})

	result = append(result, sqlAction{
		createColumn: &sqlActionCreateColumn{
			tableName:  formInterface.GetFormID(),
			schemaName: formInterface.GetDatabaseID(),
			sqlColumn:  sqlField,
		},
	})
	return result, nil
}
