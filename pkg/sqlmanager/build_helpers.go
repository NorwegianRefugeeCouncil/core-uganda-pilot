package sqlmanager

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
	"github.com/nrc-no/core/pkg/utils/sets"
)

// buildKeyFieldActions creates the SQL Actions necessary to configure an SQL Unique Constraints
// so that they represent a form's key fields.
// This mechanism will ensure that we maintain unique values for the form's "key" fields.
func buildKeyFieldActions(formInterface types.FormInterface) sqlActions {
	keyFieldIDs := sets.NewString()
	for _, fieldDefinition := range formInterface.GetFields() {
		if fieldDefinition.Key {
			keyFieldIDs.Insert(fieldDefinition.ID)
		}
	}
	result := sqlActions{}
	if keyFieldIDs.IsEmpty() {
		return result
	}

	result = append(result, sqlAction{
		createUniqueConstraint: &sqlActionCreateConstraint{
			tableName:  formInterface.GetFormID(),
			schemaName: formInterface.GetDatabaseID(),
			sqlConstraint: schema.SQLTableConstraint{
				Name: "uk_key_" + formInterface.GetFormID(),
				Unique: &schema.SQLTableConstraintUnique{
					ColumnNames: keyFieldIDs.List(),
				},
			},
		},
	})

	return result
}

// buildTableIDColumn is a helper method that creates an "id" SQL column that all tables must have.
// the ID column is also configured to be the primary key
func buildTableIDColumn() schema.SQLColumn {
	return schema.SQLColumn{
		Name: "id",
		DataType: schema.SQLDataType{
			VarChar: &schema.SQLDataTypeVarChar{
				Length: 36,
			},
		},
		Constraints: []schema.SQLColumnConstraint{
			{
				PrimaryKey: &schema.PrimaryKeySQLColumnConstraint{},
			},
		},
	}
}

// buildTableCreatedAtColumn is a helper method that creates a "created_at" SQL column that all tables must have.
// The value for this column has a default value of NOW()
func buildTableCreatedAtColumn() schema.SQLColumn {
	return schema.SQLColumn{
		Name:    "created_at",
		Default: "NOW()",
		DataType: schema.SQLDataType{
			Timestamp: &schema.SQLDataTypeTimestamp{
				Timezone: &schema.TimestampWithTimeZone,
			},
		},
		Constraints: []schema.SQLColumnConstraint{
			{
				NotNull: &schema.NotNullSQLColumnConstraint{},
			},
		},
	}
}

// getStandardSQLColumnForField returns a standard SQL sqlColumn configuration that applies to most field types.
func getStandardSQLColumnForField(fieldDefinition *types.FieldDefinition) schema.SQLColumn {
	sqlField := schema.SQLColumn{
		Name: fieldDefinition.ID,
	}
	if fieldDefinition.Required || fieldDefinition.Key {
		sqlField.Constraints = append(sqlField.Constraints, schema.SQLColumnConstraint{
			NotNull: &schema.NotNullSQLColumnConstraint{},
		})
	}
	return sqlField
}
