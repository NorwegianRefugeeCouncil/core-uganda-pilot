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

	// Handling for nullable single select fields and NULL values.
	//
	// See https://www.postgresql.org/docs/9.4/ddl-constraints.html
	//
	// Here, the sql `CHECK ("field" in ('option1','options2'))` would still evaluate to true according to the docs.
	// > It should be noted that a check constraint is satisfied if the check expression evaluates to true or the
	// > null value. Since most expressions will evaluate to the null value if any operand is null, they will not
	// > prevent null values in the constrained columns. To ensure that a column does not contain null values,
	// > the not-null constraint described in the next section can be used.
	//
	// If the field is not required, the column definition would be:
	// "fieldId" varchar(36) CHECK "fieldId" is null or "fieldId" in ('opt1','opt2')
	//
	// If the field is required, the column definition would be:
	// "fieldId" varchar(36) not null check "fieldId" in ('opt1','opt2')
	//
	canBeNullStatement := ""
	if !fieldDefinition.Required {
		canBeNullStatement = fmt.Sprintf("%s is null or ", pq.QuoteIdentifier(sqlField.Name))
	}

	checkConstraint := fmt.Sprintf("%s%s in (%s)",
		canBeNullStatement,
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
