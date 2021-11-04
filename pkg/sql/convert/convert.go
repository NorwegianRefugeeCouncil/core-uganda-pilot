package convert

import (
	"fmt"

	"github.com/nrc-no/core/pkg/api/types"
	sqlschema2 "github.com/nrc-no/core/pkg/sql/schema"
	"github.com/nrc-no/core/pkg/utils/sets"
)

func convertDatabaseToSqlSchema(database types.Database) sqlschema2.SQLSchema {
	return sqlschema2.SQLSchema{
		Name: database.ID,
	}
}

func expandSubForms(formDef *types.FormDefinition) ([]*types.FormDefinition, error) {
	var result []*types.FormDefinition
	result = append(result, formDef)
	for _, field := range formDef.Fields {
		if field.FieldType.SubForm != nil {
			expanded, err := expandSubFormsInternal(formDef, field)
			if err != nil {
				return nil, err
			}
			result = append(result, expanded...)
		}
	}
	return result, nil
}

func expandSubFormsInternal(parentForm *types.FormDefinition, subFormField *types.FieldDefinition) ([]*types.FormDefinition, error) {
	var result []*types.FormDefinition

	formDef := &types.FormDefinition{
		ID:         subFormField.ID,
		Name:       subFormField.Name,
		DatabaseID: parentForm.DatabaseID,
		Fields:     subFormField.FieldType.SubForm.Fields,
	}

	if _, err := parentForm.RemoveFieldByID(subFormField.ID); err != nil {
		return nil, err
	}

	formDef.Fields = append(formDef.Fields, &types.FieldDefinition{
		ID:       "owner_id",
		Code:     "owner_id",
		Name:     "owner_id",
		Required: true,
		FieldType: types.FieldType{
			Reference: &types.FieldTypeReference{
				DatabaseID: parentForm.DatabaseID,
				FormID:     parentForm.ID,
			},
		},
	})

	result = append(result, formDef)

	for _, field := range formDef.Fields {
		if field.FieldType.SubForm != nil {
			expanded, err := expandSubFormsInternal(formDef, field)
			if err != nil {
				return nil, err
			}
			result = append(result, expanded...)
		}
	}

	return result, nil
}

func convertFormToSqlTable(formDef *types.FormDefinition, referencedForms *types.FormDefinitionList) (sqlschema2.SQLTable, error) {
	table := sqlschema2.SQLTable{}

	table.Name = formDef.ID
	table.Schema = formDef.DatabaseID

	table.Columns = append(table.Columns, sqlschema2.SQLColumn{
		Name: "id",
		DataType: sqlschema2.SQLDataType{
			VarChar: &sqlschema2.SQLDataTypeVarChar{Length: 36},
		},
		Constraints: []sqlschema2.SQLColumnConstraint{
			{
				Name:       fmt.Sprintf("pk_%s_%s", formDef.DatabaseID, formDef.ID),
				PrimaryKey: &sqlschema2.PrimaryKeySQLColumnConstraint{},
			},
		},
	})

	table.Columns = append(table.Columns, sqlschema2.SQLColumn{
		Name: "seq",
		DataType: sqlschema2.SQLDataType{
			Serial: &sqlschema2.SQLDataTypeSerial{},
		},
	})

	table.Columns = append(table.Columns, sqlschema2.SQLColumn{
		Name:    "created_at",
		Default: "NOW()",
		DataType: sqlschema2.SQLDataType{
			Timestamp: &sqlschema2.SQLDataTypeTimestamp{
				Timezone: &sqlschema2.TimestampWithoutTimeZone,
			},
		},
		Constraints: []sqlschema2.SQLColumnConstraint{
			{
				NotNull: &sqlschema2.NotNullSQLColumnConstraint{},
			},
		},
	})

	keyFieldIDs := sets.NewString()
	for _, field := range formDef.Fields {
		table.Columns = append(table.Columns, convertFieldToSqlField(formDef, field))
		if field.Key {
			keyFieldIDs.Insert(field.ID)
		}
	}

	if !keyFieldIDs.IsEmpty() {
		table.Constraints = append(table.Constraints, sqlschema2.SQLTableConstraint{
			Name: fmt.Sprintf("uq_%s", table.Name),
			Unique: &sqlschema2.SQLTableConstraintUnique{
				ColumnNames: keyFieldIDs.List(),
			},
		})
	}

	return table, nil
}

func convertFieldToSqlField(formDef *types.FormDefinition, field *types.FieldDefinition) sqlschema2.SQLColumn {
	result := sqlschema2.SQLColumn{}
	result.Name = field.ID
	result.Comment = field.Code

	if field.Required || field.Key {
		result.Constraints = append(result.Constraints, sqlschema2.SQLColumnConstraint{
			NotNull: &sqlschema2.NotNullSQLColumnConstraint{},
		})
	}

	if field.FieldType.Text != nil {
		result.DataType = sqlschema2.SQLDataType{
			VarChar: &sqlschema2.SQLDataTypeVarChar{
				Length: 1024,
			},
		}
	} else if field.FieldType.Week != nil {
		result.DataType = sqlschema2.SQLDataType{
			Date: &sqlschema2.SQLDataTypeDate{},
		}
	} else if field.FieldType.MultilineText != nil {
		result.DataType = sqlschema2.SQLDataType{
			Text: &sqlschema2.SQLDataTypeText{},
		}
	} else if field.FieldType.Date != nil {
		result.DataType = sqlschema2.SQLDataType{
			Date: &sqlschema2.SQLDataTypeDate{},
		}
	} else if field.FieldType.Month != nil {
		result.DataType = sqlschema2.SQLDataType{
			Date: &sqlschema2.SQLDataTypeDate{},
		}
	} else if field.FieldType.Quantity != nil {
		result.DataType = sqlschema2.SQLDataType{
			Int: &sqlschema2.SQLDataTypeInt{},
		}
	} else if field.FieldType.SingleSelect != nil {
		result.DataType = sqlschema2.SQLDataType{
			VarChar: &sqlschema2.SQLDataTypeVarChar{
				Length: 1024,
			},
		}
	} else if field.FieldType.Reference != nil {
		result.Constraints = append(result.Constraints, sqlschema2.SQLColumnConstraint{
			Name: fmt.Sprintf("fkref__%s_%s__%s_id",
				formDef.ID,
				field.ID,
				field.FieldType.Reference.FormID,
			),
			Reference: &sqlschema2.ReferenceSQLColumnConstraint{
				Schema:   field.FieldType.Reference.DatabaseID,
				Table:    field.FieldType.Reference.FormID,
				Column:   "id",
				OnUpdate: sqlschema2.ActionCascade,
				OnDelete: sqlschema2.ActionCascade,
			},
		})
		result.DataType = sqlschema2.SQLDataType{
			VarChar: &sqlschema2.SQLDataTypeVarChar{Length: 36},
		}
	} else if field.FieldType.SubForm != nil {
		result.Name += "_id"
		result.DataType = sqlschema2.SQLDataType{
			VarChar: &sqlschema2.SQLDataTypeVarChar{Length: 36},
		}
		result.Constraints = append(result.Constraints, sqlschema2.SQLColumnConstraint{
			Reference: &sqlschema2.ReferenceSQLColumnConstraint{
				Schema:    formDef.DatabaseID,
				Table:     field.ID,
				Column:    "id",
				OnDelete:  sqlschema2.ActionRestrict,
				OnUpdate:  sqlschema2.ActionCascade,
				MatchType: "",
			},
		})
	}

	return result
}
