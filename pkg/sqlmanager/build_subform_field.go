package sqlmanager

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
)

// subFormFieldActions returns the SQL actions necessary to build a SubForm field
// A SubForm field does not actually create a column in the table, but rather creates
// an entirely new SQL table.
func subFormFieldActions(formInterface types.FormInterface, fieldDefinition *types.FieldDefinition) (sqlActions, error) {
	result := sqlActions{}

	// Get the form interface for the SubForm
	subFormInterface, err := formInterface.FindSubForm(fieldDefinition.ID)
	if err != nil {
		return nil, err
	}
	if subFormInterface == nil {
		return nil, fmt.Errorf("failed to find subForm with id %s", fieldDefinition.ID)
	}

	// Get the actions for creating the SubForm form
	actionsForTable, err := getSQLActionsForForm(subFormInterface)
	if err != nil {
		return nil, err
	}

	result = append(result, actionsForTable...)
	return result, nil
}

// buildSubFormOwnerColumn creates the "owner_id" column. This is only applicable for SubForms.
func buildSubFormOwnerColumn(formInterface types.SubFormInterface) schema.SQLColumn {
	return schema.SQLColumn{
		Name: keyOwnerIdColumn,
		DataType: schema.SQLDataType{
			VarChar: &schema.SQLDataTypeVarChar{
				Length: uuidFieldLength,
			},
		},
		Constraints: []schema.SQLColumnConstraint{
			{
				// The owner_id SQL column cannot accept null values
				NotNull: &schema.NotNullSQLColumnConstraint{},
			}, {
				// The owner_id SQL column must reference an existing parent
				// This adds a "foreign_key" constraint to the parent table
				Reference: &schema.ReferenceSQLColumnConstraint{
					Schema:   formInterface.GetDatabaseID(),
					Table:    formInterface.GetOwnerFormID(),
					Column:   keyIdColumn,
					OnDelete: schema.ActionCascade,
					OnUpdate: schema.ActionCascade,
				},
			},
		},
	}
}
