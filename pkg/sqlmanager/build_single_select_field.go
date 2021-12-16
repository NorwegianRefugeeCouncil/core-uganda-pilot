package sqlmanager

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
	"strings"
)

// singleSelectFieldActions returns the SQL actions necessary to build a SingleSelect field
// When a user creates a form with a Single Select Field, this will create a VARCHAR(36) column
// on the form's SQL Table to contain the selected option ID.
//
// Additionally, we add a CHECK constraint on that column to make sure that the value provided
// is one of the available option IDs.
//
func singleSelectFieldActions(formInterface types.FormInterface, fieldDefinition *types.FieldDefinition) (sqlActions, error) {
	result := sqlActions{}
	sqlField := getStandardSQLColumnForField(fieldDefinition)
	sqlField.DataType = schema.SQLDataType{
		VarChar: &schema.SQLDataTypeVarChar{
			Length: uuidFieldLength,
		},
	}

	var availableOptions []string
	for _, option := range fieldDefinition.FieldType.SingleSelect.Options {
		availableOptions = append(availableOptions, pq.QuoteLiteral(option.ID))
	}

	checkConstraint := fmt.Sprintf("%s IN (%s)",
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
