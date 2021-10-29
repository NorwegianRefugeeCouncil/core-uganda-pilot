package sqlconvert

import (
	"fmt"
	"github.com/nrc-no/core/pkg/sets"
	"github.com/nrc-no/core/pkg/sqlschema"
	"github.com/nrc-no/core/pkg/types"
)

func convertDatabaseToSqlSchema(database types.Database) sqlschema.SQLSchema {
	return sqlschema.SQLSchema{
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

func getSubFormName(parentForm *types.FormDefinition, fieldName string) string {
	return fmt.Sprintf("%s_%ss", parentForm.Name, fieldName)
}

func convertFormToSqlTable(formDef *types.FormDefinition, referencedForms *types.FormDefinitionList) (sqlschema.SQLTable, error) {
	table := sqlschema.SQLTable{}

	table.Name = formDef.ID
	table.Schema = formDef.DatabaseID

	table.Fields = append(table.Fields, sqlschema.SQLField{
		Name: "id",
		DataType: sqlschema.SQLDataType{
			VarChar: &sqlschema.SQLDataTypeVarChar{Length: 36},
		},
		Constraints: []sqlschema.SQLColumnConstraint{
			{
				Name:       fmt.Sprintf("pk_%s_%s", formDef.DatabaseID, formDef.ID),
				PrimaryKey: &sqlschema.PrimaryKeySQLColumnConstraint{},
			},
		},
	})

	table.Fields = append(table.Fields, sqlschema.SQLField{
		Name: "seq",
		DataType: sqlschema.SQLDataType{
			Serial: &sqlschema.SQLDataTypeSerial{},
		},
	})

	table.Fields = append(table.Fields, sqlschema.SQLField{
		Name: "database_id",
		DataType: sqlschema.SQLDataType{
			VarChar: &sqlschema.SQLDataTypeVarChar{
				Length: 36,
			},
		},
	})

	table.Fields = append(table.Fields, sqlschema.SQLField{
		Name: "form_id",
		DataType: sqlschema.SQLDataType{
			VarChar: &sqlschema.SQLDataTypeVarChar{
				Length: 36,
			},
		},
	})

	table.Fields = append(table.Fields, sqlschema.SQLField{
		Name:    "created_at",
		Default: "NOW()",
		DataType: sqlschema.SQLDataType{
			Timestamp: &sqlschema.SQLDataTypeTimestamp{
				Timezone: &sqlschema.TimestampWithoutTimeZone,
			},
		},
		Constraints: []sqlschema.SQLColumnConstraint{
			{
				NotNull: &sqlschema.NotNullSQLColumnConstraint{},
			},
		},
	})

	table.Constraints = append(table.Constraints, sqlschema.SQLTableConstraint{
		Name: fmt.Sprintf("fk_%s_forms", table.Name),
		ForeignKey: &sqlschema.SQLTableConstraintForeignKey{
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
	//	return sqlschema.SQLTable{}, err
	//}

	keyFieldIDs := sets.NewString()
	for _, field := range formDef.Fields {
		table.Fields = append(table.Fields, convertFieldToSqlField(formDef, field))
		if field.Key {
			keyFieldIDs.Insert(field.ID)
		}
	}

	if !keyFieldIDs.IsEmpty() {
		table.Constraints = append(table.Constraints, sqlschema.SQLTableConstraint{
			Name: fmt.Sprintf("uq_%s", table.Name),
			Unique: &sqlschema.SQLTableConstraintUnique{
				ColumnNames: keyFieldIDs.List(),
			},
		})
	}

	return table, nil
}

func convertFieldToSqlField(formDef *types.FormDefinition, field *types.FieldDefinition) sqlschema.SQLField {
	result := sqlschema.SQLField{}
	result.Name = field.ID
	result.Comment = field.Code

	if field.Required || field.Key {
		result.Constraints = append(result.Constraints, sqlschema.SQLColumnConstraint{
			NotNull: &sqlschema.NotNullSQLColumnConstraint{},
		})
	}

	if field.FieldType.Text != nil {
		result.DataType = sqlschema.SQLDataType{
			VarChar: &sqlschema.SQLDataTypeVarChar{
				Length: 1024,
			},
		}
	} else if field.FieldType.MultilineText != nil {
		result.DataType = sqlschema.SQLDataType{
			Text: &sqlschema.SQLDataTypeText{},
		}
	} else if field.FieldType.Date != nil {
		result.DataType = sqlschema.SQLDataType{
			Date: &sqlschema.SQLDataTypeDate{},
		}
	} else if field.FieldType.Quantity != nil {
		result.DataType = sqlschema.SQLDataType{
			Int: &sqlschema.SQLDataTypeInt{},
		}
	} else if field.FieldType.SingleSelect != nil {
		result.DataType = sqlschema.SQLDataType{
			VarChar: &sqlschema.SQLDataTypeVarChar{
				Length: 1024,
			},
		}
	} else if field.FieldType.Reference != nil {
		result.Constraints = append(result.Constraints, sqlschema.SQLColumnConstraint{
			Name: fmt.Sprintf("fk__%s_%s_%s__%s_%s_id",
				formDef.DatabaseID,
				formDef.ID,
				field.ID,
				field.FieldType.Reference.DatabaseID,
				field.FieldType.Reference.FormID,
			),
			Reference: &sqlschema.ReferenceSQLColumnConstraint{
				Schema:   field.FieldType.Reference.DatabaseID,
				Table:    field.FieldType.Reference.FormID,
				Column:   "id",
				OnUpdate: sqlschema.ActionCascade,
				OnDelete: sqlschema.ActionCascade,
			},
		})
		result.DataType = sqlschema.SQLDataType{
			VarChar: &sqlschema.SQLDataTypeVarChar{Length: 36},
		}
	} else if field.FieldType.SubForm != nil {
		result.Name += "_id"
		result.DataType = sqlschema.SQLDataType{
			VarChar: &sqlschema.SQLDataTypeVarChar{Length: 36},
		}
		result.Constraints = append(result.Constraints, sqlschema.SQLColumnConstraint{
			Reference: &sqlschema.ReferenceSQLColumnConstraint{
				Schema:    formDef.DatabaseID,
				Table:     field.FieldType.SubForm.ID,
				Column:    "id",
				OnDelete:  sqlschema.ActionRestrict,
				OnUpdate:  sqlschema.ActionCascade,
				MatchType: "",
			},
		})
	}

	return result
}
