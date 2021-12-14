package sqlmanager

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
)

// singleSelectFieldActions returns the SQL actions necessary to build a SingleSelect field
// When a user creates a form with a Single Select Field, there will be an SQL Table created
// to simply hold the different options for that field (ID and Name).
//
// We also add sqlActions required to populate this table with the options from the field definition.
//
// Also, we add a column to the Form SQL Table to hold the currently selected option.
// This allows us to add a Foreign Key to that field's options table for referential integrity.
//
// The options table looks like so
//
// Table Name: <FIELD_ID>_options
// +=======+==========+
// | id    | name     |
// +=======+==========+
// | uid1  | Option 1 |
// +-------+----------+
// | uid2  | Option 2 |
// +-------+----------+
//
func singleSelectFieldActions(formInterface types.FormInterface, fieldDefinition *types.FieldDefinition) (sqlActions, error) {
	result := sqlActions{}
	result = append(result, buildSelectOptionsTable(formInterface, fieldDefinition.ID, fieldDefinition.FieldType.SingleSelect.Options)...)

	sqlField := getStandardSQLColumnForField(fieldDefinition)
	sqlField.DataType = schema.SQLDataType{
		VarChar: &schema.SQLDataTypeVarChar{
			Length: uuidFieldLength,
		},
	}
	sqlField.Constraints = append(sqlField.Constraints, schema.SQLColumnConstraint{
		Reference: &schema.ReferenceSQLColumnConstraint{
			Schema:   formInterface.GetDatabaseID(),
			Table:    getFieldOptionsTableName(fieldDefinition.ID),
			Column:   keyIdColumn,
			OnDelete: schema.ActionCascade,
			OnUpdate: schema.ActionCascade,
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

func getFieldOptionsTableName(fieldID string) string {
	return fmt.Sprintf("%s_options", fieldID)
}

func buildSelectOptionsTable(formInterface types.FormInterface, fieldID string, options []*types.SelectOption) sqlActions {
	var result sqlActions
	// creating the SQL Table to hold the possible options for the single select field
	result = append(result, sqlAction{
		createTable: &sqlActionCreateTable{
			sqlTable: schema.SQLTable{
				Name:   getFieldOptionsTableName(fieldID),
				Schema: formInterface.GetDatabaseID(),
				Columns: []schema.SQLColumn{
					{
						Name: keyIdColumn,
						DataType: schema.SQLDataType{
							VarChar: &schema.SQLDataTypeVarChar{
								Length: uuidFieldLength,
							},
						},
						Constraints: []schema.SQLColumnConstraint{
							{
								PrimaryKey: &schema.PrimaryKeySQLColumnConstraint{},
							},
						},
					}, {
						Name: keyNameColumn,
						DataType: schema.SQLDataType{
							VarChar: &schema.SQLDataTypeVarChar{
								Length: 128,
							},
						},
						Constraints: []schema.SQLColumnConstraint{
							{
								NotNull: &schema.NotNullSQLColumnConstraint{},
							}, {
								Unique: &schema.UniqueSQLColumnConstraint{},
							},
						},
					},
				},
			},
		},
	})

	for _, option := range options {
		result = append(result, sqlAction{
			insertRow: &sqlActionInsertRow{
				schemaName: formInterface.GetDatabaseID(),
				tableName:  getFieldOptionsTableName(fieldID),
				columns:    []string{keyIdColumn, keyNameColumn},
				values:     []interface{}{option.ID, option.Name},
			},
		})
	}
	return result
}
