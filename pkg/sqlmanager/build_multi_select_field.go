package sqlmanager

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
)

// multiSelectFieldActions returns the SQL actions necessary to build a MultiSelect field
// When a user creates a form with a Multi Select Field, there will be two SQL Tables created
//
// The first table is the same table as the Single Select field, containing the available options.
// See singleSelectFieldActions
//
// The second table will be an association table to link a record with an associated option
// We also add sqlActions required to populate this table with the options from the field definition.
//
// Table Name: <FIELD_ID>_associations
// +=======+===========+
// | id    | option_id |
// +=======+===========+
// | uid1  | Option 1  |
// +-------+-----------+
// | uid1  | Option 2  |
// +-------+-----------+
// | uid2  | Option 2  |
// +-------+-----------+
// | uid2  | Option 2  |
// +-------+-----------+
// | uid2  | Option 3  |
// +-------+-----------+
//
func multiSelectFieldActions(formInterface types.FormInterface, fieldDefinition *types.FieldDefinition) (sqlActions, error) {
	result := sqlActions{}
	result = append(result, buildSelectOptionsTable(formInterface, fieldDefinition.ID, fieldDefinition.FieldType.MultiSelect.Options)...)
	result = append(result, buildMultiSelectAssociationTable(formInterface, fieldDefinition)...)
	return result, nil
}

func buildMultiSelectAssociationTable(formInterface types.FormInterface, fieldDefinition *types.FieldDefinition) sqlActions {
	var result sqlActions
	// creating the SQL Table to hold the associations between a record and the available options
	result = append(result, sqlAction{
		createTable: &sqlActionCreateTable{
			sqlTable: schema.SQLTable{
				Name:   getMultiSelectAssociationTableName(fieldDefinition.ID),
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
								NotNull: &schema.NotNullSQLColumnConstraint{},
							},
							{
								Reference: &schema.ReferenceSQLColumnConstraint{
									Schema:   formInterface.GetDatabaseID(),
									Table:    formInterface.GetFormID(),
									Column:   keyIdColumn,
									OnDelete: schema.ActionCascade,
									OnUpdate: schema.ActionCascade,
								},
							},
						},
					}, {
						Name: keyOptionIdColumn,
						DataType: schema.SQLDataType{
							VarChar: &schema.SQLDataTypeVarChar{
								Length: uuidFieldLength,
							},
						},
						Constraints: []schema.SQLColumnConstraint{
							{
								NotNull: &schema.NotNullSQLColumnConstraint{},
							},
							{
								Reference: &schema.ReferenceSQLColumnConstraint{
									Schema:   formInterface.GetDatabaseID(),
									Table:    getFieldOptionsTableName(fieldDefinition.ID),
									Column:   keyIdColumn,
									OnDelete: schema.ActionCascade,
									OnUpdate: schema.ActionCascade,
								},
							},
						},
					},
				},
				Constraints: []schema.SQLTableConstraint{
					{
						Name: fmt.Sprintf("uk_key_%s", fieldDefinition.ID),
						Unique: &schema.SQLTableConstraintUnique{
							ColumnNames: []string{keyIdColumn, keyOptionIdColumn},
						},
					},
				},
			},
		},
	})
	return result
}

func getMultiSelectAssociationTableName(fieldID string) string {
	return fmt.Sprintf("%s_associations", fieldID)
}
