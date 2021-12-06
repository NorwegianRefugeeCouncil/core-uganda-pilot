package sqlmanager

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
)

// referenceFieldActions returns the SQL actions necessary to create a reference field
// a reference field adds a sql column to the sql table, as well as creating a foreign key
// to the referred table. This ensures that we maintain referential integrity, and that
// users cannot enter a reference for a record that does not exist on the referred table.
func referenceFieldActions(formReference types.FormReference, fieldDefinition *types.FieldDefinition) sqlActions {
	sqlField := getStandardSQLColumnForField(fieldDefinition)
	sqlField.DataType = schema.SQLDataType{
		VarChar: &schema.SQLDataTypeVarChar{
			Length: 36,
		},
	}

	// If the reference field is not required, if the referred record is deleted
	// we use the SET NULL action. Otherwise, we restrict deleting the referred record.
	onDeleteAction := schema.ActionSetNull
	if fieldDefinition.Required || fieldDefinition.Key {
		onDeleteAction = schema.ActionRestrict
	}

	// Adding the foreign key constraint
	sqlField.Constraints = append(sqlField.Constraints, schema.SQLColumnConstraint{
		Reference: &schema.ReferenceSQLColumnConstraint{
			Schema:   fieldDefinition.FieldType.Reference.DatabaseID,
			Table:    fieldDefinition.FieldType.Reference.FormID,
			Column:   "id",
			OnDelete: onDeleteAction,
			OnUpdate: schema.ActionCascade,
		},
	})

	return []sqlAction{
		{
			createColumn: &sqlActionCreateColumn{
				tableName:  formReference.GetFormID(),
				schemaName: formReference.GetDatabaseID(),
				sqlColumn:  sqlField,
			},
		},
	}
}
