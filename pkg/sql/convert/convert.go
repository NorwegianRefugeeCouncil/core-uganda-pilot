package convert

import (
	"fmt"

	"github.com/nrc-no/core/pkg/sqlschema"
	"github.com/nrc-no/core/pkg/types"
	"github.com/nrc-no/core/pkg/utils/sets"
)

func convertDatabaseToSqlSchema(database types.Database) sqlschema2.SQLSchema {
	return sqlschema2.SQLSchema{
		Name: database.ID,
	}
}

func expandSubForms(formDef *types.FormDefinition) []*types.FormDefinition {
	var result []*types.FormDefinition
	result = append(result, formDef)
	for _, field := range formDef.Fields {
		if field.FieldType.SubForm != nil {
			result = append(result, expandSubFormsInternal(formDef, field.Name, field.FieldType.SubForm)...)
			formDef.RemoveField(field.Name)
		}
	}
	return result
}

func expandSubFormsInternal(parentForm *types.FormDefinition, fieldName string, subForm *types.FieldTypeSubForm) []*types.FormDefinition {
	var result []*types.FormDefinition

	formDef := &types.FormDefinition{
		ID:         subForm.ID,
		Code:       subForm.Code,
		Name:       subForm.Name,
		DatabaseID: parentForm.DatabaseID,
		Fields:     subForm.Fields,
	}

	formDef.RemoveField(fieldName)

	formDef.Fields = append(formDef.Fields, &types.FieldDefinition{
		ID:       "parent_id",
		Code:     "parent_id",
		Name:     "parent_id",
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
			result = append(result, expandSubFormsInternal(formDef, field.Name, field.FieldType.SubForm)...)
		}
	}

	return result
}

func convertFormToSqlTable(formDef *types.FormDefinition, referencedForms *types.FormDefinitionList) (sqlschema2.SQLTable, error) {
	table := sqlschema2.SQLTable{}

	table.Name = formDef.ID
	table.Schema = formDef.DatabaseID

	table.Fields = append(table.Fields, sqlschema2.SQLField{
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

	table.Fields = append(table.Fields, sqlschema2.SQLField{
		Name: "seq",
		DataType: sqlschema2.SQLDataType{
			Serial: &sqlschema2.SQLDataTypeSerial{},
		},
	})

	table.Fields = append(table.Fields, sqlschema2.SQLField{
		Name: "database_id",
		DataType: sqlschema2.SQLDataType{
			VarChar: &sqlschema2.SQLDataTypeVarChar{
				Length: 36,
			},
		},
	})

	table.Fields = append(table.Fields, sqlschema2.SQLField{
		Name: "form_id",
		DataType: sqlschema2.SQLDataType{
			VarChar: &sqlschema2.SQLDataTypeVarChar{
				Length: 36,
			},
		},
	})

	table.Fields = append(table.Fields, sqlschema2.SQLField{
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

	table.Constraints = append(table.Constraints, sqlschema2.SQLTableConstraint{
		Name: fmt.Sprintf("fk_%s_forms", table.Name),
		ForeignKey: &sqlschema2.SQLTableConstraintForeignKey{
			ColumnNames: []string{
				"database_id",
				"form_id",
			},
			ReferencesSchema:  "public",
			ReferencesTable:   "forms",
			ReferencesColumns: []string{"database_id", "id"},
		},
	})

	//expandedFields, err := formDef.Fields.Expand(referencedForms)
	//if err != nil {
	//	return schema.SQLTable{}, err
	//}

	keyFieldIDs := sets.NewString()
	for _, field := range formDef.Fields {
		table.Fields = append(table.Fields, convertFieldToSqlField(formDef, field))
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

func convertFieldToSqlField(formDef *types.FormDefinition, field *types.FieldDefinition) sqlschema2.SQLField {
	result := sqlschema2.SQLField{}
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
	} else if field.FieldType.MultilineText != nil {
		result.DataType = sqlschema2.SQLDataType{
			Text: &sqlschema2.SQLDataTypeText{},
		}
	} else if field.FieldType.Date != nil {
		result.DataType = sqlschema2.SQLDataType{
			Date: &sqlschema2.SQLDataTypeDate{},
		}
	} else if field.FieldType.Month != nil {
		result.DataType = sqlschema.SQLDataType{
			Date: &sqlschema.SQLDataTypeDate{},
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
			Name: fmt.Sprintf("fk__%s_%s_%s__%s_%s_id",
				formDef.DatabaseID,
				formDef.ID,
				field.ID,
				field.FieldType.Reference.DatabaseID,
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
				Table:     field.FieldType.SubForm.ID,
				Column:    "id",
				OnDelete:  sqlschema2.ActionRestrict,
				OnUpdate:  sqlschema2.ActionCascade,
				MatchType: "",
			},
		})
	}

	return result
}
